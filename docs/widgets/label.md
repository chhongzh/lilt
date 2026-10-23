# Label

`widget.Label` 是一行文字。它不可聚焦，也不响应按键；按钮、列表行、输入框里的文字都是一个 `Label`。文字放不下时，它会以跑马灯的方式来回滚动。

```go
import "github.com/chhongzh/lilt/widget"
```

## 字体、颜色与对齐

![几种样式的 Label](../assets/demos/label-styles.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_label.go:label-styles"
```

- **字体**：默认用主题的 `Font`，`widget.Font(f)` 换成别的，常用的是主题里的 `SmallFont`。
- **颜色**：默认用主题的 `Foreground`，`widget.TextColor(c)` 覆盖它。提示、说明这类次要文字用主题的 `Muted`。
- **对齐**：`widget.TextAlign(x)` 中 `0` 靠左、`0.5` 居中、`1` 靠右，竖直方向总是居中。

`Label` 测量出的宽度就是文字的宽度（不超过可用宽度），所以只有在它的盒子比文字宽时，对齐才有意义。例子里给需要对齐的标签加了 `ui.FillWidth()`，让它占满整行（[尺寸](../concepts/layout.md#size)）。最后一行的 `Play` 是一个列表行，尾部的 `6:12` 也是 `Label`，见 [ListItem](list.md#listitem)。

`Label` 实现了 `Leaf()`，表示它不接受子节点；对它调用 `ui.Add` 会 panic。

## 跑马灯 { #marquee }

文字比盒子宽时，`Label` 会先停 1 秒，以每秒 24 像素向左滚到末尾（末尾多留 2 像素空白），再停 1 秒，然后由弹簧带回开头，如此循环。何时滚动由 `widget.MarqueeMode` 决定：

| 模式 | 什么时候滚动 |
| --- | --- |
| `widget.MarqueeFocused`（默认） | 标签自己或它的某个祖先有焦点时 |
| `widget.MarqueeAlways` | 一直滚 |

![列表行的标题在聚焦时滚动，顶部的标签一直滚动](../assets/demos/label-marquee.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_label.go:label-marquee"
```

`Label` 不可聚焦，默认模式看的是祖先：列表行的标题在行被聚焦时才滚，所以一屏长标题里只有正在读的那一行在动。焦点离开后，标题滑回开头，下次聚焦重新从停顿开始。顶部的 `news` 用了 `MarqueeAlways`，焦点在哪都不影响它。

几个细节：

- 放不下时，`TextAlign` 只决定文字停在哪一端：大于 `0.5` 时停在末尾，显示文字的结尾，滚动时把它带到开头再回来；其余（包括居中）停在开头。文件路径这类结尾更重要的文字适合靠右放，见上一节例子里的路径。
- `SetText` 换了内容、`SetMarquee` 换了模式，都会让跑马灯从头开始。

`widget.MarqueeMode` 只能用在 `Label` 上。列表行、按钮的文字要取出里面的标签再设：`row.Title().SetMarquee(widget.MarqueeAlways)`。

## API

```go
func NewLabel(text string, opts ...ui.Option) *Label

func (l *Label) Text() string
func (l *Label) SetText(s string)          // 内容不变时什么也不做；变了会重新布局
func (l *Label) Font() *gfx.Font           // 没设置时返回主题的 Font
func (l *Label) SetFont(f *gfx.Font)       // 会重新布局
func (l *Label) SetColor(c gfx.Color)      // 覆盖主题的 Foreground
func (l *Label) SetAlign(x float32)        // 0 靠左，1 靠右
func (l *Label) SetMarquee(m Marquee)      // MarqueeFocused 或 MarqueeAlways
func (l *Label) Leaf()                     // 不接受子节点
```

| 选项 | 作用 |
| --- | --- |
| `widget.Font(f)` | 字体 |
| `widget.TextColor(c)` | 文字颜色 |
| `widget.TextAlign(x)` | 水平对齐，0 到 1 |
| `widget.MarqueeMode(m)` | 跑马灯模式 |
| `ui.Padding`、`ui.FillWidth` 等 | 任何 `ui.Option` 都可以用 |

默认样式取自主题：字体 `Font`，颜色 `Foreground`（见[主题与单位](../concepts/theme.md#theme)）。没有用 `TextColor` 设置过颜色的 `Label` 在每次绘制时读主题，改了主题的 `Foreground` 会跟着变；字体则在构造、`SetText`、`SetFont` 时测量，改主题字体要放在构造组件之前。
