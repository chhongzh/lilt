# 节点树

界面是一棵节点树：根下面每一层是一个页面或弹窗，页面里是容器和组件。树只属于 UI goroutine，`lilt.Start` 之前可以在任何地方建树，之后要改就得在 UI goroutine 上，或者包在 [`ui.Do`](threading.md#do) 里。

```go
import "github.com/chhongzh/lilt/ui"
```

## 节点与 Box

`ui.Node` 是一个只有未导出方法的接口，只有嵌入了 `ui.Box` 的类型才满足它。`Box` 装着一个节点共有的东西：父子链接、尺寸和位置等布局属性、背景和边框，以及叠加在布局位置上的动画偏移 `OffsetX()`、`OffsetY()`。组件额外会做什么（画自己、处理按键、测量文字），靠实现下面的[能力接口](#interfaces)，框架用类型断言找到它们。

所以 `widget.Label`、`ui.Flex`、`widget.List` 都是节点，可以互相嵌套。手里只有一个 `ui.Node` 时，用 `ui.BoxOf(n)` 拿到它的 `Box` 来读写公共属性：

```go
ui.BoxOf(n).SetGrow(1)
```

`Box` 上和树有关的只读方法：

| 方法 | 返回 |
| --- | --- |
| `Parent()` | 父节点，没挂上时为 `nil` |
| `Children()` | 子节点切片，只能读，增删要用下面的函数 |
| `Rect()` | 布局后的矩形，相对父节点左上角，不含动画偏移 |
| `ContentRect()` | 去掉内边距之后的区域，像素，相对节点自己 |

其余属性的 Set/Get 见[布局](layout.md)。

## 用选项建树 { #options }

每个构造函数的最后一个参数都是 `...ui.Option`。`Option` 就是 `func(ui.Node)`，构造函数先设好默认值，最后调用 `ui.Apply(节点, opts)` 依次执行它们，所以选项会覆盖默认值，例如给按钮传 `ui.Padding(...)` 会替换主题里的 `ButtonPadding`。

`ui.Children(...)` 也是一个选项，它把子节点加进正在构造的容器。这样一整屏可以写成一个嵌套的表达式，代码的缩进就是树的形状：

```go
screen := ui.NewColumn(ui.Fill(), ui.Children(
	widget.NewLabel("Settings"),
	widget.NewList(ui.Grow(1), ui.Children(
		widget.NewCheckBox("Wi-Fi", true),
		widget.NewCheckBox("Bluetooth", false),
	)),
))
```

`ui` 包的选项作用于任何节点，只有少数要求特定类型，用错类型会在构造时 panic：`ui.Gap` 要 `*ui.Flex` 或 `*ui.Grid`，`ui.Main`、`ui.Cross` 要 `*ui.Flex`，`ui.Columns` 要 `*ui.Grid`。`widget` 包的选项按能力匹配，例如 `widget.Font` 作用于任何有 `SetFont` 方法的组件，`widget.OnClick` 作用于任何嵌入了 `widget.Control` 的组件。

### 在嵌套里留住引用 { #ref }

嵌套表达式的缺点是里面的节点没有名字，回调里想改它就够不着。`ui.Ref(&p)` 在构造时把节点存进变量 `p`：

![按钮改动表头里的音量](../assets/demos/tree-ref.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_tree.go:tree-ref"
```

`ui.Ref` 是泛型的，`p` 的类型决定存进去的类型，这里是 `*widget.Label`，取出来就能直接调用 `SetText`。变量要在表达式之前声明。回调跑在自己的 goroutine 上，所以连 `level` 一起放进 `ui.Do`，两个按钮的回调就不会同时改它。

## 增删节点

| 函数 | 作用 |
| --- | --- |
| `ui.Add(parent, children...)` | 依次追加到末尾 |
| `ui.Insert(parent, i, child)` | 插到第 `i` 个位置，`i` 从 0 到现有子节点数 |
| `ui.Remove(n)` | 从父节点上摘下 `n`，`n` 没挂着时什么也不做 |
| `ui.RemoveChildren(parent)` | 从最后一个开始摘下所有子节点 |
| `ui.Attach(parent, child)` | 直接追加，不经过 `ContentHost` 转发，给组件搭内部结构用 |

![新建闹钟插到第二行，点一行删除它](../assets/demos/tree-add-remove.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_tree.go:tree-add-remove"
```

`row` 先声明再赋值，回调里才能引用这一行自己。节点可以在自己的回调里把自己删掉；焦点在被删的子树里时，会移到剩下的节点里离它最近的那个，上面删掉 `Alarm 2` 后光标落到了 `Alarm 3`。

几条规则：

- **一个节点只能有一个父节点。** 挂到第二个父节点上会 panic，先 `ui.Remove` 再挂。摘下的节点保留自己的状态，可以再加回树里。
- **叶子不收子节点。** 实现了 `ui.Leaf` 的节点（如 `widget.Label`）收到子节点会 panic。
- **`ContentHost` 会转发。** `ui.Add`、`ui.Insert`、`ui.RemoveChildren` 和 `ui.Children` 选项遇到实现了 `ui.ContentHost` 的节点，会改为操作它 `Content()` 返回的内部节点。`widget.List` 就是这样把行放进内部的那一列，滚动用的结构不暴露给调用方。`ui.Remove` 按节点本身操作，不受影响。
- **父节点会收到通知。** 子节点被摘下后，实现了 `ui.ChildRemover` 的父节点会收到 `ChildRemoved(child)`，子节点自己删自己时也一样。
- 每次增删都会调用 `ui.Relayout()`，下一帧整树重排。

!!! warning "回调里改树要包在 `ui.Do` 里"
    `OnClick`、`OnChange` 不在 UI goroutine 上跑，直接调用 `ui.Add`、`ui.Remove` 会和正在进行的帧冲突。包在 `ui.Do` 里，而且 `ui.Do` 里面不能再调用 `ui.Do`，见[线程模型](threading.md#do)。

### 判断位置

`ui.IsAncestor(a, n)` 报告 `a` 是不是 `n` 本身或它的祖先；`ui.Attached(n)` 报告 `n` 是否在当前显示的树里。慢回调返回时，所在的页面可能已经被关掉，先用 `ui.Attached` 判断一下再更新界面，可以避免改一个已经看不见的节点。

## 能力接口 { #interfaces }

节点的行为来自它实现的接口。框架每帧遍历树，用类型断言找到这些方法并调用；不实现就没有这项能力，也没有开销。

| 接口 | 方法 | 何时调用 |
| --- | --- | --- |
| `ui.Painter` | `Draw(c gfx.Canvas)` | 画自己，在子节点之前 |
| `ui.OverPainter` | `DrawOver(c gfx.Canvas)` | 在子节点之后再画一层 |
| `ui.Updater` | `Update(dt float32)` | 节点挂在树上且可见时，每帧一次 |
| `ui.Measurer` | `Measure(max Vec2) Vec2` | 叶子报告内容尺寸，不含内边距，见[布局](layout.md#size) |
| `ui.Layouter` | `MeasureChildren`、`ArrangeChildren` | 自己摆放子节点，`ui.Flex`、`ui.Grid` 就是这样实现的 |
| `ui.ContentHost` | `Content() Node` | 增删子节点时转发到内部节点 |
| `ui.Leaf` | `Leaf()` | 拒绝子节点 |
| `ui.ChildRemover` | `ChildRemoved(child Node)` | 子节点被摘下之后 |
| `ui.Focuser` | `FocusChanged(focused bool)` | 获得或失去焦点时，见[输入与焦点](input.md#focus) |
| `ui.KeyHandler` | `HandleKey(e *KeyEvent) bool` | 按键从焦点节点往根冒泡，见[按键分发](input.md#dispatch) |
| `ui.KeyPreviewer` | `PreviewKey(e *KeyEvent) bool` | 按键从根往焦点节点下行时先过一遍，容器可以截下来 |
| `ui.PointerHandler` | `HandlePointer(e *PointerEvent) bool` | 触摸或鼠标事件，从指针下的节点往根冒泡 |
| `ui.ArrowScroller` | `ArrowScroll(e *KeyEvent, next Node) bool` | 方向键移动焦点之前，滚动容器可以先消费它 |
| `ui.Revealer` | `Reveal(d Node)` | 子孙节点获得焦点时，滚动容器把它滚进视野 |
| `ui.CursorShaper` | `CursorShape(r Rect) Rect` | 节点聚焦时改变光标的矩形，见[自定义组件](../guides/custom-widget.md#cursorshaper) |

一个节点画出来的顺序是：`Box` 的背景和边框 → `Draw` → 子节点 → `DrawOver`。隐藏的节点和它的子树既不画，也不调用 `Update`。

怎样组合这些接口写一个新组件，见[自定义组件](../guides/custom-widget.md#interfaces)。

!!! warning "嵌入 `ui.Box` 时避开这些方法名"
    TinyGo 0.42 和 0.43 有个编译器 bug：嵌入 `ui.Box` 或 `ui.Flex` 的类型如果声明了和它们未导出方法同名的方法（`box`、`set`、`visible`、`split`、`join`、`margins`、`crossAuto`），转成接口时编译器会崩溃，报 `index out of range [0] with length 0`。自己的组件里换个名字。
