# anim

[pkg.go.dev](https://pkg.go.dev/github.com/chhongzh/lilt/anim) · `import "github.com/chhongzh/lilt/anim"`

框架里所有会动的数值都由同一种弹簧模型驱动。组件里用 `Value`，全局手感调 `Spring` 预设。

## 弹簧

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Spring struct { Response, Damping float32; Ease bool }` | SwiftUI 式参数：`Response` 是无阻尼振荡周期（秒），`Damping` 是临界阻尼的比例，1 不回弹，越小越弹 | [弹簧](../concepts/animation.md#spring) |
| `Spring.Ease` | 改用指数缓动：第一帧就全速，没有惯性也不回弹；每 `Response/2π` 秒缩小剩余距离的 1-1/e，忽略 `Damping` | [弹簧](../concepts/animation.md#spring) |
| `var Default = &Spring{Response: 0.29, Damping: 0.75}` | 列表项和布局，轻微过冲；`Value` 的默认弹簧 | [预设](../concepts/animation.md#presets) |
| `var Smooth = &Spring{Response: 0.39, Damping: 0.85}` | 整屏的推入推出 | [预设](../concepts/animation.md#presets) |
| `var Snappy = &Spring{Response: 0.31, Damping: 0.8}` | 滑块、进度这类要紧跟的值 | [预设](../concepts/animation.md#presets) |
| `var Bouncy = &Spring{Response: 0.25, Damping: 0.45}` | 勾选标记等需要弹一下落定的东西 | [预设](../concepts/animation.md#presets) |
| `var Interactive = &Spring{Response: 0.14, Damping: 1}` | 跟手：按键按下时的反馈，快且不回弹 | [预设](../concepts/animation.md#presets) |

预设是指针，运行时改一处，所有用它的值一起变。

## 动画值

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Value struct` | 朝目标弹过去的 float32；零值可用，静止在 0，用 `Default`；开始动之后不能复制 | [Value](../concepts/animation.md#value) |
| `func NewValue(v float32, s *Spring) *Value` | 静止在 `v` 的值 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Set(x float32)` | 以弹簧动画移向 `x` | [Value](../concepts/animation.md#value) |
| `func (v *Value) Jump(x float32)` | 立即到 `x` 并停止运动 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Shift(d float32)` | 当前值移动 `d`、目标不变，于是滑回原处；布局动画用它 | [布局动画](../concepts/layout.md#animate-layout) |
| `func (v *Value) Kick(velocity float32)` | 加一个速度，例如列表到头时顶一下 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Bump(px float32)` | 冲量大小按弹簧算好，静止时大约荡出 `px` 再回来；缓动值忽略 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Get() float32` | 显示用的值，含父值 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Current() float32` | 当前值，不含父值 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Target() float32` | 目标值，不含父值 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Velocity() float32` | 速度，单位每秒 | [Value](../concepts/animation.md#value) |
| `func (v *Value) Settled() bool` | 是否已静止 | [Value](../concepts/animation.md#value) |
| `func (v *Value) SetSpring(s *Spring)` | 换弹簧，nil 表示 `Default` | [Value](../concepts/animation.md#value) |
| `func (v *Value) SetParent(p *Value)` | 让 `Get` 加上 `p` 的值，子值跟着父值动 | [Value](../concepts/animation.md#value) |
| `func (v *Value) SetPrecision(p float32)` | 离目标多近算静止，默认 0.05 适合像素，0 到 1 的值要更小 | [自定义组件](../guides/custom-widget.md#animation) |

## 调度器

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Step(dt float32) bool` | 把所有在动的值推进 `dt` 秒，报告是否有值在动；帧循环每帧调用，组件不用调 | [Value](../concepts/animation.md#value) |
| `func Busy() bool` | 是否还有值在动 | [Value](../concepts/animation.md#value) |
| `var TimeScale float32 = 1` | 所有动画的时间倍率，调小放慢观察，调大让动画更快结束 | [帧率与电子墨水屏](../guides/hardware.md#fps) |

## 橡皮筋

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func RubberBand(x, dim float32) float32` | 把原始越界距离换成显示距离，越拉越费劲，永远到不了 `dim` | [滚动到边缘](../widgets/scroll.md#edges) |
| `var RubberBandCoefficient float32 = 0.55` | 橡皮筋阻力，和 UIScrollView 相同 | [滚动到边缘](../widgets/scroll.md#edges) |
