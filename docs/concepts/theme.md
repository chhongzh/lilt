# 主题与缩放

组件的颜色、字体、圆角和内边距都来自一个全局主题 `ui.CurrentTheme`。所有长度都是以 128×64 为基准的[逻辑像素](#scale)，换到更大的屏幕时由框架按比例放大。

```go
import "github.com/chhongzh/lilt/ui"
```

## Theme { #theme }

`ui.CurrentTheme` 是一个 `*ui.Theme`，默认是适合 OLED 的黑底白字：

| 字段 | 默认值 | 用在哪里 |
| --- | --- | --- |
| `Background` | `gfx.Black` | 清屏、页面和弹窗的底色 |
| `Foreground` | `gfx.White` | 文字、边框、勾和加载点 |
| `Accent` | `gfx.White` | 屏幕不能反色时，光标的描边颜色 |
| `Muted` | 灰色 `{128, 128, 128}` | 提示文字、抽屉标题 |
| `Font` | `notosans.Notosans12pt` | 正文字体，大屏上要换成[更大的](#fonts) |
| `SmallFont` | `proggy.TinySZ8pt7b` | 次要文字：`Select` 的当前值、对话框的消息、软键盘 |
| `Radius` | `3` | 按钮、输入框、弹窗和光标的圆角，逻辑像素 |
| `RowPadding` | `ui.XY(4, 1)` | 列表行、输入框、标题的文字内边距 |
| `ButtonPadding` | `ui.XY(6, 2)` | 按钮文字和边框之间 |
| `Gap` | `2` | 对话框等处上下排列的控件之间 |
| `PanelPadding` | `ui.All(2)` | 弹窗的边框和内容之间，以及自己画边框的面板 |

`Font`、`SmallFont` 是 `*gfx.Font`，也就是 tinyfont 的字体，`tinygo.org/x/tinyfont` 下的字体都能直接用。颜色是 `gfx.Color`（`color.RGBA`）。

!!! warning "先改主题，再建组件"
    组件大多在**构造时**读主题：`Button` 构造时把 `Foreground` 设成边框色、把 `ButtonPadding` 设成内边距，之后再改主题，已经建好的按钮不会跟着变。只有一部分值是每帧画的时候读的（清屏的 `Background`、文字的 `Foreground`、没单独设字体的 `Label` 用的 `Font`、光标的 `Radius` 和 `Accent`），运行中改主题会得到一半新一半旧的界面。所以主题要在建树之前改好。

## 逻辑像素与缩放 { #scale }

```go
const BaseW, BaseH = 128, 64

func SetZoom(f float32)           // 在 lilt.Start 之前调用
func ScaleFor(w, h int16) float32 // 某块屏幕会得到的缩放
func Scale() float32              // 当前屏幕的缩放：一个逻辑像素等于多少物理像素
func Dp(v float32) int16          // 逻辑像素换算成物理像素，四舍五入
```

界面按 128×64 的屏幕来写，这也是 SH1106 的大小。布局里的所有长度都是**逻辑像素**：`ui.Px` 尺寸、`Padding`、`Margin`、`Gap`、`Radius`、`ui.Position` 的平移、`Grid` 的 `ui.Fixed` 列宽，以及主题里的 `Radius`、`RowPadding`、`ButtonPadding`、`Gap`、`PanelPadding`。`ui.ScreenSize()` 返回的也是逻辑像素。

`lilt.Start` 根据屏幕的物理尺寸算出缩放：

```
Scale = min(物理宽 / 128, 物理高 / 64) × zoom
```

取宽高两个比例中较小的一个，所以长宽比和 128×64 不同的屏幕不会把界面拉变形，多出来的那一边是更多的逻辑像素。例如 320×172 的 JD9853 横屏，`min(2.5, 2.69) = 2.5`，逻辑尺寸是 128×68.8。

`zoom` 默认是 1，用 `ui.SetZoom` 手动调整：大于 1 时所有东西画得更大、一屏放下的更少，小于 1 时反过来。它必须在 `ScaleFor` 和 `lilt.Start` 之前调用。

| 屏幕 | zoom | `Scale` | 逻辑尺寸 |
| --- | --- | --- | --- |
| 128×64（SH1106） | 1 | 1 | 128×64 |
| 192×96 | 1 | 1.5 | 128×64 |
| 320×172（JD9853 横屏） | 1 | 2.5 | 128×68.8 |
| 320×172 | 1.25 | 3.125 | 102.4×约 55 |

### 字体 { #fonts }

字体是位图，只能按自己的像素大小画，不跟着缩放。所以在放大的屏幕上要换成大约 `Scale` 倍大的字体。`ui.LineHeight(f)` 和 `ui.TextWidth(f, s)` 返回的是逻辑像素，组件按它们测量，换字体后布局自动跟着变。

字体要在建树之前换好，而那时还没有 `lilt.Start`，所以用 `ui.ScaleFor(w, h)` 预先算出缩放。示例里的 `gallery.FitFonts` 就是这样按缩放挑 FreeSans 的字号：

```go
ui.SetZoom(1)          // 要调 zoom 的话先调
gallery.FitFonts(w, h) // 按 ui.ScaleFor(w, h) 换主题字体
lilt.Start(screen, gallery.Home())
```

几种字体在逻辑像素下的行高（物理行高除以 `Scale`）：

| 字体 | 物理行高 | `Scale` 1 | `Scale` 1.5 | `Scale` 2.5 |
| --- | --- | --- | --- | --- |
| `notosans.Notosans12pt`（默认） | 15 | 15 | 10 | 6 |
| `freesans.Regular9pt7b` | 22 | 22 | 约 14.7 | 8.8 |
| `freesans.Regular12pt7b` | 29 | 29 | 约 19.3 | 11.6 |
| `freesans.Regular18pt7b` | 42 | 42 | 28 | 16.8 |

同一份代码，左边是默认主题，右边是 192×96 的屏幕、字体换成 FreeSans 9pt。缩放是 1.5，行高、内边距、间距、圆角一起放大，各部分的比例不变：

=== "默认字体，128×64"

    ![默认主题下的设置页](../assets/demos/theme-default.gif){ .demo }

=== "FreeSans 9pt，192×96"

    ![放大的屏幕上的同一个设置页](../assets/demos/theme-large.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_theme.go:theme-settings"
```

```go
--8<-- "cmd/docgif/demo_theme.go:theme-large"
```

`settings` 里的长度都是逻辑像素：列表的内边距是 `ui.XY(2, 1)`，行间距取主题的 `Gap`，行和按钮的内边距由组件从主题里取。换屏幕时这些都不用改，只换字体。

自定义组件的 `Measure` 和布局用逻辑像素，`Draw` 收到的画布是物理像素：矩形用 `Pixels()` 转，单个长度用 `ui.Dp` 转，见[自定义组件](../guides/custom-widget.md#measure)。

## 共用的内边距

主题里的四个内边距让不同组件对齐到同一套间距上，自己排版时也从这里取：

| 字段 | 谁在用 | 自己排版时 |
| --- | --- | --- |
| `RowPadding` | `ListItem` 及其派生的行、`Input`、抽屉标题 | 自己做的一行文字，比如表头 |
| `ButtonPadding` | `Button` | 自己画边框、里面放一个词的控件 |
| `Gap` | 对话框里的按钮行和消息 | 竖着排的几个带边框控件之间 |
| `PanelPadding` | 层栈给每个弹窗留的边 | 自己画边框的面板 |

一个带边框的面板：

```go
t := ui.CurrentTheme
panel := ui.NewColumn(ui.Border(t.Foreground), ui.Radius(t.Radius),
	ui.Padding(t.PanelPadding), ui.Gap(t.Gap), ui.Children(
		widget.NewLabel("Battery"),
		widget.NewLabel("82%", widget.Font(t.SmallFont), widget.TextColor(t.Muted)),
	))
```

!!! note "弹窗不要自己加内边距"
    抽屉、对话框的边框由[层栈](layers.md#popups)画，`PanelPadding` 也由层栈留。放进弹窗的内容不要再加一圈内边距，否则边上会空出两倍的距离。

## 彩色屏幕

彩色屏换一套颜色即可，在建树之前改字段：

```go
t := ui.CurrentTheme
t.Background = gfx.Color{R: 250, G: 248, B: 242, A: 255}
t.Foreground = gfx.Color{R: 30, G: 30, B: 36, A: 255}
t.Accent = gfx.Color{R: 0, G: 110, B: 220, A: 255}
t.Muted = gfx.Color{R: 140, G: 140, B: 150, A: 255}
```

光标怎么画取决于屏幕：驱动实现了 `gfx.Inverter`（SH1106、silent 这类单色屏）时，光标把下面的像素反色；没实现时，光标画成 `Accent` 颜色、主题圆角的描边。`gfx.FromDisplayer` 包装出来的屏幕不能反色，所以接彩色屏时记得把 `Accent` 设成一个显眼的颜色。

单色屏只有亮和灭，颜色由 `gfx.IsLit` 判定：按亮度加权（R×3 + G×6 + B）不低于 128×10 的算亮。默认的 `Muted` 灰色刚好算亮，所以在单色屏上和 `Foreground` 看起来一样。

## 光标的弹簧

光标的手感由 `ui` 包里的五个弹簧决定，都是 `*anim.Spring`（参数含义见[弹簧](animation.md#spring)）：

| 变量 | 默认值 | 管什么 |
| --- | --- | --- |
| `ui.CursorLead` | `Response: 0.22, Ease: true` | 移动时朝前的那条边 |
| `ui.CursorTrail` | `Response: 0.52, Ease: true` | 移动时落后的那条边 |
| `ui.CursorRelease` | `Response: 0.57, Damping: 0.4` | 松开按键后光标弹回原样，带晃动 |
| `ui.CursorLeap` | `Response: 0.63, Damping: 0.45` | 切换页面时光标跳起再落下 |
| `ui.CursorLean` | `Response: 0.35, Damping: 0.7` | 按住方向键时光标朝那边探出 |

光标的四条边各自运动：朝前的边用较快的 `CursorLead`，另一边用较慢的 `CursorTrail`，所以光标先朝目标拉长，到了再收拢。按下时的收缩用的是 `anim.Interactive`，不在这张表里。

想调手感就改字段，例如 `ui.CursorTrail.Response = 0.4` 让拖尾短一些。改字段对所有用到它的地方立即生效；换掉整个指针（`ui.CursorLeap = &anim.Spring{...}`）则不保险，`CursorLeap` 和 `CursorLean` 在启动时就绑定好了，换指针不起作用。
