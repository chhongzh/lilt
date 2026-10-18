# lilt

给 TinyGo 用的 GUI 框架，手感参照 osu!framework 和 iOS：所有会动的东西都由弹簧驱动，按下有反馈，松手才触发。框架与硬件无关，屏幕和按键都是接口，单色、彩色屏都能接。模块路径是 `github.com/chhongzh/lilt`。

设计需求见 [DESIGN.md](DESIGN.md)。

## 包

```
anim      弹簧动画 Value 与全局调度器
gfx       Screen 接口、带裁剪和平移的 Canvas、字体
ui        节点树、布局（Flex / Grid / Stack）、App 主循环、屏幕栈、焦点、输入、光标、主题
widget    Label、List、Scroll、Button、CheckBox、Radio、Select、Input、软键盘
driver    SH1106 驱动（待重写）
examples  gallery 演示，以及跑 gallery 的 ESP32-S3 固件（esp32s3-generic）
cmd       preview：在主机上按脚本回放一段操作，输出 GIF 和截图拼图
```

`anim`、`gfx`、`ui`、`widget` 不 import `machine`，主机上可以直接编译和预览。

## 接入一块屏幕

驱动实现 `gfx.Screen` 即可；能反色的屏再实现 `gfx.Inverter`，光标会改成反色绘制。现成的 tinygo driver 可以用 `gfx.FromDisplayer` 包一层，但它的 `FillRect` 退化成逐点画，最好自己写适配器（参考 [examples/esp32s3-generic/main.go](examples/esp32s3-generic/main.go)）。

```go
app := ui.NewApp(screen)
ui.Push(page, ui.SlideLeft)
app.Run() // 或者在自己的循环里调用 app.Frame(dt)
```

按键由调用方自己读引脚，再交给框架：

- `ui.KeyPoller` 负责消抖和按住加速重复，每帧对每个键调用一次 `Poll`；
- 也可以直接 `app.Input.Emit(key, phase)`，任何 goroutine 都能调用；
- 触摸屏用 `EmitPointer`，物理键盘用 `EmitRune`。

## 交互约定

- 按键事件分 `Press`、`Repeat`、`Release`。按下只让光标收缩，所有动作都在**松手**且未被取消（`!Canceled`）时触发。触摸变成拖动后，那次松手会带上 `Canceled`。
- `Button` 支持 `HoldToConfirm(d)`：按住时光标从右往左缩短，缩到头才执行。`OnAction` 在后台跑，按钮显示加载点，成功显示勾，失败显示叉并抖一下。
- 列表滚到顶或底还继续按，整个列表会往那个方向被推一下，松开弹回。
- 屏幕栈有两种转场：`SlideLeft`（新屏从右侧滑入，旧屏退到左边 1/3 处）和 `Sheet`（从底部升起，高度按内容算）。被不透明屏完全盖住的屏会被隐藏，不参与更新和绘制。
- `Input` 松手后打开 `App.TextInput`；没设置时使用 `widget.Keyboard` 软键盘。

## 渲染与布局

- 任何属性变化都会标记重新布局，下一帧整树重排。
- 有变化时整屏重绘，不做脏矩形。
- 每帧不分配内存。
- 只用 float32，不依赖 `strings` 包。
- 弹簧参数采用 SwiftUI 的 Response / Damping 写法。预设（`anim.Default`、`Smooth`、`Snappy`、`Bouncy`）是指针，运行时改一处，所有用它的值一起变。

## 构建与预览

```sh
gofmt -l anim gfx ui widget examples cmd
go vet ./anim ./gfx ./ui ./widget ./examples/gallery ./cmd/preview
go run ./cmd/preview -o /tmp/preview
tinygo build -target esp32s3-generic -o /tmp/gallery.bin ./examples/esp32s3-generic
```

`go build ./...` 会因为 `driver/sh1106` 依赖 `machine` 而失败，主机上请按包名单独编译。

preview 输出两个文件：`session.gif` 是整段操作的录像，`shots.png` 是关键帧拼图。

## 示例固件

[examples/esp32s3-generic](examples/esp32s3-generic) 是一台 ESP32-S3 + SH1106 + 五向摇杆的掌机，额外带一个实时显示引脚电平的按键页。材料、接线和烧录步骤见它的 [README](examples/esp32s3-generic/README.md)。
