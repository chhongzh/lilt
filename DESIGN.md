TinyGo GUI Framework 设计需求文档

1. 项目概述
   本项目旨在基于 TinyGo 语言构建一套高性能、轻量级且具备高度可扩展性的图形用户界面（GUI）框架。该框架将借鉴 Android、Flutter 及 Chrome 的组件化渲染逻辑，提供类似 osu!lazer 的现代化 UI 体验，支持跨平台运行。
   lilt 是一个抽象的 UI 框架，与用户的接线硬件无关：屏幕和按键都是接口，屏幕不一定是单色的。

2. 核心架构设计

2.1 组件化与接口抽象
- 能力靠接口，组件靠继承：
    - 框架通过类型断言发现节点的能力。一个节点会画、会动、会测量，都是因为它实现了对应的接口，而不是因为它属于某个基类。
    - 组件通过嵌入（Go 的继承）复用实现：所有节点嵌入 `ui.Box`；行类组件嵌入 `ListItem`，`CheckBox`、`Radio`、`Select` 都建立在它上面；可交互组件额外嵌入 `widget.Control`，从而获得事件和状态机。
- 核心接口（`ui` 包）：

| 原需求里的名字 | 现在的名字 | 作用 |
| --- | --- | --- |
| （节点本身） | `Node` | 树里的任何东西。只有嵌入了 `Box` 的类型才满足它 |
| Drawable | `Painter` / `OverPainter` | 在子节点之前 / 之后绘制自己 |
| Updatable | `Updater` | 每帧推进一次逻辑，`Update(dt)` |
| Position | `Box` 的 `Anchor`、`Origin`、`Position`、`OffsetX/Y` | 相对父节点的位置，加上可动画的偏移 |
| Rect | `ui.Rect`（布局，float32）/ `gfx.Rect`（像素） | 节点的布局矩形，也是点击和焦点查找用的范围 |
| （测量） | `Measurer` / `Layouter` | 叶子报告内容尺寸 / 容器自己摆放子节点（Flex、Grid） |

- 其余能力接口：`KeyHandler`、`KeyPreviewer`、`PointerHandler`、`Focuser`、`Revealer`、`ArrowScroller`、`CursorShaper`、`ContentHost`、`Leaf`、`ChildRemover`。
- `widget` 包的选项同样按能力匹配：`Font` 作用于任何 `FontSetter`，`TextColor` 作用于任何 `Colorer`，`OnClick` / `OnChange` 作用于任何嵌入了 `Control` 的组件。

2.2 全局单例与管理器
整个 lilt 是单例：一块屏幕、一棵树，每种管理器各一个。外部只保留生命周期三个函数：

- `lilt.Start(screen, root)`：绑定屏幕，把 root 作为第一屏，启动 UI goroutine。之后框架自己管理帧循环。调用方要么 `select {}` 阻塞，要么在这里继续跑自己的业务代码（例如读引脚）。已暂停时再调用 Start 即恢复。
- `lilt.Pause()`：停止更新和绘制，状态保留。
- `lilt.Stop()`：结束 UI goroutine，丢弃所有屏幕，之后可以重新 Start。

管理器（原需求名 → 现在的名字）：

| 原需求里的名字 | 现在的名字 | 职责 |
| --- | --- | --- |
| ViewManager | 层栈中的页面，`ui.Push` / `ui.Pop` | 整页的导航栈，统一处理飞入飞出（通过 Offset） |
| （新增）DialogManager | 层栈中的弹窗，`ui.Open` / `ui.Close` | 统一的弹窗能力：四个方向的 Drawer、居中 Dialog、不抢焦点的 Toast |
| FocusManager | 焦点管理（`ui.SetFocus` / `ui.Focused`） | 焦点的获取、丢失与传递 |
| DisplayManager | `gfx.Screen` | 显示输出接口，驱动实现它即可接入 |
| RenderManager | 帧循环（`lilt` 包驱动） | 输入、更新、动画、布局、绘制，每帧按此顺序 |
| InputManager | 输入管线（`ui.Emit` / `EmitRune` / `EmitPointer`） | 统一的输入入口，任何 goroutine 都能调用 |
| Keyboard | `ui.TextInputMethod`，默认 `widget.Keyboard` | 文本输入的统一抽象，软键盘只是其中一种实现 |
| （新增） | `ui.KeyPoller` | 把原始引脚电平变成消抖后的按下、连发、松开 |

