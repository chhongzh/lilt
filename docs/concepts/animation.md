# 动画

lilt 里所有会动的东西都由同一种阻尼弹簧驱动：列表行的滑动、光标、页面飞入、勾选框的方框。动画的基本单位是 `anim.Value`，一个朝目标值弹过去的 `float32`。

```go
import "github.com/chhongzh/lilt/anim"
```

## anim.Value { #value }

`Value` 的零值就能用：停在 0，使用 `anim.Default` 弹簧。设一个目标，它就自己动过去；绘制时读 `Get()`。

```go
var x anim.Value      // 停在 0
x.SetSpring(anim.Bouncy)
x.Set(40)             // 从当前位置弹向 40
x.Get()               // 这一帧该画在哪里
```

下面的自定义节点把一个 `Value` 画成线上的一个点，0 在左端，1 在右端：

```go
--8<-- "cmd/docgif/demo_animation.go:animation-lane"
```

`Draw` 里只管读 `Get()`，不需要自己调用 `ui.Invalidate()`：只要有 `Value` 还在动，帧循环就会重画。`Value` 嵌在结构体里，结构体用指针传递，见下面的警告。

### 方法

| 方法 | 作用 |
| --- | --- |
| `Set(x)` | 以当前位置和速度为起点，弹向 `x`。目标没变时什么也不做 |
| `Jump(x)` | 立刻停在 `x`，速度清零 |
| `Shift(d)` | 当前值挪 `d`，目标不变，于是再滑回来。布局的[滑动](layout.md#animate-layout)就是这样做的 |
| `Kick(v)` | 速度加上 `v`（单位每秒），例如组件失败时的抖动 |
| `Bump(px)` | 按弹簧换算一个初速度，让静止的值大约荡开 `px` 再回来，比 `Kick` 好估计幅度 |
| `Get()` | 显示值，包括父值 |
| `Current()` | 不含父值的当前值 |
| `Target()` | 目标值，不含父值 |
| `Velocity()` | 速度，单位每秒 |
| `Settled()` | 是否已经静止 |
| `SetSpring(s)` | 换弹簧，`nil` 表示 `anim.Default` |
| `SetParent(p)` | 让 `Get()` 加上 `p.Get()`，在父值的运动上叠加自己的 |
| `SetPrecision(p)` | 距离目标多近算静止，默认 0.05 |

`anim.NewValue(v, s)` 返回一个停在 `v`、使用弹簧 `s` 的 `*Value`。

精度的默认值 0.05 是按像素定的。做 0 到 1 的进度时，0.05 就是 5%，值会在离终点还差一截时突然跳过去，所以上面的 `lane` 把精度设成 0.002，标准组件里勾选框的方框、按钮的长按进度也都设了更小的精度。值离目标小于精度、速度小于精度的 10 倍时，它直接落在目标上并停下。

`Bump` 只对弹簧有效；`Ease` 曲线（见下文）没有动量，`Kick` 和 `Bump` 加进去的速度会被忽略。

!!! warning "动起来之后不要复制 Value"
    调度器持有每个运动中 `Value` 的指针。把一个正在动的 `Value` 按值复制，副本不会动，原件还在被推进。把 `Value` 放在用指针传递的结构体里，像 `ui.Box` 和所有组件那样。`Value` 也没有加锁，只在 UI goroutine 上或 `ui.Do` 里改它。

## 调度器

`anim` 包有一个全局的运动列表。`Value` 的目标或速度一变，它就加进列表；静止后自己离开。所以静止的值不花任何时间，也不需要注册和释放。

```go
func Step(dt float32) bool // 推进所有运动中的值，返回是否有值在动
func Busy() bool           // 是否还有值在动
var TimeScale float32 = 1  // 所有动画的时间倍率
```

UI goroutine 每帧调用一次 `anim.Step(dt)`，返回 `true` 时那一帧重画。自己的代码不需要调用它。`TimeScale` 设成 0.2 可以把所有动画放慢 5 倍，调手感时用来看清楚细节。

弹簧用半隐式欧拉积分，每帧再细分成 1/240 秒的小步，帧率低、弹簧硬的时候也不会发散。因为动画按时间增量计算，帧率不同时运动的快慢一样，只是采样密度不同。

## 弹簧 { #spring }

```go
type Spring struct {
	Response float32 // 无阻尼振动的周期，秒
	Damping  float32 // 临界阻尼的比例：1 不过冲，越小越弹
	Ease     bool    // 改用指数缓动：没有动量，不过冲，忽略 Damping
}
```

参数采用 SwiftUI 的写法，比刚度和阻尼系数好调：`Response` 决定快慢，`Damping` 决定弹不弹，两者互不影响。`Damping` 为 1 时临界阻尼，最快到达且不越过目标；0.5 左右会明显回弹一次。

`Ease` 为真时不再是弹簧：值第一帧就全速出发，每过 `Response/2π` 秒缩小剩余距离的 1−1/e，没有惯性也不会越过。光标的边用它，换焦点时立即响应，不会甩出去。其中的 e^−x 用 float32 有理逼近计算，因为 ESP32-S3 的 FPU 只做单精度，`math.Exp` 的 float64 要靠软件模拟。

## 预设 { #presets }

![五个预设弹向同一个目标](../assets/demos/animation-presets.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_animation.go:animation-presets"
```

五个点在同一帧拿到同一个目标，区别只在弹簧。`Interactive` 最先到且不过冲，`Bouncy` 冲过终点又弹回，`Smooth` 最慢。改目标的代码在 `OnClick` 里，所以包在 `ui.Do` 里（见[线程模型](threading.md#do)）。

| 预设 | Response | Damping | 用在 |
| --- | --- | --- | --- |
| `anim.Default` | 0.29 | 0.75 | 滚动时行的波浪（每次滚动从它派生，离焦点越远越慢）、对话框落下，也是 `Value` 的默认弹簧 |
| `anim.Smooth` | 0.39 | 0.85 | 页面、抽屉、Toast 的进出 |
| `anim.Snappy` | 0.31 | 0.8 | 需要紧跟的东西：滚动条、长按进度、按钮里文字和勾叉的上下切换、行尾控件的进出 |
| `anim.Bouncy` | 0.25 | 0.45 | 勾选框和单选框的标记，按钮、行、输入框、Chip 行失败时的抖动 |
| `anim.Interactive` | 0.14 | 1 | 跟随按下的键或手指：光标按下时的收缩 |

除了 `Interactive`，每个预设都带一点过冲，这是有意的。

预设是指针，`Value` 保存的也是指针，所以运行时改预设的字段，所有用它的值从下一帧起立即生效。这就是全局调整手感的方式：

```go
anim.Bouncy.Damping = 0.3 // 所有勾选框弹得更厉害
anim.Smooth.Response = 0.3 // 页面切换更快
```

框架里还有几个专用的弹簧，同样是可以改的指针：光标的 `ui.CursorLead`、`CursorTrail`、`CursorRelease`、`CursorLeap`、`CursorLean`（见[光标](input.md#cursor)），Chip 横条的 `widget.ChipSpring`。

## 节点的偏移

每个节点的 `Box` 带两个 `Value`：`OffsetX()` 和 `OffsetY()`，叠加在布局给出的位置上。子节点相对父节点绘制，所以偏移会带着整棵子树一起动，页面推入推出、弹窗滑入、组件抖动都是改这两个值。偏移不影响布局，也不影响焦点查找（焦点用的是[落定后的位置](input.md#focus)）。

![给标签一个速度和一次 Bump](../assets/demos/animation-offset.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_animation.go:animation-offset"
```

`Kick` 给的是速度，幅度取决于弹簧；`Bump` 给的是大致的幅度。标准组件失败时的抖动就是 `OffsetX().Kick(-260)`，配合组件自己设的 `anim.Bouncy`。

## 回弹阻力

```go
func RubberBand(x, dim float32) float32
var RubberBandCoefficient float32 = 0.55
```

`RubberBand` 把越界的距离 `x` 换算成显示出来的距离，结果趋近但永远到不了 `dim`，推得越多动得越少。系数 0.55 和 `UIScrollView` 的越界阻力相同。滚动容器拖过头时用它。
