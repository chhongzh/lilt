# Input 与软键盘

`widget.Input` 是单行文本框。板子上通常只有几个方向键，所以松开 ++ok++ 时它打开一个输入法来编辑文字，默认是屏幕上的软键盘 `widget.Keyboard`；接了实体键盘时，字符直接打进框里。

```go
import "github.com/chhongzh/lilt/widget"
```

## 输入框 { #input }

![用软键盘输入 tea](../assets/demos/field-keyboard.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_field.go:field-keyboard"
```

框里没有文字时显示灰色的占位文字（`NewInput` 的第一个参数，颜色取主题的 `Muted`）。文字改动后，`OnChange` 在自己的 goroutine 上收到改动前后的两段文字，所以改 `results` 要包在 `ui.Do` 里。用软键盘输入时，键盘里的每一下都只改键盘自己的缓冲区，按键盘上的 `OK` 才把整段文字交给 `Input`，所以上面的例子从空到 `tea` 只调用一次 `OnChange`。

`Text()` 读当前文字，`SetText(s)` 在代码里替换文字，不调用 `OnChange`。框里没有插入点：实体键盘的字符总是加在末尾，退格删掉最后一个字符。

### 回调期间的改动

`OnChange` 还没返回时又来了改动（多见于实体键盘，一个字符一次改动），`Input` 不会再起一个回调，而是先记着，等这次回调返回后，把这期间的所有改动合成一次报告：`from` 是上次报告的结果，`to` 是现在的文字。所以回调永远不会并发执行，也不会漏掉最后的状态。

`OnChange` 返回 `error` 时，框里的文字退回这次改动之前的 `from`，期间新打的字也一起丢掉，框左右抖一下。回调超过 150ms 还没返回，框的右端显示跳动的加载点。

| 情况 | 表现 |
| --- | --- |
| 回调少于 150ms，返回 `nil` | 无额外反馈 |
| 回调 150ms 以上 | 框右端显示加载点，期间的改动攒着 |
| 回调返回 `nil`，期间有改动 | 立刻把攒下的改动作为一次 `OnChange` 报告 |
| 回调返回 `error` | 文字退回 `from`，框抖动 |

`Input` 也可以设 `OnClick`：每次松开 ++ok++ 都会执行，同时照常打开输入法。

### 实体键盘

字符和退格都通过 `ui` 包的队列进来，可以在任何 goroutine 上调用：

```go
// pump 把键盘芯片读到的字符送进界面，跑在自己的 goroutine 上。
func pump(chars <-chan rune) {
	for r := range chars {
		switch r {
		case '\b':
			ui.Emit(ui.KeyBackspace, ui.Press)
			ui.Emit(ui.KeyBackspace, ui.Release)
		case '\r':
			ui.Emit(ui.KeyOK, ui.Press)
			ui.Emit(ui.KeyOK, ui.Release)
		default:
			ui.EmitRune(r)
		}
	}
}
```

`ui.EmitRune(r)` 发出一个 `ui.KeyChar` 事件，字符在 `KeyEvent.Rune` 里；`ui.KeyBackspace` 在按下和连发时各删一个字符，松开不删。字符送给获得焦点的节点：焦点在 `Input` 上时打进框里，软键盘开着时打进键盘的缓冲区，在别处则被忽略。事件队列有 32 格，满了会丢弃新事件，不会阻塞发送方（见[按键分发](../concepts/input.md#dispatch)）。

## 软键盘 { #keyboard }

`widget.Keyboard` 是 `*widget.SoftKeyboard` 类型的全局变量，`Input` 在没有设置 `ui.TextInput` 时打开它。它从底部升起，作为 `ui.DrawerBottom` 弹窗打开，在 128x64 的屏幕上正好占满。

| 行 | 字母层 | 符号层 |
| --- | --- | --- |
| 预览 | 已输入的文字和 `_` 插入点 | 同左 |
| 1 | `q w e r t y u i o p` | `1 2 3 4 5 6 7 8 9 0` |
| 2 | `a s d f g h j k l` | `- / : ; ( ) $ & @` |
| 3 | `^ z x c v b n m <` | `= . , ? ! ' " # <` |
| 4 | `123  space  OK` | `abc  space  OK` |

- 方向键在键之间移动，用的是全局的[焦点查找](../concepts/input.md#focus)，所以长短不一的行也能上下对齐到最近的键。松开 ++ok++ 按下光标下的键。
- `^` 是一次性的 Shift：开着时字母显示成大写、`^` 下面有一道线，打完一个字母自动关掉。
- `123` 切到数字和符号，`abc` 切回字母；`<` 删除最后一个字符，`space` 输入空格。
- 文字超出预览框时改成右对齐，插入点始终可见。
- `OK` 把文字交给目标并关掉键盘。按 ++back++ 也会关掉键盘，但文字不交回去，相当于取消。

键盘不只给 `Input` 用。任何实现了 `ui.TextTarget`（`Text()` 和 `SetText(s)`）的东西都可以拿去编辑，例如直接改一个 `Label`：

```go
// title lets the keyboard edit a label in place.
type title struct{ l *widget.Label }

func (t title) Text() string     { return t.l.Text() }
func (t title) SetText(s string) { t.l.SetText(s) }
```

在 UI goroutine 上调用 `widget.Keyboard.Open(title{l})` 即可，回调里要包在 `ui.Do` 里。gallery 的待办示例（`examples/gallery/todo.go`）就是这样给任务改名的。

## 自定义输入法 { #text-input }

`Input` 打开的是全局变量 `ui.TextInput`，为 `nil` 时才用 `widget.Keyboard`。实现一个 `ui.TextInputMethod` 并在 `lilt.Start` 之前赋给它，就能换掉所有输入框的编辑方式：例如按键太少的板子只给几个预设选项，或者接了实体键盘的板子打开一个只有预览和提示的界面。

![用预设名字代替键盘](../assets/demos/field-method.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_field.go:field-method-type"
```

```go
--8<-- "cmd/docgif/demo_field.go:field-method"
```

`ui.TextInputMethod` 只有一个方法 `Open(t ui.TextTarget)`，`Input` 在 UI goroutine 上调用它，所以 `Open` 里可以直接建树、打开弹窗。之后把结果交回去要调用 `t.SetText`，它会改动输入框并触发 `OnChange`，必须在 UI goroutine 上调用：上面的列表项回调跑在自己的 goroutine 上，所以和 `d.Close()` 一起包在 `ui.Do` 里。不调用 `SetText` 就关掉，输入框保持原样。

## API

```go
func NewInput(placeholder string, opts ...ui.Option) *Input

func (in *Input) Text() string
func (in *Input) SetText(s string)          // 不调用 OnChange
func (in *Input) SetFont(f *gfx.Font)
func (in *Input) OnChange[F ChangeFunc](fn F) // 来自嵌入的 Control
func (in *Input) Busy() bool
func (in *Input) Loading() bool

var Keyboard = &SoftKeyboard{}              // Input 默认打开的软键盘

func (k *SoftKeyboard) Open(t ui.TextTarget) // 打开键盘编辑 t
```

`ui` 包里相关的部分：

```go
type TextTarget interface {
	Text() string
	SetText(s string)
}

type TextInputMethod interface{ Open(t TextTarget) }

var TextInput TextInputMethod // Input 打开的输入法，nil 表示 widget.Keyboard

func EmitRune(r rune) // 送入一个字符，任何 goroutine 都可以调用
```

| 选项 | 作用 |
| --- | --- |
| `widget.OnChange(fn)` | 文字改动后执行，`fn` 是 `func(from, to string)` 或 `func(from, to string) error` |
| `widget.OnClick(fn)` | 每次激活时执行 |
| `widget.Font(f)` | 文字字体 |
| `ui.Width`、`ui.Grow` 等 | 任何 `ui.Option` 都可以用 |

默认样式取自主题：边框用 `Foreground`，圆角 `Radius`，内边距 `RowPadding`，占位文字用 `Muted`（见[主题与单位](../concepts/theme.md#theme)）。软键盘的键和预览用 `SmallFont`。
