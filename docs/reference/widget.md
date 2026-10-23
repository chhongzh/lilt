# widget

[pkg.go.dev](https://pkg.go.dev/github.com/chhongzh/lilt/widget) · `import "github.com/chhongzh/lilt/widget"`

建立在 `ui` 之上的标准组件。所有构造函数最后一个参数都是 `...ui.Option`，下面的选项和任何 `ui.Option` 都能传进去。

## 事件与 Control

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Control struct` | 可交互组件嵌入它，得到回调、Busy 和 Loading 状态 | [线程模型](../concepts/threading.md#control) |
| `func (c *Control) OnClick[F ClickFunc](fn F)` | 设置激活时运行的回调 | [Button](../widgets/button.md) |
| `func (c *Control) OnChange[F ChangeFunc](fn F)` | 设置用户改值后运行的回调，参数是旧值和新值；类型和组件的值类型不符时 panic | [线程模型](../concepts/threading.md#control) |
| `func (c *Control) Busy() bool` | 回调是否还在运行 | [线程模型](../concepts/threading.md#control) |
| `func (c *Control) Loading() bool` | 回调是否已运行超过 150ms，组件正在显示加载点 | [Button](../widgets/button.md) |
| `type ClickFunc interface` | `OnClick` 回调的类型：`func()` 或 `func() error`，返回 error 表示失败 | [Button](../widgets/button.md) |
| `type ChangeFunc interface` | `OnChange` 回调的类型，见下 | [线程模型](../concepts/threading.md#control) |
| `func OnClick[F ClickFunc](fn F) ui.Option` | 选项版的 `OnClick`，用于嵌入了 `Control` 的组件 | [激活](../concepts/input.md#activation) |
| `func OnChange[F ChangeFunc](fn F) ui.Option` | 选项版的 `OnChange` | [线程模型](../concepts/threading.md#control) |

`ChangeFunc` 允许的回调类型，以及哪些组件报告哪种值：

```go
type ChangeFunc interface {
	func(from, to bool) | func(from, to bool) error | // CheckBox
		func(from, to int) | func(from, to int) error | // Choice、Select、RadioGroup
		func(from, to string) | func(from, to string) error | // Input
		func(from, to uint32) | func(from, to uint32) error // Toggles，第 i 位是第 i 个 Chip
}
```

回调返回 error 时，值类组件恢复旧值。`widget` 包以外的组件嵌入 `Control` 拿不到这套机制，见[自定义组件](../guides/custom-widget.md#control)。

## 通用选项与能力接口

选项按接口生效：节点有对应的方法就能用，没有就 panic。

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Font(f *gfx.Font) ui.Option` | 设置 Label 或组件内部文字的字体，要求 `FontSetter` | [Label](../widgets/label.md) |
| `func TextColor(c gfx.Color) ui.Option` | 文字颜色，要求 `Colorer` | [Label](../widgets/label.md) |
| `func TextAlign(x float32) ui.Option` | 水平对齐，0 靠左到 1 靠右，要求 `Aligner` | [Label](../widgets/label.md) |
| `func MarqueeMode(m Marquee) ui.Option` | 跑马灯模式，只用于 `*Label` | [跑马灯](../widgets/label.md#marquee) |
| `func Trailing(t ui.Node) ui.Option` | 在行尾放一个节点，要求 `Trailer` | [ListItem](../widgets/list.md#listitem) |
| `func HoldToConfirm(dur time.Duration) ui.Option` | 按住 `dur` 才触发，只用于 `*Button` | [Button](../widgets/button.md) |
| `type FontSetter interface{ SetFont(f *gfx.Font) }` | 能换字体 | [Label](../widgets/label.md) |
| `type Colorer interface{ SetColor(c gfx.Color) }` | 能换文字颜色 | [Label](../widgets/label.md) |
| `type Aligner interface{ SetAlign(x float32) }` | 能水平对齐文字 | [Label](../widgets/label.md) |
| `type Texter interface{ Text() string; SetText(s string) }` | 有一行文字，例如 `Label` | [Label](../widgets/label.md) |
| `type Trailer interface{ SetTrailing(d ui.Node) }` | 能在行尾放节点 | [ListItem](../widgets/list.md#listitem) |

## Label

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewLabel(text string, opts ...ui.Option) *Label` | 一行文字，在盒子里垂直居中 | [Label](../widgets/label.md) |
| `func (l *Label) Text() string`、`SetText(s string)` | 读写文字 | [Label](../widgets/label.md) |
| `func (l *Label) Font() *gfx.Font`、`SetFont(f *gfx.Font)` | 字体，未设置时用主题的 `Font` | [Label](../widgets/label.md) |
| `func (l *Label) SetColor(c gfx.Color)` | 覆盖主题的前景色 | [Label](../widgets/label.md) |
| `func (l *Label) SetAlign(x float32)` | 水平对齐 | [Label](../widgets/label.md) |
| `func (l *Label) SetMarquee(m Marquee)` | 跑马灯模式 | [跑马灯](../widgets/label.md#marquee) |
| `type Marquee uint8` | 放不下的 Label 何时滚动 | [跑马灯](../widgets/label.md#marquee) |
| `MarqueeFocused` | 默认：自己或祖先有焦点时滚动 | [跑马灯](../widgets/label.md#marquee) |
| `MarqueeAlways` | 一直滚动 | [跑马灯](../widgets/label.md#marquee) |
| `Label` 的 `Draw`、`Measure`、`Update`、`Leaf` | `ui` 能力接口的实现，由框架调用 | [能力接口](../guides/custom-widget.md#interfaces) |

## List、ListItem 与 Scroll

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewList(opts ...ui.Option) *List` | 纵向滚动的一列，`ui.Add` 把行加进去；默认铺满父节点 | [List](../widgets/list.md#list) |
| `func (l *List) Column() *ui.Flex` | 放行的那一列，用来调间距和对齐 | [List](../widgets/list.md#list) |
| `func NewListItem(title string, opts ...ui.Option) *ListItem` | 可聚焦的一行：标题加可选的尾部节点 | [ListItem](../widgets/list.md#listitem) |
| `func (li *ListItem) Title() *Label` | 标题 | [ListItem](../widgets/list.md#listitem) |
| `func (li *ListItem) SetFont(f *gfx.Font)` | 标题字体 | [ListItem](../widgets/list.md#listitem) |
| `func (li *ListItem) SetTrailing(d ui.Node)` | 换掉尾部节点 | [ListItem](../widgets/list.md#listitem) |
| `ListItem` 的 `HandleKey`、`Update`、`DrawOver` | 能力接口的实现：松开 OK 时激活，加载时尾部滑走、换成加载点 | [能力接口](../guides/custom-widget.md#interfaces) |
| `func NewScroll(content ui.Node, axes Axes, opts ...ui.Option) *Scroll` | 裁剪 `content` 并按方向滚动 | [Scroll](../widgets/scroll.md) |
| `type Axes uint8` | 滚动方向 | [Scroll](../widgets/scroll.md) |
| `ScrollVertical`、`ScrollHorizontal`、`ScrollBoth` | 纵向、横向、两个方向 | [Scroll](../widgets/scroll.md) |
| `func (s *Scroll) Content() ui.Node` | `ui.Add` 实际添加到的节点 | [Scroll](../widgets/scroll.md) |
| `func (s *Scroll) ScrollPosition() ui.Vec2` | 沿 x、y 已滚动的距离 | [Scroll](../widgets/scroll.md) |
| `func (s *Scroll) SetScrollBar(show bool)` | 显示或隐藏边缘的滚动条 | [Scroll](../widgets/scroll.md) |
| `func (s *Scroll) SetStagger(f float32)` | 离焦点每远一步跟得慢多少，0 整体移动 | [波浪](../widgets/scroll.md#stagger) |
| `Scroll` 的 `ArrowScroll`、`Reveal`、`HandlePointer`、`MeasureChildren`、`ArrangeChildren`、`DrawOver` | 能力接口的实现：方向键滚动、焦点跟随、拖动、布局、滚动条 | [嵌套滚动](../widgets/scroll.md#nested)、[边缘回弹](../widgets/scroll.md#edges) |

## Button

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewButton(text string, opts ...ui.Option) *Button` | 带边框的按钮，松开 OK 时运行 `OnClick` | [Button](../widgets/button.md) |
| `func (b *Button) Label() *Label` | 按钮里的文字 | [Button](../widgets/button.md) |
| `func (b *Button) SetFont(f *gfx.Font)` | 文字字体 | [Button](../widgets/button.md) |
| `func (b *Button) SetHold(d time.Duration)` | 同 `HoldToConfirm`，0 关闭 | [Button](../widgets/button.md) |
| `func (b *Button) CursorShape(r ui.Rect) ui.Rect` | 长按时从左边切掉光标 | [改变光标](../guides/custom-widget.md#cursorshaper) |
| `Button` 的 `Draw`、`HandleKey`、`Update` | 能力接口的实现 | [能力接口](../guides/custom-widget.md#interfaces) |

## CheckBox 与 Radio

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewCheckBox(title string, checked bool, opts ...ui.Option) *CheckBox` | 松开 OK 切换勾选的一行，`OnChange` 收 `bool` | [CheckBox](../widgets/toggle.md#checkbox) |
| `func (cb *CheckBox) Checked() bool` | 是否勾选 | [CheckBox](../widgets/toggle.md#checkbox) |
| `func (cb *CheckBox) SetChecked(v bool)` | 改勾选状态，不触发 `OnChange` | [CheckBox](../widgets/toggle.md#checkbox) |
| `func NewRadioGroup(selected int) *RadioGroup` | 一组单选，`-1` 表示都不选；`OnChange` 收 `int` | [Radio](../widgets/toggle.md#radio) |
| `func (g *RadioGroup) Selected() int` | 选中的序号 | [Radio](../widgets/toggle.md#radio) |
| `func (g *RadioGroup) Select(i int)` | 选中第 `i` 个，不触发 `OnChange` | [Radio](../widgets/toggle.md#radio) |
| `func (g *RadioGroup) Radios() []*Radio` | 组里的所有行 | [Radio](../widgets/toggle.md#radio) |
| `func NewRadio(g *RadioGroup, title string, opts ...ui.Option) *Radio` | 属于 `g` 的一行；它的 `OnClick` 在组的 `OnChange` 之后运行 | [Radio](../widgets/toggle.md#radio) |

## Select

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewSelect(title string, options []string, selected int, opts ...ui.Option) *Select` | 显示当前选项的一行，激活时打开一个单选的底部抽屉；`OnChange` 收 `int` | [Select](../widgets/select.md) |
| `func (s *Select) Selected() int` | 选中的序号 | [Select](../widgets/select.md) |
| `func (s *Select) SetSelected(i int)` | 显示第 `i` 项，不触发 `OnChange` | [Select](../widgets/select.md) |

## Chip、Choice 与 Toggles

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewChip(text string, opts ...ui.Option) *Chip` | 一个可开关的词，开着时下面有横条 | [Chip](../widgets/chips.md#chip) |
| `func (c *Chip) On() bool`、`SetOn(on bool)` | 读、改开关，横条动画过渡 | [Chip](../widgets/chips.md#chip) |
| `func (c *Chip) Label() *Label`、`SetFont(f *gfx.Font)` | 文字和字体 | [Chip](../widgets/chips.md#chip) |
| `Chip` 的 `Draw`、`HandleKey`、`Update` | 能力接口的实现 | [能力接口](../guides/custom-widget.md#interfaces) |
| `var ChipSpring = &anim.Spring{Response: 0.262, Damping: 0.5}` | 横条从中间长出、略超过文字宽度再落定的弹簧 | [Chip](../widgets/chips.md#chip) |
| `func NewChoice(labels []string, chosen int, opts ...ui.Option) *Choice` | 一排 Chip 里恰好一个开着；`OnChange` 收旧、新序号 | [Choice](../widgets/chips.md#choice) |
| `func (c *Choice) Chosen() int` | 选中的序号 | [Choice](../widgets/chips.md#choice) |
| `func (c *Choice) Choose(i int)` | 切到第 `i` 个，不触发 `OnChange`；`-1` 全关 | [Choice](../widgets/chips.md#choice) |
| `func (c *Choice) Chips() []*Chip`、`SetFont(f *gfx.Font)` | 所有 Chip；统一换字体 | [Choice](../widgets/chips.md#choice) |
| `func NewToggles(labels []string, opts ...ui.Option) *Toggles` | 一排 Chip 任意开关；`OnChange` 收位掩码，最多 32 个 | [Toggles](../widgets/chips.md#toggles) |
| `func (t *Toggles) Mask() uint32`、`SetMask(m uint32)` | 读、改整组开关，第 `i` 位是第 `i` 个；`SetMask` 不触发 `OnChange` | [Toggles](../widgets/chips.md#toggles) |
| `func (t *Toggles) On(i int) bool`、`SetOn(i int, on bool)` | 读、改单个，`SetOn` 不触发 `OnChange` | [Toggles](../widgets/chips.md#toggles) |
| `func (t *Toggles) Chips() []*Chip`、`SetFont(f *gfx.Font)` | 所有 Chip；统一换字体 | [Toggles](../widgets/chips.md#toggles) |
| `Choice`、`Toggles` 的 `Update` | 能力接口的实现 | [能力接口](../guides/custom-widget.md#interfaces) |

## Input 与软键盘

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewInput(placeholder string, opts ...ui.Option) *Input` | 单行输入框，激活时打开 `ui.TextInput` 或软键盘；`OnChange` 收旧、新文字 | [Input](../widgets/input.md#input) |
| `func (in *Input) Text() string` | 当前文字 | [Input](../widgets/input.md#input) |
| `func (in *Input) SetText(s string)` | 替换文字，不触发 `OnChange` | [Input](../widgets/input.md#input) |
| `func (in *Input) SetFont(f *gfx.Font)` | 字体 | [Input](../widgets/input.md#input) |
| `Input` 的 `HandleKey`、`Update`、`DrawOver` | 能力接口的实现；直接接收物理键盘的字符和退格 | [物理键盘](../guides/hardware.md#keys) |
| `type SoftKeyboard struct` | 从底部抽屉弹出的屏幕键盘，方向键用焦点查找在按键间移动 | [软键盘](../widgets/input.md#keyboard) |
| `func (k *SoftKeyboard) Open(t ui.TextTarget)` | 打开键盘编辑 `t` | [软键盘](../widgets/input.md#keyboard) |
| `var Keyboard = &SoftKeyboard{}` | `ui.TextInput` 为 nil 时 `Input` 使用的键盘 | [文字输入](../widgets/input.md#text-input) |

## 弹窗

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func NewDrawer(edge Edge, title string, opts ...ui.Option) *Drawer` | 从屏幕一边滑入的一列，`title` 为空时没有标题 | [Drawer](../widgets/popups.md#drawer) |
| `func (d *Drawer) Open()`、`Close()` | 滑入、滑出；`Open` 在 UI goroutine 上调用 | [Drawer](../widgets/popups.md#drawer) |
| `func (d *Drawer) Edge() Edge` | 从哪一边滑入 | [Drawer](../widgets/popups.md#drawer) |
| `func (d *Drawer) Title() *Label` | 标题，没有时为 nil | [Drawer](../widgets/popups.md#drawer) |
| `func (d *Drawer) TitleHeight() float32` | 标题占的高度，没有标题为 0 | [Drawer](../widgets/popups.md#drawer) |
| `type Edge uint8` | 抽屉的边 | [Drawer](../widgets/popups.md#drawer) |
| `Bottom`、`Top`、`Left`、`Right` | 下、上、左、右 | [Drawer](../widgets/popups.md#drawer) |
| `func NewDialog(title, message string, buttons ...*Button) *Dialog` | 屏幕中央的对话框：标题、一行消息、平分宽度的一排按钮；按钮回调成功返回后自动关闭 | [Dialog](../widgets/popups.md#dialog) |
| `func (d *Dialog) Open()`、`Close()` | 打开（第一个按钮获得焦点）、关闭；`Open` 在 UI goroutine 上调用 | [Dialog](../widgets/popups.md#dialog) |
| `func (d *Dialog) Title() *Label`、`Message() *Label` | 标题和消息 | [Dialog](../widgets/popups.md#dialog) |
| `func Toast(text string)` | 在底部显示 `text` 两秒，不抢焦点、不接按键，替换还在显示的 Toast；在 UI goroutine 上调用 | [Toast](../widgets/popups.md#toast) |