框架与硬件无关：用户自己读引脚，再调用 Emit 或 KeyPoller 交给框架。

2.3 节点树与生命周期
- Parent-Child 关系：所有 UI 元素（Child）必须拥有父节点（Parent），形成树状结构。
- 根节点：全局单例，层栈里的每一层都挂在它下面。
- 自我销毁机制：节点可以 `ui.Remove` 自己，父节点会收到 `ChildRemoved` 通知。
- Parent 也可以移除 Child，但要看组件自己实现：实现了 `Leaf` 的节点（例如 Label）拒绝子节点。
- 焦点控制：节点可配置是否可被聚焦（Focusable）。滚动容器（Scroll）处理焦点穿透与捕获：先让焦点在内容里移动并滚动显示，到边界再交出去。

2.4 线程模型
- 树只属于 UI goroutine。其他 goroutine 改树必须包在 `ui.Do(fn)` 里，它在两帧之间拿到帧锁再执行 fn。fn 里不能再调用 Do。
- 用户回调永远不在 UI goroutine 上跑。OnClick、OnChange 各自在独立的 goroutine 里执行，慢回调不会卡住动画；回调要改界面就用 `ui.Do`。
- 实现框架接口的代码（自定义组件、`TextTarget`、`TextInputMethod`）属于组件本身，在 UI goroutine 上运行。

3. 渲染与布局系统

3.1 布局哲学
- 相对坐标系统：尺寸用 `Rel`（父内容区的比例）或 `Auto`（按内容），位置用锚点和原点。
- 相对间距：Padding、Margin、Gap 和圆角都以“单位”给出，一个单位等于主题字体行高与默认字体行高之比。默认字体下一个单位就是一个像素；换大字体配大屏时，所有间距跟着放大。只有 `Px` 尺寸是绝对像素。
- 统一风格：行、输入框、标题的文字内边距取主题的 `RowPadding`，按钮取 `ButtonPadding`，并排控件之间取 `Gap`，边框统一用前景色加主题 `Radius`，次要文字统一用小字体加 `Muted` 颜色。
- DOM Tree 机制：任何属性变化都标记重新布局，下一帧整树重排。
- 无响应式样式：不采用 CSS 响应式设计思路，而是通过代码逻辑控制布局。
- 数据访问：属性一律走 Set/Get 方法，确保数据流向清晰。

3.2 布局组件支持
- Flex & Grid：支持弹性盒子（Flex）与网格（Grid）布局算法，另有按锚点摆放的 Stack。
- Flutter-like 创建方式：组件创建采用嵌套式风格，构造函数接收 `ui.Option`。
- 指针传递：组件创建与传递全部使用指针，节省内存并方便管理对象生命周期。

3.3 动画系统
- `anim.Value`：动画系统的核心原语，支持数值插值，可以有 Parent，用来叠加反馈变换而不必手动管理。
- Parent Offset 继承：子元素相对父元素绘制，自然继承父元素的偏移，用于 Push 等转场。
- Delta Based 动画：基于时间增量（Delta Time）计算动画帧，确保不同帧率下速度一致。
- 物理效果：统一的阻尼弹簧，参数采用 SwiftUI 的 Response / Damping 写法。预设（`anim.Default`、`Smooth`、`Snappy`、`Bouncy`、`Interactive`）是指针，运行时改一处全局生效。
- 列表滑到顶部后用户仍然往上，整个列表给一个向下的 push 视觉效果，松开弹回。

