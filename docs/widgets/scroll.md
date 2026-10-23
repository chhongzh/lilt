# Scroll

`widget.Scroll` 把一个比自己大的节点裁在框里，焦点移到哪里就滚到哪里。它不需要单独的“滚动键”：方向键照常移动焦点，`Scroll` 在旁边负责让焦点保持可见，没有焦点可去时按步滚动，到了尽头用橡皮筋回弹表示到底了。[`List`](list.md#list) 就是一个竖向的 `Scroll`，[`Choice` 和 `Toggles`](chips.md) 是横向的。

```go
import "github.com/chhongzh/lilt/widget"
```

## 竖向滚动与边缘 { #edges }

![滚到列表底部再往下推](../assets/demos/scroll-edges.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_scroll.go:scroll-edges"
```

`NewScroll(content, axes)` 的第二个参数选滚动方向：`widget.ScrollVertical`、`widget.ScrollHorizontal` 或两个方向都滚的 `widget.ScrollBoth`。沿滚动方向，内容按不受限的大小测量；`Scroll` 自己的大小不超过父节点给的空间，所以要给它尺寸，比如上面的 `ui.FillWidth()`，高度则由外面的列给。`ui.Add(s, child)` 会加到内容里（`Scroll` 实现了 [`ui.ContentHost`](../concepts/tree.md#interfaces)），`Content()` 返回内容节点本身。

滚动条是沿右边缘或下边缘的一道细线，有边框时画在边框里面。它和光标都在右边，所以例子给右边留了 3 个单位，免得细线压在光标上；`List` 默认就这样留。

按方向键时 `Scroll` 按这个顺序处理：

1. 下一个焦点在内容里、离当前视野不远（不超过半屏）：焦点照常移过去，`Scroll` 把它滚进视野，四周留 2 个单位。
2. 下一个焦点离得很远，或者后面已经没有可聚焦的节点：每按一下（或每次连发）滚 16 像素，直到滚到头。
3. 已经滚到头：按下的瞬间内容被往外拉 16 像素，按住时每次连发再多拉 16 像素，松手弹回。

拉出的距离经过 `anim.RubberBand` 换算：越往外推越费力，显示的位移趋近但永远不会达到视野的长度，和 iOS 的滚动视图一样（系数 `anim.RubberBandCoefficient`，默认 0.55）。动图后半段就是在最后一行按住 ++arrow-down++：按下时整列被拉上去一点，连发时继续拉，松手后弹回原处，滚动条在拉出时会缩短。

中间那条规则针对的是一大段不能聚焦的内容，比如一页说明文字：如果直接把焦点跳到文字后面的按钮，读者会跳过整段文字。下面的例子里焦点一直停在 `Later` 上，每按一下 ++arrow-down++ 文字滚一步，`Install` 进入半屏范围后焦点才移过去。

![按步滚过更新说明](../assets/demos/scroll-text.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_scroll.go:scroll-text"
```

触摸屏上可以直接拖动内容。手指移动超过 4 像素且方向是 `Scroll` 支持的，它就接管这次拖动，内容跟手走，拖过边缘同样有橡皮筋；松手时按最后一次移动的距离乘 8 甩出去，再停在边界之内。

## 横向滚动

![一排地图卡片横向滚动](../assets/demos/scroll-shelf.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_scroll.go:scroll-shelf"
```

横向时滚动条在下边缘，所以给下面留 3 个单位。这样的一排放进竖向的列表里，上下键在排与排之间移动，左右键在排里移动，和电视菜单一样。在一排的两端继续按左右键时，即使上下别的行还有更靠边的节点，焦点也不会斜着跳过去，而是这一排回弹。

## 波浪 { #stagger }

滚动时内容的每个子节点各用一个弹簧追上来：离焦点所在的子节点越远，弹簧越慢，整列像一道波浪那样跟上，沿滚动方向还会冲过一点。

![把层次加大到 0.6](../assets/demos/scroll-stagger.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_scroll.go:scroll-stagger"
```

离焦点 `k` 个子节点的弹簧，`Response` 是 `anim.Default` 的 `1 + stagger×k` 倍，`k` 最多算到 7，阻尼和 `anim.Default` 相同。默认 `stagger` 是 0.2，`Choice` 和 `Toggles` 用 0.15，设成 0 时整个内容作为一块移动。这些弹簧每次滚动时都从 `anim.Default` 重新算，所以调整 `anim.Default`（见[弹簧预设](../concepts/animation.md#presets)）也会改变滚动的手感。内容没有子节点时，内容本身用 `anim.Default` 滑过去。

波浪只作用在内容的直接子节点上。想让一组节点一起动，把它们放进同一个 `Row` 或 `Column` 当作一个子节点。

## 嵌套 { #nested }

`Scroll` 可以放在另一个 `Scroll` 里，比如可滚动页面里一个固定高度、自己也能滚的收件箱。

![收件箱滚到底后交给外层页面](../assets/demos/scroll-nested.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_scroll.go:scroll-nested"
```

方向键先交给离焦点最近的 `Scroll`。内层滚到头以后：

- 外面还有可聚焦的节点，焦点照常移出去，外层把它滚进视野，动图里就是从 `Message 5` 移到 `Mark all read`。
- 没有可聚焦的节点，但外层同一方向还能滚，内层不回弹，把按键让给外层，由外层按步滚动。
- 外层也到头了，才由内层回弹。

拖动时，接管的是手指下面最内层、支持这个方向的 `Scroll`。横向的一排放在竖向列表里互不干扰：横排只认左右，竖向列表只认上下。

## API

```go
type Axes uint8

const (
	ScrollVertical Axes = 1 << iota
	ScrollHorizontal
	ScrollBoth = ScrollVertical | ScrollHorizontal
)

func NewScroll(content ui.Node, axes Axes, opts ...ui.Option) *Scroll

func (s *Scroll) Content() ui.Node           // 被滚动的内容
func (s *Scroll) ScrollPosition() ui.Vec2    // 两个方向上已经滚过的距离
func (s *Scroll) SetStagger(f float32)       // 波浪的层次，0 为整块移动
func (s *Scroll) SetScrollBar(show bool)     // 显示或隐藏滚动条
func (s *Scroll) Reveal(d ui.Node)           // 把内容里的 d 滚进视野
func (s *Scroll) ArrowScroll(e *ui.KeyEvent, next ui.Node) bool // 处理方向键
func (s *Scroll) HandlePointer(e *ui.PointerEvent) bool         // 处理拖动
```

`Reveal` 和 `ArrowScroll` 分别实现了 `ui.Revealer` 和 `ui.ArrowScroller`，由焦点系统在焦点变化和按方向键时调用（见[按键分发](../concepts/input.md#dispatch)），一般不需要手动调用。`MeasureChildren`、`ArrangeChildren`、`DrawOver` 是布局和绘制的钩子。

| 选项 | 作用 |
| --- | --- |
| `ui.FillWidth()`、`ui.Height(...)`、`ui.Grow(...)` | 给 `Scroll` 尺寸，它不会自己撑满 |
| `ui.Padding(...)` | 给滚动条留位置 |
| `ui.Border(...)`、`ui.Radius(...)` | 带边框的滚动区域，滚动条画在边框里面 |

滚动条用主题的 `Foreground`，弹簧是 `anim.Snappy`（见[主题与单位](../concepts/theme.md#theme)）。内容始终至少和视野一样大，所以短内容的对齐、`Grow` 照常生效。
