# ui

[pkg.go.dev](https://pkg.go.dev/github.com/chhongzh/lilt/ui) · `import "github.com/chhongzh/lilt/ui"`

节点树、布局、层栈、焦点、输入、光标和主题。除了标注“任何 goroutine”的函数，其余都只能在 UI goroutine 上调用：`lilt.Start` 之前、`ui.Do` 里、或者组件的方法里。

## 节点与 Box

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Node interface` | 树上的任何节点；只有一个未导出方法，只有嵌入了 `Box` 的类型满足它 | [节点树](../concepts/tree.md) |
| `type Box struct` | 每个节点都嵌入的基础：父子链接、尺寸、内边距、动画偏移 | [节点树](../concepts/tree.md) |
| `func BoxOf(d Node) *Box` | 取 `d` 嵌入的 `Box` | [节点树](../concepts/tree.md) |
| `func (b *Box) Parent() Node` | 父节点 | [节点树](../concepts/tree.md) |
| `func (b *Box) Children() []Node` | 子节点，不要直接修改，用 `Add`、`Insert`、`Remove` | [节点树](../concepts/tree.md) |
| `func (b *Box) Width() Size`、`SetWidth(s Size)` | 宽度 | [尺寸](../concepts/layout.md#size) |
| `func (b *Box) Height() Size`、`SetHeight(s Size)` | 高度 | [尺寸](../concepts/layout.md#size) |
| `func (b *Box) Anchor() Vec2`、`SetAnchor(v Vec2)` | 节点挂在父节点内容区的哪一点 | [锚点](../concepts/layout.md#anchor) |
| `func (b *Box) Origin() Vec2`、`SetOrigin(v Vec2)` | 节点自己的哪一点放在锚点上 | [锚点](../concepts/layout.md#anchor) |
| `func (b *Box) Position() Vec2`、`SetPosition(v Vec2)` | 锚定之后再加的偏移，像素 | [锚点](../concepts/layout.md#anchor) |
| `func (b *Box) Margin() Insets`、`SetMargin(i Insets)` | 外边距，单位 | [尺寸](../concepts/layout.md#size) |
| `func (b *Box) Padding() Insets`、`SetPadding(i Insets)` | 内边距，单位 | [尺寸](../concepts/layout.md#size) |
| `func (b *Box) Grow() float32`、`SetGrow(g float32)` | 在 `Flex` 里分剩余空间的份数 | [Flex](../concepts/layout.md#flex) |
| `func (b *Box) Visible() bool`、`SetVisible(v bool)` | 隐藏的节点不占空间、不绘制、不能有焦点 | [节点树](../concepts/tree.md) |
| `func (b *Box) Focusable() bool`、`SetFocusable(f bool)` | 能否获得焦点 | [焦点](../concepts/input.md#focus) |
| `func (b *Box) SetClip(c bool)` | 把自己和子节点的绘制限制在自己的矩形里 | [能力接口](../guides/custom-widget.md#interfaces) |
| `func (b *Box) SetAnimateLayout(a bool)` | 布局引起的位置变化改为滑过去 | [布局动画](../concepts/layout.md#animate-layout) |
| `func (b *Box) SetBackground(c gfx.Color)` | 用 `c` 填充矩形，按圆角 | [主题](../concepts/theme.md#theme) |
| `func (b *Box) SetBorder(c gfx.Color)` | 一像素的边框 | [主题](../concepts/theme.md#theme) |
| `func (b *Box) Bordered() bool` | 是否画边框 | [主题](../concepts/theme.md#theme) |
| `func (b *Box) SetRadius(r int16)` | 背景和边框的圆角，单位 | [单位](../concepts/theme.md#unit) |
| `func (b *Box) Rect() Rect` | 布局得到的矩形，相对父节点，不含动画偏移 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (b *Box) ContentRect() Rect` | 内边距以内的区域，像素，自己的坐标系 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (b *Box) OffsetX() *anim.Value`、`OffsetY() *anim.Value` | 叠加在布局位置上的动画平移，子节点跟着动 | [动画](../guides/custom-widget.md#animation) |

## 构造选项

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Option func(Node)` | 构造时配置节点，让整棵树能写成一个嵌套表达式 | [选项](../concepts/tree.md#options) |
| `func Apply(d Node, opts []Option)` | 依次应用选项；组件的构造函数最后调用它 | [自定义组件](../guides/custom-widget.md#type) |
| `func Children(nodes ...Node) Option` | 添加子节点 | [选项](../concepts/tree.md#options) |
| `func Ref[T Node](p *T) Option` | 把正在构造的节点存进 `*p`，在嵌套的树里留个句柄 | [Ref](../concepts/tree.md#ref) |
| `func Width(s Size) Option`、`Height(s Size) Option` | 宽、高 | [尺寸](../concepts/layout.md#size) |
| `func Fill() Option` | 宽高都等于父节点内容区 | [尺寸](../concepts/layout.md#size) |
| `func FillWidth() Option`、`FillHeight() Option` | 只铺满一个方向 | [尺寸](../concepts/layout.md#size) |
| `func Anchor(v Vec2) Option`、`Origin(v Vec2) Option` | 锚点、原点 | [锚点](../concepts/layout.md#anchor) |
| `func Align(v Vec2) Option` | 锚点和原点设成同一点，把节点钉到父节点的边或角 | [锚点](../concepts/layout.md#anchor) |
| `func Position(x, y float32) Option` | 附加偏移 | [锚点](../concepts/layout.md#anchor) |
| `func Margin(i Insets) Option`、`Padding(i Insets) Option` | 外边距、内边距 | [尺寸](../concepts/layout.md#size) |
| `func Grow(g float32) Option` | `Flex` 里的伸展份数 | [Flex](../concepts/layout.md#flex) |
| `func Gap(g float32) Option` | `Flex` 或 `Grid` 子节点间距，单位；用在别的节点上 panic | [Flex](../concepts/layout.md#flex) |
| `func Main(m MainAlign) Option`、`Cross(c CrossAlign) Option` | `Flex` 主轴、交叉轴对齐 | [Flex](../concepts/layout.md#flex) |
| `func Columns(cols ...Track) Option` | `Grid` 的列 | [Grid](../concepts/layout.md#grid) |
| `func Focusable() Option` | 可获得焦点 | [焦点](../concepts/input.md#focus) |
| `func Hidden() Option` | 一开始就隐藏 | [节点树](../concepts/tree.md) |
| `func Clip() Option` | 裁剪到自己的矩形 | [能力接口](../guides/custom-widget.md#interfaces) |
| `func AnimateLayout() Option` | 布局变化时滑动 | [布局动画](../concepts/layout.md#animate-layout) |
| `func Background(c gfx.Color) Option`、`Border(c gfx.Color) Option` | 背景色、边框色 | [主题](../concepts/theme.md#theme) |
| `func Radius(r int16) Option` | 圆角，单位 | [单位](../concepts/theme.md#unit) |

## 树操作

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Add(parent Node, children ...Node)` | 追加子节点；`parent` 实现了 `ContentHost` 时加到它的内部节点 | [节点树](../concepts/tree.md) |
| `func Insert(parent Node, i int, child Node)` | 插入到第 `i` 个位置 | [节点树](../concepts/tree.md) |
| `func Attach(parent, child Node)` | 直接追加，忽略 `ContentHost`；组件搭内部结构用 | [节点树](../concepts/tree.md) |
| `func Remove(d Node)` | 从父节点摘下并通知父节点；焦点在里面时移到最近的剩余节点 | [节点树](../concepts/tree.md) |
| `func RemoveChildren(parent Node)` | 摘下所有子节点（或 `Content` 的） | [节点树](../concepts/tree.md) |
| `func IsAncestor(a, d Node) bool` | `a` 是 `d` 本身或包含 `d` | [节点树](../concepts/tree.md) |
| `func Attached(d Node) bool` | `d` 是否在活动的树上 | [节点树](../concepts/tree.md) |

## 线程与刷新

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Do(fn func())` | 在两帧之间运行 `fn`，其他 goroutine 借此改树；任何 goroutine，`fn` 里不能再调用 `Do` | [ui.Do](../concepts/threading.md#do) |
| `func Relayout()` | 安排一次整树布局和重绘 | [布局](../concepts/layout.md) |
| `func Invalidate()` | 安排一次重绘 | [布局](../concepts/layout.md) |

## 尺寸与几何

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Size struct` | 一个方向的尺寸：固定像素、父节点内容区的比例，或 `Auto` | [尺寸](../concepts/layout.md#size) |
| `func Px(v float32) Size` | 固定像素 | [尺寸](../concepts/layout.md#size) |
| `func Rel(f float32) Size` | 父节点内容区的比例 | [尺寸](../concepts/layout.md#size) |
| `var Auto = Size{}` | 按内容；`Size` 的零值 | [尺寸](../concepts/layout.md#size) |
| `const Unbounded float32 = 1 << 24` | 滚动方向上的可用空间；这里的比例尺寸退回 `Auto` | [尺寸](../concepts/layout.md#size) |
| `type Vec2 struct{ X, Y float32 }` | 布局空间的点或尺寸 | [锚点](../concepts/layout.md#anchor) |
| `TopLeft`、`TopCenter`、`TopRight`、`CenterLeft`、`Center`、`CenterRight`、`BottomLeft`、`BottomCenter`、`BottomRight` | 九个对齐点，`Vec2{0, 0}` 到 `Vec2{1, 1}` | [锚点](../concepts/layout.md#anchor) |
| `type Rect struct{ X, Y, W, H float32 }` | 布局矩形 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (r Rect) Right() float32`、`Bottom() float32` | 右边、下边 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (r Rect) Center() Vec2` | 中心点 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (r Rect) Pixels() gfx.Rect` | 四条边分别取整，相邻矩形取整后仍相邻 | [测量与绘制](../guides/custom-widget.md#measure) |
| `type Insets struct{ Left, Top, Right, Bottom float32 }` | 外边距或内边距，单位 | [单位](../concepts/theme.md#unit) |
| `func All(v float32) Insets` | 四边相同 | [单位](../concepts/theme.md#unit) |
| `func XY(x, y float32) Insets` | 左右 `x`、上下 `y` | [单位](../concepts/theme.md#unit) |
| `func LTRB(l, t, r, b float32) Insets` | 分别指定 | [单位](../concepts/theme.md#unit) |
| `func (i Insets) Horizontal() float32`、`Vertical() float32` | 左右之和、上下之和 | [单位](../concepts/theme.md#unit) |
| `func (i Insets) Scale(f float32) Insets` | 四边乘 `f`，乘 `Unit()` 得到像素 | [单位](../concepts/theme.md#unit) |
| `func Unit() float32` | 一个单位等于多少像素：主题字体行高除以默认字体行高 | [单位](../concepts/theme.md#unit) |
| `func Round(x float32) int16` | 四舍五入，.5 远离零 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func ScreenSize() (w, h int16)` | 绑定的屏幕尺寸，像素 | [屏幕接口](../guides/hardware.md#screen) |

## 容器

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Flex struct` | 沿一个方向排列子节点 | [Flex](../concepts/layout.md#flex) |
| `func NewRow(opts ...Option) *Flex`、`NewColumn(opts ...Option) *Flex` | 横向、纵向的 `Flex` | [Flex](../concepts/layout.md#flex) |
| `func (f *Flex) Axis() Axis`、`SetAxis(a Axis)` | 方向；嵌入 `Flex` 的组件用 `SetAxis` 变成一列 | [Flex](../concepts/layout.md#flex) |
| `func (f *Flex) Gap() float32`、`SetGap(g float32)` | 子节点间距，单位 | [Flex](../concepts/layout.md#flex) |
| `func (f *Flex) SetMainAlign(m MainAlign)`、`SetCrossAlign(c CrossAlign)` | 主轴、交叉轴对齐 | [Flex](../concepts/layout.md#flex) |
| `type Axis uint8`：`Horizontal`、`Vertical` | `Flex` 的方向 | [Flex](../concepts/layout.md#flex) |
| `type MainAlign uint8`：`MainStart`、`MainCenter`、`MainEnd`、`MainSpaceBetween` | 没有子节点伸展时，主轴上的剩余空间怎么分 | [Flex](../concepts/layout.md#flex) |
| `type CrossAlign uint8`：`CrossStretch`、`CrossStart`、`CrossCenter`、`CrossEnd` | 交叉轴上的位置；`CrossStretch` 拉满该方向为 `Auto` 的子节点 | [Flex](../concepts/layout.md#flex) |
| `type Grid struct` | 按列逐行摆放，行高取该行最高的子节点，`Auto` 尺寸的子节点填满格子 | [Grid](../concepts/layout.md#grid) |
| `func NewGrid(opts ...Option) *Grid` | 没设置列时为一列 `Fr(1)` | [Grid](../concepts/layout.md#grid) |
| `func (g *Grid) SetColumns(cols ...Track)` | 设置列 | [Grid](../concepts/layout.md#grid) |
| `func (g *Grid) SetGap(gap float32)` | 行列间距，单位 | [Grid](../concepts/layout.md#grid) |
| `type Track struct` | 一列的宽度 | [Grid](../concepts/layout.md#grid) |
| `func Fixed(px float32) Track` | 固定像素 | [Grid](../concepts/layout.md#grid) |
| `func Fr(n float32) Track` | 剩余宽度的份数 | [Grid](../concepts/layout.md#grid) |
| `type Stack struct` | 普通容器，子节点按锚点和原点摆放 | [Stack](../concepts/layout.md#stack) |
| `func NewStack(opts ...Option) *Stack` | 新建 `Stack` | [Stack](../concepts/layout.md#stack) |
| `Flex`、`Grid` 的 `MeasureChildren`、`ArrangeChildren` | `Layouter` 的实现，由布局调用 | [能力接口](../guides/custom-widget.md#interfaces) |

## 布局函数

自己实现 `Layouter` 时用来测量、摆放子节点。

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Measure(d Node, max Vec2) Vec2` | `d` 在 `max`（扣除外边距后）里想要的尺寸 | [能力接口](../guides/custom-widget.md#interfaces) |
| `func Arrange(d Node, r Rect)` | 给 `d` 最终矩形（相对父节点）并布局它的子节点 | [能力接口](../guides/custom-widget.md#interfaces) |
| `func ArrangeFree(c Node, area Rect)` | 按锚点、原点和偏移把 `c` 放进 `area`；非 `Layouter` 节点的默认摆法 | [Stack](../concepts/layout.md#stack) |
| `func RestRect(d Node) Rect` | 所有动画停下后 `d` 在屏幕上的矩形，像素 | [焦点](../concepts/input.md#focus) |

## 页面与弹窗

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Push(p Node)` | 把 `p` 作为新的整屏页面从右侧推入 | [页面](../concepts/layers.md#push) |
| `func Pop() bool` | 关闭最上层的页面或弹窗；根页面不关 | [Back](../concepts/layers.md#back) |
| `func Open(n Node, p Popup)` | 以弹窗打开 `n`；抽屉和对话框关闭前独占焦点和按键 | [弹窗](../concepts/layers.md#popups) |
| `func Close(n Node) bool` | 关闭显示 `n` 的页面或弹窗 | [弹窗](../concepts/layers.md#popups) |
| `type Popup uint8` | 弹窗的位置和进出方向 | [弹窗](../concepts/layers.md#popups) |
| `DrawerBottom`、`DrawerTop` | 从下、上边缘滑入，高度按内容 | [Drawer](../widgets/popups.md#drawer) |
| `DrawerLeft`、`DrawerRight` | 从左、右边缘滑入，宽度按内容 | [Drawer](../widgets/popups.md#drawer) |
| `Dialog` | 落到屏幕中央 | [Dialog](../widgets/popups.md#dialog) |
| `Toast` | 从底部浮起，不抢焦点、不接按键 | [Toast](../widgets/popups.md#toast) |

## 焦点

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func SetFocus(d Node)` | 把焦点移到 `d`；隐藏或不可聚焦的节点被忽略 | [焦点](../concepts/input.md#focus) |
| `func Focused() Node` | 当前焦点节点，没有时为 nil | [焦点](../concepts/input.md#focus) |
| `func FindFocus(scope, from Node, dir Key) Node` | 按 Android FocusFinder 的规则找 `dir` 方向的下一个焦点 | [焦点](../concepts/input.md#focus) |
| `func FirstFocusable(d Node) Node` | 按树序第一个可聚焦节点 | [焦点](../concepts/input.md#focus) |

## 按键与触摸

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Key uint8` | 逻辑键，引脚到键的映射由板子代码负责 | [按键](../concepts/input.md#keys) |
| `KeyNone`、`KeyUp`、`KeyDown`、`KeyLeft`、`KeyRight`、`KeyOK`、`KeyBack`、`KeyMenu`、`KeyBackspace`、`KeyChar` | 所有逻辑键；`KeyChar` 的字符在 `KeyEvent.Rune` 里 | [按键](../guides/hardware.md#keys) |
| `func (k Key) Arrow() bool` | 是否方向键 | [按键](../concepts/input.md#keys) |
| `type Phase uint8`：`Press`、`Repeat`、`Release` | 按下、连发、松开 | [激活](../concepts/input.md#activation) |
| `type KeyEvent struct{ Key; Phase; Rune; Canceled; Repeated }` | 一个按键事件；`Canceled` 是触摸变拖动后的松开，`Repeated` 是连发之后的松开 | [激活](../concepts/input.md#activation) |
| `func (e *KeyEvent) Fires() bool` | 是否该执行动作：未取消且前面没有连发的 `Release`，或每个 `Repeat` | [激活](../concepts/input.md#activation) |
| `func Emit(k Key, p Phase)` | 发一个按键事件；任何 goroutine | [按键](../guides/hardware.md#keys) |
| `func EmitRune(r rune)` | 发一个键入的字符；任何 goroutine | [物理键盘](../guides/hardware.md#keys) |
| `func EmitPointer(x, y int16, p PointerPhase)` | 发一个触摸或鼠标事件，屏幕像素；任何 goroutine | [触摸屏](../guides/hardware.md#keys) |
| `type PointerEvent struct{ X, Y int16; Phase PointerPhase }` | 一个指针事件 | [分发](../concepts/input.md#dispatch) |
| `type PointerPhase uint8`：`PointerDown`、`PointerMove`、`PointerUp` | 按下、移动、抬起 | [触摸屏](../guides/hardware.md#keys) |
| `type KeyPoller struct` | 把按键电平变成消抖后的按下、连发、松开，连发会加速 | [KeyPoller](../concepts/input.md#keypoller) |
| `func NewKeyPoller() *KeyPoller` | 默认 `Debounce` 2、`RepeatDelay` 0.2、`RepeatInterval` 0.07、`MinInterval` 0.03、`Accel` 0.9 | [KeyPoller](../concepts/input.md#keypoller) |
| `func (p *KeyPoller) Poll(k Key, down bool, dt float32)` | 喂一次采样，`dt` 是距上次的秒数 | [按键](../guides/hardware.md#keys) |
| `KeyPoller.Debounce`、`RepeatDelay`、`RepeatInterval`、`MinInterval`、`Accel` | 消抖采样数、连发延迟、首次间隔、最小间隔、间隔倍率 | [按键](../guides/hardware.md#keys) |

## 文字输入

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type TextTarget interface{ Text() string; SetText(s string) }` | 输入法编辑的对象 | [文字输入](../widgets/input.md#text-input) |
| `type TextInputMethod interface{ Open(t TextTarget) }` | 为 `t` 打开一个编辑器 | [文字输入](../widgets/input.md#text-input) |
| `var TextInput TextInputMethod` | `Input` 激活时打开的编辑器；nil 时用 `widget.Keyboard` | [文字输入](../widgets/input.md#text-input) |

## 能力接口

组件实现这些接口获得对应能力，框架通过类型断言调用。

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Painter interface{ Draw(c gfx.Canvas) }` | 画自己，在子节点之前 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type OverPainter interface{ DrawOver(c gfx.Canvas) }` | 在子节点之后再画 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type Updater interface{ Update(dt float32) }` | 节点在树上且可见时每帧一次 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type Measurer interface{ Measure(max Vec2) Vec2 }` | 叶子节点的内容尺寸，不含内边距 | [测量与绘制](../guides/custom-widget.md#measure) |
| `type Layouter interface{ MeasureChildren; ArrangeChildren }` | 自己测量、摆放子节点；`content` 是内边距以内的区域 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type KeyHandler interface{ HandleKey(e *KeyEvent) bool }` | 按键从焦点节点向上冒泡；`Press` 时返回 true 就认领这个键 | [按键分发](../concepts/input.md#dispatch) |
| `type KeyPreviewer interface{ PreviewKey(e *KeyEvent) bool }` | 按键从根向焦点节点下行时先过目，容器借此拦截 | [按键分发](../concepts/input.md#dispatch) |
| `type PointerHandler interface{ HandlePointer(e *PointerEvent) bool }` | 指针事件从指下节点向上冒泡；移动时返回 true 就接管指针 | [按键分发](../concepts/input.md#dispatch) |
| `type Focuser interface{ FocusChanged(focused bool) }` | 获得或失去焦点时通知 | [焦点](../concepts/input.md#focus) |
| `type Revealer interface{ Reveal(d Node) }` | 后代获得焦点时通知，用来滚动到它 | [焦点](../concepts/input.md#focus) |
| `type ArrowScroller interface{ ArrowScroll(e *KeyEvent, next Node) bool }` | 焦点移动前祖先先拿到方向键；`next` 是焦点将去的节点，边缘时为 nil | [嵌套滚动](../widgets/scroll.md#nested) |
| `type CursorShaper interface{ CursorShape(r Rect) Rect }` | 焦点节点改变光标的矩形，屏幕像素 | [改变光标](../guides/custom-widget.md#cursorshaper) |
| `type ContentHost interface{ Content() Node }` | 让 `Add`、`RemoveChildren` 转到内部节点 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type Leaf interface{ Leaf() }` | 不接受子节点 | [能力接口](../guides/custom-widget.md#interfaces) |
| `type ChildRemover interface{ ChildRemoved(child Node) }` | 子节点摘下时通知，包括子节点移除自己 | [能力接口](../guides/custom-widget.md#interfaces) |

## 光标

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Cursor struct` | 标示焦点的唯一高亮，四条边各自运动，生活在屏幕坐标里 | [光标](../concepts/input.md#cursor) |
| `func (c *Cursor) Lean(dx, dy float32)` | 把光标偏移并保持，按键按下时的预备动作；`Lean(0, 0)` 放开 | [光标](../concepts/input.md#cursor) |
| `func (c *Cursor) Leap()` | 跳一下：抬起、变大、落地压扁；切换屏幕时调用 | [光标](../concepts/input.md#cursor) |
| `func (c *Cursor) Nudge(px float32)` | 横向甩出约 `px` 再弹回 | [光标](../concepts/input.md#cursor) |
| `var CursorLead = &anim.Spring{Response: 0.22, Ease: true}` | 光标朝运动方向那条边的缓动 | [光标](../concepts/input.md#cursor) |
| `var CursorTrail = &anim.Spring{Response: 0.52, Ease: true}` | 光标后面那条边的缓动，更慢，于是先拉长再收拢 | [光标](../concepts/input.md#cursor) |
| `var CursorRelease = &anim.Spring{Response: 0.57, Damping: 0.4}` | 松开 OK 后光标带着晃动恢复；按下时用 `anim.Interactive` | [光标](../concepts/input.md#cursor) |
| `var CursorLeap = &anim.Spring{Response: 0.63, Damping: 0.45}` | 跨屏幕跳跃 | [光标](../concepts/input.md#cursor) |
| `var CursorLean = &anim.Spring{Response: 0.35, Damping: 0.7}` | 朝按住的键偏过去再回来 | [光标](../concepts/input.md#cursor) |

`Cursor` 的实例是包内的单例，没有导出的访问方式；框架自己在切换屏幕、按 Back、移动焦点时调用这三个方法。组件要改光标的形状，用 `CursorShaper`。

## 主题

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Theme struct` | 组件绘制用的颜色、字体和间距 | [主题](../concepts/theme.md#theme) |
| `Theme.Background`、`Foreground` | 背景色、前景色 | [主题](../concepts/theme.md#theme) |
| `Theme.Accent` | 屏幕不能反色时光标的颜色 | [反色光标](../guides/hardware.md#inverter) |
| `Theme.Muted` | 提示和禁用文字 | [主题](../concepts/theme.md#theme) |
| `Theme.Font`、`SmallFont` | 正文字体、小字体；`Font` 的行高决定 `Unit()` | [单位](../concepts/theme.md#unit) |
| `Theme.Radius` | 带边框控件的圆角，单位 | [主题](../concepts/theme.md#theme) |
| `Theme.RowPadding` | 列表行、输入框、标题的文字内边距 | [主题](../concepts/theme.md#theme) |
| `Theme.ButtonPadding` | 按钮文字和边框之间的距离 | [主题](../concepts/theme.md#theme) |
| `Theme.Gap` | 抽屉、对话框、带内边距的列表里控件之间的间距 | [主题](../concepts/theme.md#theme) |
| `Theme.PanelPadding` | 控件和外框之间的距离：每个弹窗、每个带边框的面板 | [主题](../concepts/theme.md#theme) |
| `var CurrentTheme = &Theme{...}` | 当前主题，默认黑底白字，适合 OLED；组件在构造时读取 | [主题](../concepts/theme.md#theme) |
