# Choice 与 Toggles

一排横着的词，开着的词下面有一道短横线。`widget.Choice` 是单选：同一时间只开一个，适合难度、速率这类互斥的选项；`widget.Toggles` 是多选：每个词各开各的，适合一组独立的开关。两者都由 `widget.Chip` 组成，左右键在词之间移动，松开 ++ok++ 时切换光标下的那个。

```go
import "github.com/chhongzh/lilt/widget"
```

## Chip { #chip }

`widget.Chip` 是一个可聚焦的词，开着时下面画一道 2 像素高的横线。`SetOn` 切换它时，横线从词的正中间往两边长出来，由 `widget.ChipSpring`（`Response: 0.262, Damping: 0.5`）驱动：阻尼小，横线会先冲过词宽，最多再长出四分之一，再弹回词宽，落定时有一下“咔哒”的感觉。关掉时横线往中间缩回去。

`Chip` 自己不决定开着是什么意思。单独使用时，松开 ++ok++ 只执行它的 `OnClick`，开关要你在回调里用 `ui.Do` 调 `SetOn`。通常不单独用，而是交给 `Choice` 或 `Toggles` 管：它们会在 UI goroutine 上先切换词，再跑回调。

想改动画手感，改 `widget.ChipSpring` 的字段即可：所有词都指向这一个弹簧，改了立刻生效（见[弹簧](../concepts/animation.md#spring)）。

## 单选 { #choice }

![在难度之间切换](../assets/demos/chips-choice.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_chips.go:chips-choice"
```

`NewChoice` 的第二个参数是一开始开着的序号，传 `-1` 表示一个都不开。初始状态直接画好，不播放动画。松开 ++ok++ 时旧选项下的横线缩回、新选项下的长出来，然后 `OnChange` 在自己的 goroutine 上收到旧序号和新序号，所以改 `stars` 的文字要包在 `ui.Do` 里。在已经开着的词上按 ++ok++ 不算改动，不会调用 `OnChange`。

选项比屏幕宽时，`Choice` 会横向滚动，光标移到哪里就把哪个词滚进屏幕，动图后半段就是这样滚到 `Expert` 的。它本身就是一个横向的 [`Scroll`](scroll.md)，只是关掉了滚动条（细线会被误认成某个词下面的横线），层次感也比列表轻一些（`SetStagger(0.15)`，见[波浪](scroll.md#stagger)）。在一排的两端继续按左右键，这一排会被拉开一点再弹回来，光标不会斜着跳到上下别的行。

在代码里用 `Choose(i)` 切换选项，它不调用 `OnChange`。`Chips()` 返回所有 `Chip`，可以用来改某个词的文字或给它单独设 `OnClick`。

## 多选 { #toggles }

![打开和关闭几个模组](../assets/demos/chips-toggles.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_chips.go:chips-toggles"
```

`Toggles` 一开始全关。松开 ++ok++ 翻转光标下的词，`OnChange` 收到的是改动前后的整组状态，用 `uint32` 位掩码表示：第 `i` 位对应第 `i` 个词。所以一排最多 32 个词。

代码里读写状态用这几个方法，都不调用 `OnChange`：

| 方法 | 作用 |
| --- | --- |
| `Mask()` | 当前开着的词组成的位掩码 |
| `SetMask(m)` | 按 `m` 把每个词设成开或关 |
| `On(i)` | 第 `i` 个词是否开着 |
| `SetOn(i, on)` | 只切换第 `i` 个词 |

## 加载与回滚

`OnChange` 可以返回 `error`。回调跑得久，触发改动的那个词下面的横线换成跳动的加载点；回调失败，整排左右抖一下，选项回到改动之前。

![风扇档位：Eco 成功，Turbo 失败后退回 Eco](../assets/demos/chips-fail.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_chips.go:chips-fail"
```

| 回调用时 | 返回 | 表现 |
| --- | --- | --- |
| 少于 150ms | `nil` | 横线直接落到新选项下 |
| 150ms 以上 | `nil` | 新选项下先是加载点，回调返回后换成横线 |
| 任意 | `error` | 加载点（若超过 150ms）→ 整排抖动 → 恢复改动前的选项 |

回调没返回之前，整排不再接受新的切换，按 ++ok++ 不会改动选项。选项是先切过去、失败再退回来的，所以界面不用等回调就能立刻响应；这套 Busy、Loading 的规则来自嵌入的 [`widget.Control`](../concepts/threading.md#control)。`Toggles` 失败时用 `SetMask` 退回改动前的整组状态。

## API

```go
var ChipSpring = &anim.Spring{Response: 0.262, Damping: 0.5} // 横线的弹簧

func NewChip(text string, opts ...ui.Option) *Chip

func (c *Chip) Label() *Label        // 词本身
func (c *Chip) SetFont(f *gfx.Font)
func (c *Chip) On() bool
func (c *Chip) SetOn(on bool)        // 切换，横线动画过去

func NewChoice(labels []string, chosen int, opts ...ui.Option) *Choice

func (c *Choice) Chosen() int        // 开着的序号，-1 表示没有
func (c *Choice) Choose(i int)       // 切换到 i，不调用 OnChange；-1 全关
func (c *Choice) Chips() []*Chip
func (c *Choice) SetFont(f *gfx.Font) // 所有词的字体

func NewToggles(labels []string, opts ...ui.Option) *Toggles

func (t *Toggles) Mask() uint32       // 第 i 位是第 i 个词
func (t *Toggles) SetMask(m uint32)   // 不调用 OnChange
func (t *Toggles) On(i int) bool
func (t *Toggles) SetOn(i int, on bool) // 不调用 OnChange
func (t *Toggles) Chips() []*Chip
func (t *Toggles) SetFont(f *gfx.Font)
```

`Choice`、`Toggles` 和 `Chip` 都嵌入了 `Control`，也有 `OnClick`、`OnChange`、`Busy`、`Loading`。

| 选项 | 作用 |
| --- | --- |
| `widget.OnChange(fn)` | `Choice` 用 `func(from, to int)`，`Toggles` 用 `func(from, to uint32)`，都可以多返回一个 `error` |
| `widget.OnClick(fn)` | `Chip` 被激活时执行；在 `Choice`、`Toggles` 里，它在切换之后执行 |
| `widget.Font(f)` | 字体，用在 `Choice`、`Toggles` 上时作用于每个词 |
| `ui.Width`、`ui.Margin` 等 | 任何 `ui.Option` 都可以用 |

`OnChange` 的类型和组件不符时会 panic。默认样式取自主题：文字用 `Font`，横线和加载点用 `Foreground`（见[主题与单位](../concepts/theme.md#theme)）。每个词左右各留 3 个单位，词与词之间隔 2 个单位，`Choice` 和 `Toggles` 默认铺满父节点的宽度。
