# 自定义组件

标准组件不够用时，自己写一个。组件就是一个嵌入了 `ui.Box` 的结构体，再按需要实现几个接口：要画东西就实现 `Draw`，要接按键就实现 `HandleKey`。框架每帧用类型断言找这些接口，没有注册表，也没有基类方法要覆盖。

这一页从头写一个滑块 `Slider`，再写一个电池图标 `Battery`，把用到的接口、画布和动画逐个讲清楚。

## 一个滑块

![左右键拖动的滑块](../assets/demos/custom-slider.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-slider"
```

左右键移动滑块，按住连发时一格一格走；上下键照常移动焦点。光标不罩住整行，而是只套在圆钮上，跟着它走。下面几节就是 `Slider` 的全部代码。

## 类型与构造函数 { #type }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-slider-type"
```

- **嵌入 `ui.Box`**。`ui.Node` 接口只有一个未导出的方法，只有嵌入了 `ui.Box` 的类型才满足它。`Box` 带着父子链接、尺寸、内边距和动画偏移；要把子节点排成一行或一列，就改为嵌入 `ui.Flex`（`ListItem` 就是这样）。
- **构造函数收 `...ui.Option`，`ui.Apply` 放最后**。先设好默认值，再应用选项，调用方传的 `ui.Width`、`ui.Padding` 才能覆盖默认值。所有标准组件都是这个顺序，见[选项](../concepts/tree.md#options)。
- **自己的选项就是 `func(ui.Node)`**。`OnSlide` 断言出 `*Slider` 再写字段，用在别的节点上会 panic，和 `widget` 包的选项一样。
- **默认内边距从主题取**。这里用 `RowPadding`，所以它和列表行左右对齐，换主题时一起变，见[主题与单位](../concepts/theme.md#theme)。
- **状态走 Set/Get**。`SetValue` 是代码里改值的入口，按约定它不调用回调。

改了状态之后要告诉框架：影响尺寸的调用 `ui.Relayout()`，只影响外观的调用 `ui.Invalidate()`。`SetValue` 两个都没调，因为它只改了 `anim.Value` 的目标，而帧循环在任何 `anim.Value` 还在动时都会重绘，见[动画](#animation)。

!!! warning "TinyGo 下不要用这些方法名"
    嵌入 `ui.Box` 或 `ui.Flex` 的类型，如果声明了和它们未导出方法同名的方法：`box`、`set`、`visible`、`split`、`join`、`margins`、`crossAuto`，TinyGo 0.42 和 0.43 在把它转成接口时会崩溃（`index out of range [0] with length 0`）。主机上的 Go 编译正常，所以要到 `tinygo build` 才会发现。给自己的方法换个名字，例如 `Slider` 的 `move`、`knobX`。

## 能力接口 { #interfaces }

`ui` 包定义了一组只有一两个方法的接口，组件实现哪个，就获得哪种能力。都在 UI goroutine 上调用。

| 接口 | 方法 | 什么时候调用 |
| --- | --- | --- |
| `Painter` | `Draw(c gfx.Canvas)` | 绘制自己，在子节点之前 |
| `OverPainter` | `DrawOver(c gfx.Canvas)` | 在子节点之后再画一层，例如滚动条、加载点 |
| `Updater` | `Update(dt float32)` | 每帧一次，节点在树上且可见时；`dt` 是秒 |
| `Measurer` | `Measure(max Vec2) Vec2` | 宽或高是 `Auto` 时问内容尺寸，不含内边距 |
| `Layouter` | `MeasureChildren(max Vec2) Vec2`<br>`ArrangeChildren(content Rect)` | 自己测量、摆放子节点，`Flex`、`Grid` 都是 |
| `KeyHandler` | `HandleKey(e *KeyEvent) bool` | 按键从焦点节点往根冒泡 |
| `KeyPreviewer` | `PreviewKey(e *KeyEvent) bool` | 按键从根往焦点节点下行，先于 `HandleKey` |
| `PointerHandler` | `HandlePointer(e *PointerEvent) bool` | 触摸从指针下的节点往根冒泡 |
| `Focuser` | `FocusChanged(focused bool)` | 自己获得或失去焦点 |
| `Revealer` | `Reveal(d Node)` | 某个后代获得焦点，滚动容器借此把它滚进视野 |
| `ArrowScroller` | `ArrowScroll(e *KeyEvent, next Node) bool` | 焦点移动之前，焦点节点的祖先先拿到方向键 |
| `CursorShaper` | `CursorShape(r Rect) Rect` | 获得焦点时每帧一次，改光标的矩形 |
| `ContentHost` | `Content() Node` | `ui.Add`、`ui.RemoveChildren` 改去操作返回的内部节点 |
| `Leaf` | `Leaf()` | 标记不接受子节点，`ui.Attach` 到它上面会 panic |
| `ChildRemover` | `ChildRemoved(child Node)` | 子节点被移除，包括子节点移除自己 |

几条规则：

- **只有可聚焦的节点能拿到焦点和按键**。在构造函数里 `SetFocusable(true)`，或者让调用方传 `ui.Focusable()`。
- **`Draw` 和 `DrawOver` 收到的画布已经平移到节点左上角**（含 `OffsetX`、`OffsetY` 的动画偏移）；只有 `SetClip(true)` 的节点才会被裁剪到自己的矩形里。`SetBackground`、`SetBorder` 的底色和边框由 `Box` 先画好，`Draw` 画在它们上面。
- **`Measure` 的返回值是像素**，`max` 是内边距以内还剩的空间。宽高都是固定值或比例时不会调用它。
- **`HandleKey` 在 `Press` 上返回 `true`，就认领了这个键**，之后的 `Repeat` 和 `Release` 直接送到这个节点，不再冒泡。返回 `false` 交给祖先，最后由框架处理：方向键移动焦点，Back 关闭当前页面或弹窗。分发路径见[按键分发](../concepts/input.md#dispatch)。
- **`Updater` 只在需要逐帧推进时实现**，例如计时器、跑马灯。只是让一个数值动起来，用 `anim.Value` 就够了，不需要 `Update`。

## 测量与绘制 { #measure }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-slider-draw"
```

`Slider` 的宽度要多少给多少，高度刚好放下圆钮。几个容易弄混的地方：

- **单位和像素**。`Padding`、`Margin`、`Gap`、`Radius` 用的是单位，`Measure`、`Draw` 里用的是像素。`ui.Unit()` 是一个单位等于多少像素，默认字体下是 1，换大字体时变大。圆钮半径写成 `3 * ui.Unit()`，所以字体放大时它跟着放大，见[单位](../concepts/theme.md#unit)。
- **`ContentRect()` 已经扣掉了内边距**，是节点自己坐标系里的像素，类型是 `ui.Rect`（float32）。画图前用 `Pixels()` 转成 `gfx.Rect`，它对四条边分别取整，相邻的矩形取整后仍然相邻。
- **`knobX` 用 `s.knob.Get()`**，即圆钮此刻的位置，而不是目标位置，这样画出来的就是动画中的样子。
- **颜色取自 `ui.CurrentTheme`**，不要写死。单色屏上 `Background` 和 `Foreground` 就是灭和亮；先用 `Background` 画一圈稍大的圆，再画圆钮，圆钮和填充条之间就留出一道缝。

## 画布 { #canvas }

`gfx.Canvas` 是屏幕的一个平移、裁剪过的视图。它是值类型：`Translate` 和 `Clip` 返回新的画布，原来的不变，所以嵌套绘制不用保存和恢复状态。坐标全是 `int16` 像素，以画布原点为准。

| 方法 | 作用 |
| --- | --- |
| `Translate(dx, dy)` | 原点移动 `dx`、`dy`，返回新画布 |
| `Clip(r)` | 只允许画在 `r` 以内，返回新画布 |
| `Absolute()` | 原点回到屏幕左上角，保留裁剪 |
| `Origin()`、`Bounds()` | 原点的屏幕坐标；裁剪区在本地坐标里的矩形 |
| `Visible(r)` | `r` 是否有一部分在裁剪区内，用来跳过画不出来的东西 |
| `Pixel(x, y, col)` | 一个点 |
| `HLine`、`VLine` | 一像素宽的横线、竖线 |
| `Line(x0, y0, x1, y1, col)` | 任意方向的直线（Bresenham），画勾、叉这类小图标 |
| `FillRect`、`StrokeRect` | 实心、空心矩形 |
| `FillRoundRect`、`StrokeRoundRect` | 圆角矩形，`radius` 是像素 |
| `FillCircle`、`StrokeCircle` | 圆心加半径 |
| `Text(f, x, y, s, col)` | 一行字，`y` 是基线，返回最后一个字形之后的 `x` |
| `InvertRect`、`InvertRoundRect` | 把区域里的像素反色，屏幕实现了 `gfx.Inverter` 才有效 |
| `CanInvert()` | 反色是否有效 |

文字按基线定位。要在一个盒子里垂直居中，照 `Label` 的算法：

```go
y := box.Y + (box.H-gfx.LineHeight(f))/2 + gfx.Ascent(f)
c.Text(f, x, y, s, col)
```

`gfx.TextWidth(f, s)` 量文字宽度，可以拿来右对齐或在 `Measure` 里报尺寸。字体是 tinyfont 的位图字体，主题里有 `Font` 和 `SmallFont` 两种。

反色只在实现了 `gfx.Inverter` 的屏幕上生效，其他屏幕上 `InvertRect` 什么也不做。需要反色效果的组件先问 `CanInvert()`，不能反色时换一种画法，光标就是这样：能反色时把焦点节点反色，不能时画一个 `Accent` 色的轮廓。哪些屏幕能反色见[接入硬件](hardware.md#inverter)。

## 按键 { #keys }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-slider-keys"
```

这段照着框架的[交互约定](../concepts/input.md#activation)写：

- **`Press` 只给反馈**。圆钮朝按键方向探出 0.3 格，表示“收到了”，值不变。
- **`e.Fires()` 时才动**。它在没被取消、前面也没有连发的 `Release` 上为真，按住时的每个 `Repeat` 上也为真。按住右键就一格一格走，松手时不会多走一格。
- **其余情况把圆钮放回原位**。被取消的松开（触摸变成拖动）和连发之后的松开都走这里。
- **上下键返回 `false`**，交给框架移动焦点，所以两个滑块之间能上下切换。
- **回调用 `go` 调用**。用户回调永远不在 UI goroutine 上跑，慢回调不会卡住动画；回调里改树要包在 `ui.Do` 里，就像使用示例里更新 `level` 那样，见[线程模型](../concepts/threading.md#do)。

`Slider` 认领左右键，所以焦点在滑块上时左右键不再移动焦点。如果滑块放在横向的 `Scroll` 里，这也意味着它会挡住横向滚动，设计布局时要考虑到。

## 改变光标 { #cursorshaper }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-slider-cursor"
```

焦点节点实现了 `ui.CursorShaper` 时，光标每帧把节点的矩形交给 `CursorShape`，用它返回的矩形代替。参数和返回值都是屏幕像素坐标；参数已经被开了裁剪的祖先裁过，滚出一半的行只给可见的那部分。

`Slider` 返回圆钮周围的一个小方块，光标就只罩住圆钮。因为用的是 `knobX()` 此刻的位置，圆钮动，光标的目标也跟着动；光标四条边各有自己的弹簧（`ui.CursorLead`、`ui.CursorTrail`），会朝运动方向先拉长再收拢。[`Button`](../widgets/button.md) 的长按确认也是这个接口：按住时从左边把光标一点点切掉。

光标的按下收缩、换焦点时的拉伸、跨屏幕的跳跃都在 `CursorShape` 的结果上叠加，组件不用自己处理，见[光标](../concepts/input.md#cursor)。

## 动画 { #animation }

`Slider` 的圆钮是一个 `anim.Value`：一个朝目标弹过去的 float32。零值就能用，静止在 0，用 `anim.Default` 弹簧。

| 调用 | 作用 |
| --- | --- |
| `v.Set(x)` | 以弹簧动画移向 `x` |
| `v.Jump(x)` | 立即到 `x`，停止运动；构造函数里设初值用它 |
| `v.Get()` | 此刻的值，画图时用 |
| `v.Target()` | 正在去的值 |
| `v.SetSpring(s)` | 换弹簧，例如 `anim.Snappy` |
| `v.SetPrecision(p)` | 离目标多近算停下；默认 0.05 适合像素，0 到 1 的进度要更小 |
| `v.Bump(px)`、`v.Kick(vel)` | 给一个冲量，弹出去再回来 |

`Set` 让值加入调度器的活动列表，帧循环每帧推进所有活动值，只要有一个还在动就重绘整屏，停下后自动退出列表。所以组件只要在 `Draw` 里读 `Get()`，不用实现 `Updater`，也不用调用 `ui.Invalidate()`。`Value` 被调度器用指针引用，开始动之后不能再复制；把它作为字段放在通过指针使用的组件结构体里就没问题。弹簧参数和预设见[动画](../concepts/animation.md#value)。

要让整个节点动，不必自己存偏移：`Box` 自带 `OffsetX()` 和 `OffsetY()` 两个 `anim.Value`，绘制时叠加在布局位置上，子节点跟着一起动。标准组件失败时左右抖一下，就是 `OffsetX().Kick(...)`。

## 事件回调与 `widget.Control` { #control }

标准组件的 `OnClick`、`OnChange`、加载点和失败抖动都来自嵌入的 `widget.Control`。但 `Control` 里保存回调的字段、启动回调的 `run`、推进加载计时的 `tick` 都是未导出的，`widget` 包以外的代码碰不到。在你自己的包里嵌入 `widget.Control`：

| 能做 | 不能做 |
| --- | --- |
| 类型上有 `OnClick` 方法，`widget.OnClick(fn)` 选项也能用在它上面，回调会被存下来 | 调用存下来的回调：字段和 `run` 都未导出，存了也不会执行 |
| | 用 `OnChange`：`Control` 不知道你的值类型，`widget.OnChange(fn)` 和 `OnChange` 方法都会 panic |
| | 用 `Busy()`、`Loading()`：没有回调在跑，它们永远是 `false` |

所以在别的包里写组件，有两条路。

**一、交给标准组件处理激活**。把自己的节点放进 `ListItem` 的尾部（或放在一个 `Button` 旁边），点击、加载点、失败抖动都由 `ListItem` 的 `Control` 负责，自己的组件只管显示：

![电池图标放在列表行里](../assets/demos/custom-battery.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_custom.go:custom-battery"
```

`Battery` 只实现了 `Measure`、`Draw` 和 `Leaf`，没有一行按键代码：

```go
--8<-- "cmd/docgif/demo_custom.go:custom-battery-type"
```

回调超过 150ms 时，`ListItem` 把尾部的电池滑走、换成加载点；回调在 `ui.Do` 里改了电量，返回后电池滑回来，填充条再弹到新电量。

**二、自己管回调**。像 `Slider` 那样存一个函数字段，在 UI goroutine 上判断好该不该调用，然后 `go` 出去。这样没有加载点，也不会在回调运行期间屏蔽激活；需要这些时，自己记一个 `busy` 标志，在回调结束时用 `ui.Do` 清掉它。

!!! note "不要把标准组件按值嵌入"
    `widget.NewListItem`、`widget.NewButton` 在构造时创建内部的子节点并把它们挂到自己身上，初始化方法也是未导出的。在别的包里按值嵌入 `widget.ListItem` 得到的是一个没有标题、不可聚焦的空行；复制 `*NewListItem(...)` 的结果则会让子节点的父指针指向原来那个对象。用组合，不要用嵌入。

## 下一步

- [API 速查：ui](../reference/ui.md)：所有接口和 `Box` 方法的签名。
- [接入硬件](hardware.md)：组件画到的屏幕从哪来，按键怎么进来。
