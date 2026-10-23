# 布局

每个节点的位置和大小由布局算出来：父节点给出一块区域，子节点按自己的尺寸、锚点和边距在里面占一块。布局用 `float32` 计算，弹簧推动的位置可以停在两个像素之间，画的时候才取整。

```go
import "github.com/chhongzh/lilt/ui"
```

有三种容器：`ui.Stack` 按锚点摆放子节点，`ui.Flex` 把子节点排成一行或一列，`ui.Grid` 按列排成表格。组件也是节点，大多数组件内部就是这三种容器之一。

## 尺寸 { #size }

宽和高各是一个 `ui.Size`，三种取值：

| 写法 | 含义 |
| --- | --- |
| `ui.Px(v)` | 固定 `v` 像素，不随[单位](theme.md#unit)缩放 |
| `ui.Rel(f)` | 父节点内容区（去掉内边距和自己的外边距之后）的 `f` 倍 |
| `ui.Auto` | 默认值，按内容定 |

`Auto` 的内容尺寸这样算：实现了 `ui.Measurer` 的叶子（如 `Label`）报告文字大小；实现了 `ui.Layouter` 的容器（`Flex`、`Grid`）汇总子节点；其他节点取所有可见子节点（连同外边距和位置偏移）的外接范围。最后再加上内边距。

| 选项 | 作用 |
| --- | --- |
| `ui.Width(s)`、`ui.Height(s)` | 设一个方向 |
| `ui.Fill()` | 两个方向都是 `Rel(1)` |
| `ui.FillWidth()`、`ui.FillHeight()` | 一个方向 `Rel(1)` |

对应的方法是 `SetWidth`、`SetHeight`、`Width()`、`Height()`。

滚动容器沿滚动方向给子节点的空间是 `ui.Unbounded`（2^24），无穷大的几分之几没有意义，所以这个方向上的 `Rel` 按 `Auto` 处理。把 `List` 放进 `List` 时就会遇到，见[嵌套滚动](../widgets/scroll.md#nested)。

## 锚点与原点 { #anchor }

在 `Stack` 和其他不自己排布子节点的容器里，子节点靠两个点定位：

- **锚点**（`Anchor`）：父节点内容区里的一点；
- **原点**（`Origin`）：子节点自己身上的一点，它被放到锚点上。

两者都是 `ui.Vec2`，`{0, 0}` 是左上角，`{1, 1}` 是右下角。`ui.Align(v)` 把锚点和原点设成同一个值，这是贴边、贴角最常用的写法：`ui.Align(ui.BottomRight)` 让节点的右下角对准父节点的右下角。

![按钮让标签在六个位置之间滑动](../assets/demos/layout-anchor.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_layout.go:layout-anchor"
```

预定义的九个点：

| | 左 | 中 | 右 |
| --- | --- | --- | --- |
| 上 | `ui.TopLeft` | `ui.TopCenter` | `ui.TopRight` |
| 中 | `ui.CenterLeft` | `ui.Center` | `ui.CenterRight` |
| 下 | `ui.BottomLeft` | `ui.BottomCenter` | `ui.BottomRight` |

锚点和原点可以不同：`ui.Anchor(ui.BottomCenter), ui.Origin(ui.TopCenter)` 把节点挂在父节点下边缘的外面。`ui.Position(x, y)` 在此基础上再平移 `x`、`y` **像素**。

用到的选项和方法：`ui.Anchor`、`ui.Origin`、`ui.Align`、`ui.Position`，`SetAnchor`、`SetOrigin`、`SetPosition` 及对应的读取方法。

!!! note "什么时候锚点不起作用"
    `Flex` 和 `Grid` 自己排布子节点，忽略子节点的锚点、原点和位置。`ui.Push` 推入的页面总是铺满整屏，页面根节点的锚点也不起作用，要在页面里放一个 `Stack` 再对齐。

## 外边距与内边距

外边距（`Margin`）是节点外面留的空，内边距（`Padding`）是节点边框和内容之间留的空。两者都是 `ui.Insets`，以**单位**计，绘制时乘 `ui.Unit()` 变成像素，见[主题与单位](theme.md#unit)。

| 写法 | 结果 |
| --- | --- |
| `ui.All(v)` | 四边都是 `v` |
| `ui.XY(x, y)` | 左右 `x`，上下 `y` |
| `ui.LTRB(l, t, r, b)` | 分别指定 |

`Insets` 还有 `Horizontal()`（左加右）、`Vertical()`（上加下）和 `Scale(f)`（每边乘 `f`）。节点的 `Px` 和 `Rel` 尺寸包含内边距，`Auto` 尺寸是内容加内边距。

自定义组件画内容时用 `ContentRect()`，它返回去掉内边距的区域，已经换算成像素，坐标相对节点自己。

行、按钮、弹窗共用的内边距放在主题里（`RowPadding`、`ButtonPadding`、`PanelPadding`），自己排版时也从主题取，界面才统一。

## Stack { #stack }

`ui.NewStack(opts...)` 是最简单的容器，子节点各自按[锚点](#anchor)定位，互相重叠，后加的画在上面。它适合角标、浮在内容上的提示，或者让一个节点在一块区域里居中。

## Flex { #flex }

`ui.NewRow` 把子节点从左到右排，`ui.NewColumn` 从上到下排。沿排列方向叫主轴，另一个方向叫交叉轴。

![切换主轴对齐，三个方块滑到新位置](../assets/demos/layout-flex.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_layout.go:layout-flex"
```

子节点之间的间距用 `ui.Gap(g)` 或 `SetGap(g)` 设，单位是单位。主轴上的空余空间按 `MainAlign` 分配：

| 值 | 效果 |
| --- | --- |
| `ui.MainStart` | 默认，靠起点 |
| `ui.MainCenter` | 居中 |
| `ui.MainEnd` | 靠终点 |
| `ui.MainSpaceBetween` | 首尾贴边，空余空间平分到间隙里 |

交叉轴上按 `CrossAlign` 放：

| 值 | 效果 |
| --- | --- |
| `ui.CrossStretch` | 默认，交叉轴尺寸是 `Auto` 的子节点撑满；尺寸固定的子节点靠起点 |
| `ui.CrossStart` | 靠起点 |
| `ui.CrossCenter` | 居中 |
| `ui.CrossEnd` | 靠终点 |

选项是 `ui.Main(m)`、`ui.Cross(c)`，方法是 `SetMainAlign`、`SetCrossAlign`。`Axis()` 返回方向（`ui.Horizontal` 或 `ui.Vertical`），`SetAxis` 可以在运行时把行改成列。隐藏的子节点不参与排列，也不占间距。

### 分配剩余空间

`ui.Grow(g)` 让子节点按比例分走主轴上的剩余空间：一个 `Grow(1)` 拿走全部，两个 `Grow(1)` 平分，`Grow(2)` 和 `Grow(1)` 按二比一分。

![打开 Grow 的方块撑开，两个一起打开时平分](../assets/demos/layout-grow.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_layout.go:layout-grow"
```

只要有一个子节点在长，`MainAlign` 就不起作用，因为没有剩余空间可分。

内容超出主轴长度时：

- **没有子节点在长**：每个子节点保持自己的尺寸，排到容器外面去。容器设了[裁剪](#clip)就被切掉，否则会画到容器外面。
- **有子节点在长**：超出的部分从会长的子节点身上扣，最多扣到 0。一行里标题 `Grow(1)`、右边是开关，标题太长时缩短的是标题，开关不会被挤出去；`ListItem` 就是这样排的，被截短的标题在这一行聚焦时由[跑马灯](../widgets/label.md)滚动显示。

!!! warning "在 Flex 里填满剩余空间用 Grow，不用 Rel"
    `Rel` 相对的是整个容器的内容区，不是兄弟节点剩下的空间。一列里放一个标题，再放一个 `Height(Rel(1))` 的列表，两者加起来比容器高。给列表 `ui.Grow(1)` 才能让它占满标题下面的部分。（`widget.List` 默认是 `Rel(1)`，加上 `Grow(1)` 之后超出的部分从它身上扣，所以也能用。）

## Grid { #grid }

`ui.NewGrid` 把子节点一行一行填进若干列。列宽用 `ui.Columns(...)` 或 `SetColumns(...)` 给出，每列是一个 `ui.Track`：

| 写法 | 含义 |
| --- | --- |
| `ui.Fixed(px)` | 固定宽度，**像素** |
| `ui.Fr(n)` | 固定列和间距分完之后，剩下的宽度按 `n` 的比例分 |

![三行灯的开关，光标在格子之间上下左右移动](../assets/demos/layout-grid.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_layout.go:layout-grid"
```

- 不设列时只有一列 `Fr(1)`。
- 每行和这一行最高的子节点一样高。
- 宽是 `Auto` 的子节点撑满列宽，高是 `Auto` 的撑满行高，所以按钮的宽度跟着列走。
- `ui.Gap` 同时是行距和列距，以单位计。
- 隐藏的子节点不占格子，后面的子节点往前补。

方向键在格子之间按位置找最近的节点，不需要额外设置，见[输入与焦点](input.md#focus)。

## 显示与隐藏 { #visible }

`SetVisible(false)` 或构造时的 `ui.Hidden()` 隐藏一个节点。隐藏的节点不占空间、不画、不调用 `Update`，也不能持有焦点：焦点在里面时会移到最近的可见节点。

![关掉 Wi-Fi 时 Network 行消失，下面一行滑上来](../assets/demos/layout-visible.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_layout.go:layout-visible"
```

`network` 隐藏之后，`Bluetooth` 这一行因为设了 `ui.AnimateLayout()` 而滑上去，而不是直接跳上去。回调里用 `to` 而不是去读复选框，因为 `OnChange` 已经把新旧值都传进来了。

## 位置变化时滑动 { #animate-layout }

默认情况下，布局变了，节点下一帧就出现在新位置。`ui.AnimateLayout()`（或 `SetAnimateLayout(true)`）让它滑过去：重新布局时，新旧位置的差被加进节点的 `OffsetX()`、`OffsetY()`，这两个动画值再用弹簧回到 0。

- 只有位置会滑，宽高立即变成新值。
- 偏移用的是 `OffsetX()`、`OffsetY()` 自己的弹簧，默认是 [`anim.Default`](animation.md#presets)。有些组件为了抖动换了弹簧，例如 `Button` 的 `OffsetX` 用 `anim.Bouncy`，它滑动时也会带一点回弹。
- 子节点相对父节点绘制，父节点滑动时子节点跟着走，不需要各自设置。
- 第一次布局没有旧位置，不会滑。节点被 `ui.Remove` 后再加回树里，也算第一次。

前面[锚点](#anchor)、[Flex](#flex)、[显示与隐藏](#visible)的动图都用了它。

## 裁剪、背景与边框 { #clip }

| 选项 | 方法 | 作用 |
| --- | --- | --- |
| `ui.Clip()` | `SetClip(bool)` | 节点和子节点只画在自己的矩形里，矩形外的部分也不响应触摸 |
| `ui.Background(c)` | `SetBackground(c)` | 用颜色 `c` 填满矩形 |
| `ui.Border(c)` | `SetBorder(c)`、`Bordered()` | 描一像素宽的边框 |
| `ui.Radius(r)` | `SetRadius(r)` | 背景和边框的圆角，以单位计 |

背景和边框在节点自己的 `Draw` 之前画，不需要组件自己处理。聚焦的节点被裁剪祖先遮住一部分时，光标只框住看得见的那部分。

`Button`、`ListItem` 都开了裁剪，里面的文字滑动或滚动时不会画到边框外。

## 重新布局 { #relayout }

```go
func Relayout()   // 下一帧重新布局并重画
func Invalidate() // 下一帧只重画
```

`Box` 和各组件的 Set 方法在值真的变了时会自己调用这两个函数：影响尺寸或位置的（`SetWidth`、`SetPadding`、`SetText`、增删节点）调用 `Relayout`，只影响外观的（`SetBackground`、`SetClip`）调用 `Invalidate`。用现成组件时不需要手动调用；自定义组件改了影响测量结果的字段要调用 `Relayout`，只改了画法调用 `Invalidate`。

`Relayout` 不区分是谁变了，下一帧从根开始**整棵树重新布局**。单片机屏幕上的树很小，追踪哪些子树变脏要多写的代码比省下的时间更多，所以 lilt 没有做脏标记。代价是同一帧里调用多少次 `Relayout` 都只排一次，但只要调用一次，就是整棵树。

一次布局里同一个节点可能被测量不止一次（`Flex` 先汇总尺寸再摆放），所以 `Measure` 要快、不要分配内存。自己实现 `ui.Layouter` 时可以用导出的 `ui.Measure(n, max)`、`ui.Arrange(n, rect)` 和 `ui.ArrangeFree(n, area)` 测量和摆放子节点，见[自定义组件](../guides/custom-widget.md#interfaces)。
