# API 速查

每个包一页，按用途分组列出全部导出的标识符：签名、一句话说明，以及讲它的页面。完整的文档注释见各包在 pkg.go.dev 上的页面。

| 包 | 内容 | pkg.go.dev |
| --- | --- | --- |
| [`lilt`](lilt.md) | `Start`、`Pause`、`Stop` 和帧率上限 `FPS` | [lilt](https://pkg.go.dev/github.com/chhongzh/lilt) |
| [`ui`](ui.md) | 节点树、布局、层栈、焦点、输入、光标、主题 | [ui](https://pkg.go.dev/github.com/chhongzh/lilt/ui) |
| [`widget`](widget.md) | 标准组件：Label、List、Button、CheckBox、Select、Chip、Input、Drawer 等 | [widget](https://pkg.go.dev/github.com/chhongzh/lilt/widget) |
| [`anim`](anim.md) | 弹簧 `Spring`、动画值 `Value` 和调度器 | [anim](https://pkg.go.dev/github.com/chhongzh/lilt/anim) |
| [`gfx`](gfx.md) | `Screen` 接口、带裁剪和平移的 `Canvas`、字体 | [gfx](https://pkg.go.dev/github.com/chhongzh/lilt/gfx) |
| [`driver/silent`](../guides/drivers.md#silent) | 只画进内存的单色屏幕，导出 PNG 和 GIF | [driver/silent](https://pkg.go.dev/github.com/chhongzh/lilt/driver/silent) |
| [`driver/desktop`](../guides/drivers.md#desktop) | 桌面窗口和浏览器里的模拟器 | [driver/desktop](https://pkg.go.dev/github.com/chhongzh/lilt/driver/desktop) |

`driver/sh1106` 是示例固件用的 SH1106 驱动，依赖 `machine`，之后会重写，这里不列出。

## driver/silent 与 driver/desktop

这两个包的用法都在[驱动与模拟器](../guides/drivers.md)里讲，这里只列签名。

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func silent.New(w, h int) *silent.Screen` | 一块全灭的单色屏幕，实现 `gfx.Screen` 和 `gfx.Inverter` | [silent](../guides/drivers.md#silent) |
| `Screen` 的 `Size`、`SetPixel`、`FillRect`、`InvertRect`、`Flush` | 屏幕接口的实现；`Flush` 什么也不做 | [silent](../guides/drivers.md#silent) |
| `func (s *Screen) Image(scale int) *image.Paletted` | 整屏转成图像 | [silent](../guides/drivers.md#silent) |
| `func (s *Screen) Region(r image.Rectangle, scale int) *image.Paletted` | 只转 `r` 以内，保留位置 | [silent](../guides/drivers.md#silent) |
| `func (s *Screen) Pixels() []bool` | 缓冲区本身，逐行，`true` 为亮 | [silent](../guides/drivers.md#silent) |
| `func (s *Screen) Diff(prev []bool) image.Rectangle` | 和之前的副本相比变化部分的外接矩形 | [silent](../guides/drivers.md#silent) |
| `var silent.Palette` | 灭、亮两种颜色 | [silent](../guides/drivers.md#silent) |
| `type silent.Recorder struct{ Scale int }` | 把一帧帧画面攒成 GIF | [Recorder](../guides/drivers.md#recorder) |
| `func (r *Recorder) Capture(s *Screen, delay int)` | 记下当前画面，显示 `delay` 个 1/100 秒 | [Recorder](../guides/drivers.md#recorder) |
| `func (r *Recorder) Frames() int` | 已存的帧数 | [Recorder](../guides/drivers.md#recorder) |
| `func (r *Recorder) Encode(w io.Writer) error` | 写出 GIF | [Recorder](../guides/drivers.md#recorder) |
| `func desktop.Run(cfg desktop.Config, root ui.Node) error` | 在窗口里显示 `root`，阻塞到窗口关闭 | [desktop](../guides/drivers.md#desktop) |
| `type desktop.Config struct{ Title string; Width, Height, Scale int }` | 窗口设置，零值取默认 | [desktop](../guides/drivers.md#desktop) |
