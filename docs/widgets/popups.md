# Drawer、Dialog 与 Toast

三种现成的弹窗：抽屉 `widget.Drawer` 从屏幕边缘滑进来，放一组操作或选项；对话框 `widget.Dialog` 落在屏幕中间，请用户确认一件事；`widget.Toast` 在底部冒出一行字，两秒后自己消失。它们都建立在层栈的 [`ui.Open`](../concepts/layers.md#popups) 上，边框和内边距由层栈画、由层栈留，弹窗内容自己不要再加。

```go
import "github.com/chhongzh/lilt/widget"
```

## 抽屉 { #drawer }

![从底部升起的抽屉，选一项后关闭并提示](../assets/demos/popups-drawer.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_popups.go:popups-drawer"
```

`NewDrawer(edge, title)` 建一个竖向排列的抽屉，标题是一行灰色小字，传空字符串就没有标题。用 `ui.Add` 往里放内容，标题下面依次排列；它本身是一个 [`ui.Flex`](../concepts/layout.md#flex)，`ui.Gap` 之类的选项也能用。`Open` 和 `Close` 都要在 UI goroutine 上调用，所以例子在回调里包了 `ui.Do`；列表项的回调里先 `Close` 再弹 `Toast`，两件事在同一个 `ui.Do` 里完成。

抽屉的大小由内容决定：

| 边 | 位置与大小 | 弹窗种类 |
| --- | --- | --- |
| `widget.Bottom` | 贴底，宽度占满，高度等于内容 | `ui.DrawerBottom` |
| `widget.Top` | 贴顶，宽度占满，高度等于内容 | `ui.DrawerTop` |
| `widget.Left` | 贴左，高度占满，宽度等于内容 | `ui.DrawerLeft` |
| `widget.Right` | 贴右，高度占满，宽度等于内容 | `ui.DrawerRight` |

抽屉的边框最多和屏幕一样大，内容再多也不会自动滚动。选项很多时放一个 [`List`](list.md#list) 并给它固定高度，[`Select`](select.md) 就是这样做的：它把列表高度限制在屏幕高度的 80% 减去 `TitleHeight()`（标题连同内边距的高度）和上边的内边距。

侧边抽屉适合做导航菜单：

![从左边滑出的菜单，按住 Back 再松开关闭](../assets/demos/popups-side.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_popups.go:popups-side"
```

抽屉可以先建好、多次打开，`ui.Do(menu.Open)` 直接把方法值交给 `ui.Do`。抽屉打开期间独占焦点和按键：打开时焦点移到它的第一个可聚焦节点，关闭后回到打开前的位置。按下 ++back++ 时抽屉朝来的方向让开一点，光标跟着歪过去，松开就关掉；这套返回手势对所有页面和弹窗都一样（见[返回](../concepts/layers.md#back)）。

抽屉贴着的那一边没有边框和内边距，边框从屏幕外绕过来，只露出里侧的两个圆角，圆角是主题 `Radius` 的两倍。进出用 `anim.Smooth`。

## 对话框 { #dialog }

![删除前确认，删除需要一点时间](../assets/demos/popups-dialog.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_popups.go:popups-dialog"
```

`NewDialog(title, message, buttons...)` 从上往下排：居中的标题、居中的一行小字说明、一排按钮。说明只有一行，放不下时一直滚动（[跑马灯](label.md#marquee)），传空字符串则不显示。按钮平分对话框的宽度，对话框本身占满屏幕宽度减去两边各 4 个单位。

对话框落进屏幕中间，第一个按钮获得焦点。每个按钮的 `OnClick` 照常执行，加载点和失败的抖动都和普通 [`Button`](button.md) 一样；回调返回 `nil` 后对话框自动关闭，返回 `error` 时对话框留着，用户可以重试或换一个按钮。没有设 `OnClick` 的按钮（例子里的 `Keep`）按下就关闭对话框。按 ++back++ 也会关闭。

所以回调里不需要自己关对话框，只管做事：例子里 `Delete` 的回调睡 700ms 模拟删除文件，期间按钮显示加载点，返回后对话框从上方退出，回调里弹出的 `Toast` 留在页面上。

`Title()` 和 `Message()` 返回标题和说明的 `Label`，可以在打开前后改文字。

## Toast { #toast }

![两个 Toast，后一个替换前一个](../assets/demos/popups-toast.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_popups.go:popups-toast"
```

`widget.Toast(text)` 在屏幕底部显示一行小字，2 秒后自动关闭。它不拿焦点、不收按键，所以显示期间用户照常操作页面，动图里光标一直可以上下移动。新的 Toast 会立刻替换还在显示的旧的，不会排队。文字太长时一直滚动。

`Toast` 会改动节点树，必须在 UI goroutine 上调用，回调里要包在 `ui.Do` 里。

## 与层栈的关系

三种组件都是对 `ui.Open(node, kind)` 的包装，`kind` 决定弹窗的位置、进出方向以及是否独占按键：

| `ui.Popup` | 位置 | 独占焦点和按键 | 对应组件 |
| --- | --- | --- | --- |
| `ui.DrawerBottom`、`ui.DrawerTop`、`ui.DrawerLeft`、`ui.DrawerRight` | 贴边 | 是 | `Drawer`、软键盘 |
| `ui.Dialog` | 屏幕中间 | 是 | `Dialog` |
| `ui.Toast` | 底部 | 否 | `Toast` |

要做别的弹窗，直接把任意节点交给 `ui.Open`，用 `ui.Close(node)` 关掉，`ui.Pop()` 关掉最上面的一层；[自定义输入法](input.md#text-input)里的抽屉就是例子。弹窗和 `ui.Push` 推入的页面共用一个层栈，规则见[页面与弹窗](../concepts/layers.md#push)。

## API

```go
type Edge uint8

const (
	Bottom Edge = iota
	Top
	Left
	Right
)

func NewDrawer(edge Edge, title string, opts ...ui.Option) *Drawer

func (d *Drawer) Open()                 // 滑入，在 UI goroutine 上调用
func (d *Drawer) Close()                // 滑出
func (d *Drawer) Edge() Edge
func (d *Drawer) Title() *Label         // 标题，没有标题时为 nil
func (d *Drawer) TitleHeight() float32  // 标题占的高度（像素），没有标题时为 0

func NewDialog(title, message string, buttons ...*Button) *Dialog

func (d *Dialog) Open()                 // 显示并聚焦第一个按钮，在 UI goroutine 上调用
func (d *Dialog) Close()
func (d *Dialog) Title() *Label
func (d *Dialog) Message() *Label

func Toast(text string)                 // 显示 2 秒，在 UI goroutine 上调用
```

| 选项 | 作用 |
| --- | --- |
| `ui.Gap(g)` | 抽屉里各项之间的间距 |
| `ui.Padding`、`ui.Width` 等 | 任何 `ui.Option` 都可以用在 `NewDrawer` 上 |

`NewDialog` 不接受选项，需要时对返回值调用 `ui.Apply`。

默认样式取自主题：弹窗底色 `Background`，边框 `Foreground`，圆角 `Radius`，边框内留 `PanelPadding`；抽屉标题和 Toast 用 `SmallFont` 和 `RowPadding`，抽屉标题的颜色是 `Muted`；对话框的说明用 `SmallFont`，各行之间和按钮之间隔 `Gap`（见[主题与单位](../concepts/theme.md#theme)）。
