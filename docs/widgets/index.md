# 组件

`widget` 包是建在 `ui` 之上的标准组件。它们都是普通的节点：用 `ui.Add` 放进节点树，用 `ui.Option` 配置，和自己写的组件没有区别。

```go
import "github.com/chhongzh/lilt/widget"
```

## 总览

| 组件 | 用途 | `OnChange` 值类型 | 页面 |
| --- | --- | --- | --- |
| `Label` | 一行文字，放不下时跑马灯 | 无 | [Label](label.md) |
| `List` | 竖向滚动的一列，放列表行 | 无 | [List 与 ListItem](list.md#list) |
| `ListItem` | 可聚焦的列表行：标题加可选的尾部组件，松开 ++ok++ 执行 `OnClick` | 无 | [List 与 ListItem](list.md#listitem) |
| `Button` | 带边框的按钮，可要求长按确认 | 无 | [Button](button.md) |
| `CheckBox` | 带勾选框的列表行 | `bool` | [CheckBox 与 Radio](toggle.md#checkbox) |
| `RadioGroup`、`Radio` | 一组单选行，`OnChange` 挂在组上 | `int`（下标） | [CheckBox 与 Radio](toggle.md#radio) |
| `Select` | 显示当前选项的列表行，激活后从底部弹出单选抽屉 | `int`（下标） | [Select](select.md) |
| `Chip` | `Choice`、`Toggles` 里的一个选项，单独用时只有 `OnClick` | 无 | [Choice 与 Toggles](chips.md) |
| `Choice` | 横向一排标签，只能开一个 | `int`（下标） | [Choice 与 Toggles](chips.md#choice) |
| `Toggles` | 横向一排标签，可以开任意个 | `uint32`（第 i 位是第 i 个） | [Choice 与 Toggles](chips.md#toggles) |
| `Input` | 单行输入框，激活后打开输入法 | `string` | [Input 与软键盘](input.md) |
| `SoftKeyboard`、`Keyboard` | 屏幕软键盘，`Input` 在没有设置 `ui.TextInput` 时用它 | 无 | [Input 与软键盘](input.md#keyboard) |
| `Scroll` | 裁剪内容并按焦点滚动，可横向、纵向或两向 | 无 | [Scroll](scroll.md) |
| `Drawer` | 从屏幕一边滑入的抽屉 | 无 | [弹窗](popups.md#drawer) |
| `Dialog` | 屏幕中间的对话框，带一排按钮 | 无 | [弹窗](popups.md#dialog) |
| `Toast` | 屏幕底部显示两秒的一行提示 | 无 | [弹窗](popups.md#toast) |

## 共同规则

所有可交互的组件都嵌入了 [`widget.Control`](../concepts/threading.md#control)，所以下面几条对它们都成立，各组件页不再重复：

- **松手才激活**：按下 ++ok++ 只让光标收缩，松开且没有被取消时才触发（[交互约定](../concepts/input.md#activation)）。
- **回调在自己的 goroutine 上跑**：`OnClick`、`OnChange` 都不在 UI goroutine 上执行，慢回调不会卡住画面。回调里改树要包在 [`ui.Do`](../concepts/threading.md#do) 里。
- **值变了才调 `OnChange`**：回调拿到旧值和新值 `(from, to)`。`CheckBox`、`Radio`、`Select` 这几种列表行接着还会调 `OnClick`，值没变（比如再选一次同一个选项）就只调 `OnClick`；`OnChange` 返回 `error` 时 `OnClick` 不再执行。`Choice`、`Toggles`、`Input` 的规则见各自页面。
- **忙时不响应**：回调返回之前组件忽略激活；超过 150ms 还没返回，显示跳动的加载点。
- **失败就抖**：回调返回 `error` 时组件左右抖一下；有值的组件同时把值改回去。值在激活时立即改变，回调失败才回滚，所以界面不用等回调。
- **代码改值不触发回调**：`SetChecked`、`Select`、`SetSelected`、`Choose`、`SetMask`、`Input.SetText` 这些方法只改显示，不调 `OnChange`。

`OnChange` 的参数类型必须和组件报告的值类型一致（见上表），否则设置时直接 panic。回调可以是 `func(from, to T)`，也可以是 `func(from, to T) error`；`OnClick` 同理，可以是 `func()` 或 `func() error`。

## 选项 { #options }

`widget` 的选项和 `ui` 的选项一样是 `ui.Option`，可以混在一起传给构造函数（[选项](../concepts/tree.md#options)）。每个选项只要求节点有它需要的方法，不关心节点是什么类型；节点没有这个方法时，选项在构造时 panic，信息是 `widget: option X does not apply to this node`。

| 选项 | 需要的方法（接口） | 适用的标准组件 |
| --- | --- | --- |
| `widget.Font(f)` | `SetFont`（`FontSetter`） | `Label`、`Button`、`ListItem` 及其派生、`Chip`、`Choice`、`Toggles`、`Input` |
| `widget.TextColor(c)` | `SetColor`（`Colorer`） | `Label` |
| `widget.TextAlign(x)` | `SetAlign`（`Aligner`） | `Label` |
| `widget.MarqueeMode(m)` | 必须是 `*Label` | `Label` |
| `widget.OnClick(fn)` | 嵌入 `Control` | 所有可交互组件 |
| `widget.OnChange(fn)` | 嵌入 `Control` | 有值的组件（上表第三列不为“无”的） |
| `widget.HoldToConfirm(d)` | 必须是 `*Button` | `Button` |
| `widget.Trailing(n)` | `SetTrailing`（`Trailer`） | `ListItem` 及其派生 |

“`ListItem` 及其派生”指 `ListItem`、`CheckBox`、`Radio`、`Select`，后三个嵌入了 `ListItem`，方法跟着继承下来。所以 `widget.Font` 用在这些行上改的是标题字体。

!!! note "组件里面的文字"
    `TextColor`、`TextAlign`、`MarqueeMode` 只对 `Label` 本身有效，用在 `Button` 上会 panic。要改按钮或列表行里的文字，先取出里面的 `Label`：`b.Label().SetColor(c)`、`row.Title().SetMarquee(widget.MarqueeAlways)`。

能力接口定义在 `widget` 包里，自定义组件实现了对应方法就能直接用这些选项（[自定义组件](../guides/custom-widget.md#interfaces)）：

```go
type Texter interface {
	Text() string
	SetText(s string)
}
type FontSetter interface{ SetFont(f *gfx.Font) }
type Colorer interface{ SetColor(c gfx.Color) }
type Aligner interface{ SetAlign(x float32) }
type Trailer interface{ SetTrailing(d ui.Node) }
```

`Texter` 目前没有对应的选项，`Label` 和 `Input` 实现了它，用来以统一的方式读写一行文字。
