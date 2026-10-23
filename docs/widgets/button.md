# Button

`widget.Button` 是带边框的按钮，松开 ++ok++ 时执行 `OnClick`。它自己管理加载、成功、失败三种状态，回调只管干活、返回 `error`。

```go
import "github.com/chhongzh/lilt/widget"
```

## 点击

![点击计数的按钮](../assets/demos/button-click.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_button.go:button-click"
```

按下 ++ok++ 只让光标收缩，**松手**才触发（[交互约定](../concepts/input.md#activation)）。回调跑在自己的 goroutine 上，所以改按钮文字要包在 `ui.Do` 里；`tap` 先声明再赋值，是为了在回调里引用按钮自己。

## 加载、成功与失败

回调返回之前，按钮不再响应激活。回调超过 150ms 还没返回，文字换成跳动的加载点；返回后亮一个勾。返回 `error` 时按钮左右抖一下并显示叉，不管等没等过。

![慢回调与失败的回调](../assets/demos/button-states.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_button.go:button-states"
```

| 回调用时 | 返回 | 表现 |
| --- | --- | --- |
| 少于 150ms | `nil` | 文字不变 |
| 150ms 以上 | `nil` | 加载点 → 勾，1.2 秒后换回文字 |
| 任意 | `error` | 抖动 → 叉，1.2 秒后换回文字 |

!!! tip "不用自己写状态机"
    Busy、Loading 的判断都在嵌入的 [`widget.Control`](../concepts/threading.md#control) 里，`Button`、`ListItem`、`CheckBox` 等所有可交互组件共用同一套规则。

## 长按确认

`widget.HoldToConfirm(d)` 让按钮必须按住 `d` 才触发：按住时光标从左往右被吃掉，吃完立即执行回调。中途松手，光标弹回原样，什么也不发生。适合删除、重置这类危险操作。

![长按删除](../assets/demos/button-hold.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_button.go:button-hold"
```

光标的形状由 `Button` 实现的 [`ui.CursorShaper`](../guides/custom-widget.md#cursorshaper) 决定，自定义组件也能用同样的办法改变光标。

## API

```go
func NewButton(text string, opts ...ui.Option) *Button

func (b *Button) Label() *Label                 // 按钮里的文字
func (b *Button) SetFont(f *gfx.Font)
func (b *Button) SetHold(d time.Duration)       // 同 HoldToConfirm，0 关闭
func (b *Button) OnClick[F ClickFunc](fn F)     // 来自嵌入的 Control
func (b *Button) Busy() bool
func (b *Button) Loading() bool
```

| 选项 | 作用 |
| --- | --- |
| `widget.OnClick(fn)` | 激活时执行，`fn` 是 `func()` 或 `func() error` |
| `widget.HoldToConfirm(d)` | 按住 `d` 才触发 |
| `widget.Font(f)` | 文字字体 |
| `ui.Width`、`ui.Grow` 等 | 任何 `ui.Option` 都可以用 |

默认样式取自主题：边框用 `Foreground`，圆角 `Radius`，内边距 `ButtonPadding`（见[主题与单位](../concepts/theme.md)）。在 `Row` 里并排放几个按钮时，给它们 `ui.Grow(1)` 平分宽度，[`Dialog`](popups.md#dialog) 就是这样排按钮的。