3.4 层栈与弹窗
- 页面和弹窗共用一个按绘制顺序排列的层栈。只有最上面的模态层接收焦点和指针；被其他页面完全盖住的页面隐藏，不参与更新和绘制。
- 页面：`ui.Push` 从右侧滑入，下层页面退到左边 1/3 处。
- Drawer：`ui.Open(n, ui.DrawerBottom)` 等四个方向，从对应边缘滑入，沿边铺满，另一方向按内容大小。`widget.Drawer` 是带标题的列容器，`Select` 是 Drawer 的封装。
- Dialog：从上方落到屏幕中央。`widget.Dialog` 有标题、一行消息和一排按钮，按钮回调成功后自动关闭。
- Toast：底部浮起，不抢焦点也不接按键，画在光标之上。
- Back：按下时最上层朝它离开的方向探出一点，松手关闭；根页面只轻轻顶一下。

4. 交互与输入系统

4.1 焦点查找算法
- FocusFinder：类似 Android 的焦点查找，按方向键在“光束”内优先，再比较 13·主轴² + 次轴² 距离。

4.2 事件分发机制
- 按键先自上而下给 `KeyPreviewer` 预览，再从焦点节点自下而上冒泡给 `KeyHandler`。认领了 Press 的节点独占该键后续的 Repeat 和 Release。
- 所有可激活的组件都在 OK 的 Release 且未取消时触发；Press 只给反馈（光标收缩）。
- 事件：
    - `OnClick(fn)`：组件被激活时调用。fn 可以是 `func()` 或 `func() error`。
    - `OnChange(fn)`：用户改变了组件的值时调用，参数是旧值和新值 `func(from, to T)`，也可以返回 error。T 由组件决定：CheckBox 是 bool，Choice、Select、RadioGroup 是 int，Input 是 string，Toggles 是 uint32 位掩码。代码里调用 Set 改值不会触发 OnChange。
    - 两者都是 `Control` 上的泛型方法（Go 1.27，TinyGo 0.42 起支持），也有同名的 `widget.OnClick` / `widget.OnChange` 选项。一次激活里先跑 OnChange（值变了才跑），再跑 OnClick。
- SingleFlight：回调在跑时，组件忽略新的激活。

4.3 组件状态机
每个可交互组件都按同一个状态机运行，用户不用自己写：

- Normal：可以激活。
- Busy：回调已经在后台开始，组件忽略激活，外观不变。值类组件已经乐观地显示新值。
- Loading：回调跑了超过 150ms，显示加载点。按钮的文字换成加载点，行的尾部控件换成加载点，Chip 的横条换成加载点。
- 回调结束：成功则回到 Normal（按钮如果显示过加载，会先亮一个勾）；返回 error 则组件抖一下，值类组件恢复旧值，按钮显示叉。

4.4 光标与反馈
- 全局统一光标：类似 iOS 的指针吸附效果，换目标时先朝新位置拉长再收拢，跨屏时跳过去。
- 系统级键盘唤起：Input 激活时打开 `ui.TextInput`，没设置就用软键盘。

5. 标准组件库

框架提供一套开箱即用的基础组件，保持统一的设计语言：

- Button：
    - 自动集成 Loading 动画（见 4.3）。
    - 支持长按确认（类似 osu!lazer 的确认按钮：长按时光标从左往右缩短，长度代表确认进度，缩完按钮才触发）。
    - 内置状态切换动画：普通态、按下态、加载态、成功态、失败态。
- 基础控件：Label（放不下就跑马灯）、ListItem、List、CheckBox、Radio、Select、Chip、Choice、Toggles、Input、软键盘。
- 弹窗：Drawer（四个方向）、Dialog、Toast。
- Scroll：
    - 支持视口裁切（Clipping）。
    - 支持水平、垂直及双向滚动。
    - 支持嵌套：内层滚到头且没有可聚焦的下一项时，把按键交给外层还能滚的滚动容器，而不是原地回弹。

6. 参考设计与技术对标

- UI 风格：参考 osu!lazer 的现代化、高动态 UI 框架。
- 交互逻辑：深度对标 Android 的焦点管理与 View 系统。
- 渲染架构：参考 Flutter 的组件树构建方式与 Chrome 的渲染管线逻辑。
