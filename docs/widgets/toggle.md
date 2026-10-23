# CheckBox 与 Radio

`widget.CheckBox` 是带勾选框的列表行，松开 ++ok++ 切换开关。`widget.Radio` 是带圆圈的列表行，同一个 `widget.RadioGroup` 里的行只能选中一个。两者都建在 [`ListItem`](list.md#listitem) 上，通常放在 [`List`](list.md#list) 里。

```go
import "github.com/chhongzh/lilt/widget"
```

## CheckBox { #checkbox }

![打开 Wi-Fi；云同步的回调失败，勾选框弹回](../assets/demos/toggle-checkbox.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_toggle.go:toggle-checkbox"
```

`NewCheckBox` 的第二个参数是初始状态。激活时勾选框**立即**翻转，然后在自己的 goroutine 上执行 `OnChange(from, to)`，接着执行 `OnClick`。回调返回 `error` 时勾选框翻回去，整行抖一下。

也就是说，回调运行期间界面已经显示新状态，失败才回滚。回调超过 150ms 时，勾选框像 `ListItem` 的其他尾部组件一样滑出，换成加载点（[加载与失败](list.md#loading)）。

| 回调返回 | 勾选框 |
| --- | --- |
| `nil` | 保持新状态 |
| `error` | 翻回旧状态，整行抖一下 |

`CheckBox` 的 `OnChange` 报告 `bool`，回调必须是 `func(from, to bool)` 或 `func(from, to bool) error`，类型不对时设置就会 panic。

用 `SetChecked` 从代码里改状态，不会触发 `OnChange`，适合在启动时或收到外部消息后同步显示：

```go
ui.Do(func() { wifi.SetChecked(connected) })
```

## Radio { #radio }

![选中 Light 再选中 Auto，上方的说明跟着更新](../assets/demos/toggle-radio.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_toggle.go:toggle-radio"
```

先建 `RadioGroup`，再把它传给每个 `NewRadio`。`NewRadioGroup` 的参数是初始选中的下标，`-1` 表示一个都不选。下标按 `NewRadio` 的调用顺序编号，和行在节点树里的位置无关，所以同一组的行可以分散在页面的不同地方。

`RadioGroup` 不是节点，不能接收选项，`OnChange` 要直接调用它的方法：`theme.OnChange(fn)`。回调拿到旧下标和新下标，类型是 `int`。

选中一行时：

1. 组内的圆圈立即切到这一行；
2. 下标变了的话，执行组的 `OnChange(from, to)`；
3. 执行这一行自己的 `OnClick`（用 `widget.OnClick` 传给 `NewRadio`）。

回调运行期间，加载点显示在被选中的那一行，整组都不响应激活。回调返回 `error` 时，选中状态回到旧的那一行，被选中的行抖一下。再选一次已经选中的行，`OnChange` 不会执行，`OnClick` 照常执行。

`Select(i)` 从代码里选中第 `i` 行，不触发 `OnChange`；`Selected()` 读当前下标；`Radios()` 按下标返回组里的所有行。

## API

```go
func NewCheckBox(title string, checked bool, opts ...ui.Option) *CheckBox

func (cb *CheckBox) Checked() bool
func (cb *CheckBox) SetChecked(v bool)            // 不触发 OnChange

func NewRadioGroup(selected int) *RadioGroup      // -1 表示都不选

func (g *RadioGroup) Selected() int               // 当前下标，-1 表示都没选
func (g *RadioGroup) Select(i int)                // 不触发 OnChange
func (g *RadioGroup) Radios() []*Radio            // 按下标排列
func (g *RadioGroup) OnChange[F ChangeFunc](fn F) // 来自嵌入的 Control，值类型 int

func NewRadio(g *RadioGroup, title string, opts ...ui.Option) *Radio
```

`CheckBox` 和 `Radio` 嵌入了 `ListItem`，`Title`、`SetFont`、`OnClick`、`Busy`、`Loading` 等方法都可以用。

| 选项 | 作用 |
| --- | --- |
| `widget.OnChange(fn)` | 只用于 `CheckBox`：开关变化后执行，`fn` 是 `func(from, to bool)` 或带 `error` 返回值的版本 |
| `widget.OnClick(fn)` | 激活后执行，在 `OnChange` 之后 |
| `widget.Font(f)` | 标题字体 |
| `ui.Width`、`ui.Grow` 等 | 任何 `ui.Option` 都可以用 |

勾选框和圆圈边长 9 个单位，用主题的 `Foreground` 画，开关时内部的方块或圆点由 `anim.Bouncy` 弹簧放大、缩小（[弹簧预设](../concepts/animation.md#presets)）。行的内边距同 `ListItem`，取自主题的 `RowPadding`。
