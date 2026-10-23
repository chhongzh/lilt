# Select

`widget.Select` 是显示当前选项的列表行。松开 ++ok++ 时从屏幕底部滑出一个抽屉，里面每个选项是一行 [`Radio`](toggle.md#radio)；选中一项，抽屉关闭，行上换成新选项。适合选项太多、一行 [`Choice`](chips.md#choice) 放不下的设置。

```go
import "github.com/chhongzh/lilt/widget"
```

## 选择

![打开 Speed，移到 Ludicrous 并选中](../assets/demos/select-pick.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_select.go:select-pick"
```

`NewSelect` 的参数依次是标题、选项文字和初始选中的下标。行尾用主题的 `SmallFont` 显示当前选项；下标越界（例如 `-1`）时行尾为空。

打开抽屉时：

- 抽屉的标题就是这一行的标题；
- 焦点落在当前选中的那一项上；
- 抽屉的高度按选项行数算，最多到屏幕高度的 80%，选项更多时抽屉里的列表可以滚动。

在抽屉里移动焦点不会改变选择，松开 ++ok++ 才算选中。按返回键 `ui.KeyBack` 关闭抽屉（[返回](../concepts/layers.md#back)），什么也不改，回调也不执行。抽屉每次打开时按选项重新建一遍，关闭后就丢掉。

选中一项后，抽屉关闭，行尾立即换成新选项，然后在自己的 goroutine 上依次执行：

1. `OnChange(from, to)`，参数是旧下标和新下标，仅当下标变了；
2. `OnClick`。

所以 `OnChange` 的回调类型是 `func(from, to int)` 或 `func(from, to int) error`。例子里在回调中弹 [`Toast`](popups.md#toast)，要包在 `ui.Do` 里。

## 失败与回滚

![回调失败，Region 回到 Europe](../assets/demos/select-fail.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_select.go:select-fail"
```

回调运行期间，加载状态显示在这一行上，而不是抽屉里：超过 150ms 时行尾的选项滑出，换成加载点，这期间再按 ++ok++ 不会打开抽屉。回调返回 `error` 时，行尾换回旧选项，整行抖一下。

| 回调返回 | 行尾 |
| --- | --- |
| `nil` | 保持新选项 |
| `error` | 换回旧选项，整行抖一下 |

`SetSelected(i)` 从代码里换选项，不触发 `OnChange`：

```go
ui.Do(func() { region.SetSelected(2) })
```

选项列表在构造时传入，之后不能修改。

## API

```go
func NewSelect(title string, options []string, selected int, opts ...ui.Option) *Select

func (s *Select) Selected() int         // 当前下标
func (s *Select) SetSelected(i int)     // 不触发 OnChange；越界时行尾为空
```

`Select` 嵌入了 [`ListItem`](list.md#listitem)，`Title`、`SetFont`、`OnClick`、`OnChange`、`Busy`、`Loading` 等方法都可以用。

| 选项 | 作用 |
| --- | --- |
| `widget.OnChange(fn)` | 选中了不同的选项后执行，`fn` 是 `func(from, to int)` 或带 `error` 返回值的版本 |
| `widget.OnClick(fn)` | 每次在抽屉里选中一项后执行，在 `OnChange` 之后 |
| `widget.Font(f)` | 标题字体；行尾的当前选项固定用 `SmallFont` |
| `ui.Width`、`ui.Grow` 等 | 任何 `ui.Option` 都可以用 |

默认样式取自主题：行的内边距 `RowPadding`，标题字体 `Font`，当前选项 `SmallFont`（见[主题与单位](../concepts/theme.md#theme)）。抽屉的样式见 [Drawer](popups.md#drawer)。
