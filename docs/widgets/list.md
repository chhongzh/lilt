# List 与 ListItem

`widget.ListItem` 是可聚焦的一行：左边标题，右边可选的尾部组件，松开 ++ok++ 执行 `OnClick`。`widget.List` 是放这些行的竖向滚动列表。设置页、菜单、文件列表都是这两个组件拼出来的；[`CheckBox`、`Radio`](toggle.md)、[`Select`](select.md) 也是在 `ListItem` 上加了自己的行为。

```go
import "github.com/chhongzh/lilt/widget"
```

## ListItem { #listitem }

![慢回调显示加载点，失败的回调让行抖动](../assets/demos/list-items.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_list.go:list-items"
```

`widget.Trailing(n)` 把任意节点放在行尾，常见的是一个显示当前值的小号 `Label`。标题占满剩下的宽度，过长时在行被聚焦时[滚动](label.md#marquee)。`version` 先声明出来，是为了在回调里改它；回调跑在自己的 goroutine 上，所以改文字要包在 `ui.Do` 里（[线程模型](../concepts/threading.md#do)）。

### 加载与失败 { #loading }

和 [`Button`](button.md#加载成功与失败) 一样，回调返回之前行不再响应激活。不同的是列表行没有勾和叉，等待和失败都在行尾表现：

| 回调用时 | 返回 | 表现 |
| --- | --- | --- |
| 少于 150ms | `nil` | 没有变化 |
| 150ms 以上 | `nil` | 尾部组件向上滑出，加载点出现在行尾；返回后尾部组件从下方滑回 |
| 150ms 以上 | `error` | 同上，滑回时整行左右抖一下 |
| 少于 150ms | `error` | 只抖一下 |

没有尾部组件的行，加载点直接画在行尾，盖住标题末尾的几个字。例子里 `Update` 的回调在返回前已经把 `version` 改成了 `2.4`，所以滑回来的是新值。

!!! warning "`Trailing` 会替换原来的尾部"
    一行只有一个尾部组件，再设一次会把旧的移出节点树，`SetTrailing(nil)` 则只移除。`CheckBox`、`Radio` 的勾选框和 `Select` 的当前值都是它们的尾部组件，对这几种行用 `widget.Trailing` 会把它们换掉。

`ListItem` 的内边距取自主题的 `RowPadding`，标题和尾部之间隔 4 个单位。行会裁掉超出自己的内容，滑出、滑入和抖动都不会画到别的行上。

## List { #list }

`widget.NewList()` 默认占满父节点的宽和高，`ui.Add(list, ...)` 把行追加到它内部的一列里。行比屏幕多时，列表跟着焦点滚动：

![连续下移时各行像波浪一样跟上](../assets/demos/list-wave.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_list.go:list-wave"
```

滚动时每一行都在自己的弹簧上移动：焦点所在的行最快，每远一行，弹簧的响应时间多出基准的 20%（算到 7 行为止，更远的行和第 7 行一样），于是整列像波浪一样跟上焦点，停下时沿滚动方向略微过冲。这些弹簧每次滚动都从 [`anim.Default`](../concepts/animation.md#presets) 派生，调整 `anim.Default` 会一起改变列表的手感。`list.SetStagger(0)` 让整列作为一块移动，详见 [Scroll](scroll.md#stagger)。

到了第一行或最后一行还继续按方向键，内容会被拉出一段再弹回；右侧的细条是滚动条，被拉出时会缩短。这些都是 `List` 嵌入的 [`Scroll`](scroll.md) 提供的。

`list.Column()` 返回放行的那个 `ui.Flex`，要调行距、内边距或对齐时改它：

```go
list := widget.NewList()
list.Column().SetGap(ui.CurrentTheme.Gap)
list.Column().SetPadding(ui.XY(2, 1))
```

放按钮、输入框这类带边框的组件时就这样做，让边框和屏幕边缘、彼此之间留出空隙。`NewList` 在右侧留了 3 个单位的内边距，给滚动条让位，免得它落在光标下面。

## API

```go
func NewListItem(title string, opts ...ui.Option) *ListItem

func (li *ListItem) Title() *Label               // 标题
func (li *ListItem) SetFont(f *gfx.Font)         // 标题字体
func (li *ListItem) SetTrailing(d ui.Node)       // 行尾组件，nil 移除
func (li *ListItem) OnClick[F ClickFunc](fn F)   // 来自嵌入的 Control
func (li *ListItem) Busy() bool
func (li *ListItem) Loading() bool

func NewList(opts ...ui.Option) *List

func (l *List) Column() *ui.Flex                 // 放行的那一列
func (l *List) SetStagger(f float32)             // 来自嵌入的 Scroll：波浪的强度，0 关闭
func (l *List) SetScrollBar(show bool)           // 显示或隐藏滚动条
func (l *List) ScrollPosition() ui.Vec2          // 已滚动的距离
```

| 选项 | 作用 |
| --- | --- |
| `widget.OnClick(fn)` | 行被激活时执行，`fn` 是 `func()` 或 `func() error` |
| `widget.Trailing(n)` | 行尾组件 |
| `widget.Font(f)` | 标题字体 |
| `ui.Width`、`ui.Grow` 等 | 任何 `ui.Option` 都可以用，`List` 和 `ListItem` 都接受 |

默认样式取自主题：行的内边距 `RowPadding`，标题字体 `Font`（见[主题与单位](../concepts/theme.md#theme)）。光标就是选中行的高亮，`ListItem` 自己不画背景。
