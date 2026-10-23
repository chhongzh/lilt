# 驱动与模拟器

没有板子时，界面照样能跑、能看、能录。仓库带两个主机上用的驱动和三个命令：

| 名字 | 是什么 | 用途 |
| --- | --- | --- |
| `driver/silent` | 只画进内存的单色屏幕 | 测试、截图、录 GIF |
| `driver/desktop` | 桌面窗口，也能编成 wasm | 边写边调，浏览器里试玩 |
| `cmd/sim` | 在窗口或浏览器里跑 gallery | 看所有组件的实际手感 |
| `cmd/preview` | 按脚本回放 gallery，录 GIF 和截图 | README 的动图 |
| `cmd/docgif` | 录这个站点的全部动图 | 文档 |

## silent：内存里的屏幕 { #silent }

`silent.New(w, h)` 返回一块单色帧缓冲，它同时实现了 `gfx.Screen` 和 `gfx.Inverter`，所以光标和 OLED 上一样是反色的。颜色按 `gfx.IsLit` 变成亮或灭，`Flush` 什么也不做。

```go
package main

import (
	"image/png"
	"log"
	"os"
	"time"

	"github.com/chhongzh/lilt"
	"github.com/chhongzh/lilt/driver/silent"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func main() {
	list := widget.NewList()
	ui.Add(list, widget.NewCheckBox("Wi-Fi", true), widget.NewCheckBox("Bluetooth", false))

	scr := silent.New(128, 64)
	lilt.Start(scr, list)
	time.Sleep(time.Second) // 让光标长出来、动画停下
	if err := lilt.Stop(); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("shot.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, scr.Image(4)); err != nil { // 每个屏幕像素放大成 4x4
		log.Fatal(err)
	}
}
```

| API | 作用 |
| --- | --- |
| `silent.New(w, h)` | 一块 `w`×`h` 的屏幕，全灭 |
| `s.Image(scale)` | 整屏转成 `*image.Paletted`，每个像素放大 `scale` 倍 |
| `s.Region(r, scale)` | 只转 `r` 以内的像素；图像保留它在放大后的屏幕里的位置，GIF 帧可以只带变化的那块 |
| `s.Pixels()` | 缓冲区本身，逐行排列，`true` 是亮；测试里直接断言像素 |
| `s.Diff(prev)` | 和之前从 `Pixels` 复制的一份相比，变化部分的外接矩形 |
| `silent.Palette` | 灭和亮两种颜色，依次排列；改它就改了所有导出图像的配色 |

!!! warning "不能边跑边读"
    `silent.Screen` 不是并发安全的。UI goroutine 在画的时候，别的 goroutine 读 `Pixels`、调用 `Image` 都是数据竞争。要么像上面那样 `lilt.Stop` 之后再读（缓冲区里留着最后一帧），要么先 `lilt.Pause`，读完再 `lilt.Start` 恢复。

## Recorder：录 GIF { #recorder }

`silent.Recorder` 把屏幕的一帧帧画面攒成动画 GIF：

| API | 作用 |
| --- | --- |
| `Recorder{Scale: n}` | 每个屏幕像素在 GIF 里占 `n`×`n`，0 当作 1 |
| `r.Capture(s, delay)` | 记下屏幕此刻的画面，显示 `delay` 个 1/100 秒 |
| `r.Frames()` | 已经存了多少帧 |
| `r.Encode(w)` | 把整段 GIF 写进 `w` |

每帧只存和上一帧不同的那块矩形，画在上一帧之上；和上一帧完全相同时不新增帧，只把上一帧的显示时间加长。所以静止的画面几乎不占体积，录长一点也没关系。

`Capture` 读屏幕，和截图一样不能在 UI 运行时调用。在这个仓库以外，可以用 `lilt.Pause` 和 `lilt.Start` 把帧循环停一下再采样。接着上面的例子：

```go
scr := silent.New(128, 64)
rec := &silent.Recorder{Scale: 3}
lilt.Start(scr, list)
// One GIF frame every 20ms: browsers slow down shorter GIF delays.
for i := range 100 { // two seconds
	switch i {
	case 25:
		ui.Emit(ui.KeyDown, ui.Press)
	case 30:
		ui.Emit(ui.KeyDown, ui.Release)
	}
	time.Sleep(20 * time.Millisecond)
	lilt.Pause() // returns once the frame loop is idle
	rec.Capture(scr, 2)
	lilt.Start(nil, nil) // resumes; the arguments are ignored while paused
}
if err := lilt.Stop(); err != nil {
	log.Fatal(err)
}
f, err := os.Create("demo.gif")
if err != nil {
	log.Fatal(err)
}
defer f.Close()
if err := rec.Encode(f); err != nil {
	log.Fatal(err)
}
```

`Pause` 要等帧循环收到命令才返回，此后它一直停着，读屏幕是安全的；暂停期间来的按键在队列里等着。恢复时帧循环重新计时，暂停的那段时间不算进动画。这种做法按真实时间采样，帧与帧的间隔会有几毫秒的误差，做文档和演示足够了。

仓库里的 `cmd/preview` 和 `cmd/docgif` 录得更精确：它们通过 `internal/reel` 不经 `lilt.Start`，自己一帧一帧推进 UI，每一帧都落在固定的时间点上，同时按真实时间放慢节奏，让在自己 goroutine 上 sleep 的回调和画面对得上。`internal/` 下的包按 Go 的规则只能被本模块的代码 import，你自己的项目用不了，用上面的办法即可。

## desktop：桌面窗口 { #desktop }

`desktop.Run` 开一个窗口当屏幕，键盘当按键。它基于 [Ebitengine](https://ebitengine.org)，不需要 cgo。

```go
package main

import (
	"log"

	"github.com/chhongzh/lilt/driver/desktop"
)

func main() {
	// Run 占住主 goroutine，直到窗口关闭。
	if err := desktop.Run(desktop.Config{Title: "hello", Scale: 6}, hello()); err != nil {
		log.Fatal(err)
	}
}
```

| `Config` 字段 | 默认值 | 作用 |
| --- | --- | --- |
| `Title` | `"lilt"` | 窗口标题 |
| `Width`、`Height` | 128、64 | 屏幕像素数，和 SH1106 一样 |
| `Scale` | 6 | 打开窗口时一个屏幕像素占多大 |

零值字段取默认值。窗口可以拖动改变大小，画面按整数倍放大并居中，像素边缘保持清晰。

| 键盘 | 对应 |
| --- | --- |
| ++w++ ++a++ ++s++ ++d++ 或方向键 | 上、左、下、右 |
| ++space++、++enter++ | OK（摇杆按下） |
| ++j++ | Back（板子上的 SET 键） |
| ++k++ | Menu（板子上的 RST 键） |

几点行为：

- **`Run` 内部调用 `lilt.Start`**，不要自己再调；窗口关闭时它调用 `lilt.Stop`，返回 Ebitengine 或 `Stop` 的错误。所以它之后的代码只在窗口关掉后才执行，后台任务要在调用 `Run` 之前用 goroutine 启动，改树照样用 `ui.Do`。
- **必须在 `main` 里调用**：窗口系统要占用主 goroutine。
- **屏幕是 `silent.Screen`**，单色、能反色，颜色用 `silent.Palette`。在这里调彩色主题只能看到亮灭，彩色效果要上真机。
- **按键经过 `ui.KeyPoller`**，按住同样会加速连发；Ebitengine 每次 update 采样一次键盘，本身没有抖动，所以消抖关掉了（`Debounce = 0`）。

### 编成 wasm

同一份程序设 `GOOS=js GOARCH=wasm` 就能在浏览器里跑：

```bash
mkdir out
GOOS=js GOARCH=wasm go build -o out/lilt.wasm ./cmd/sim
cp cmd/sim/index.html "$(go env GOROOT)/lib/wasm/wasm_exec.js" out/
python -m http.server -d out
```

`cmd/sim/index.html` 是加载 `lilt.wasm` 的页面，换成自己的程序时把 `./cmd/sim` 换成你的包。[在线体验](../play.md)页面嵌的就是这样编出来的 gallery。

## cmd/sim 与 cmd/preview

`cmd/sim` 就是用 `desktop.Run` 跑 `examples/gallery`：

```bash
go run ./cmd/sim
```

`cmd/preview` 在主机上按一段写死的按键脚本操作 gallery，把整段录成 GIF，并把标记的关键时刻拼成一张图。改了组件以后跑一遍，看 `shots.png` 就能检查有没有画坏：

```bash
go run ./cmd/preview -o /tmp/pv              # session.gif 和 shots.png
go run ./cmd/preview -o /tmp/pv -highlight   # 半分钟的精华片段 highlight.gif
```

| 参数 | 默认值 | 作用 |
| --- | --- | --- |
| `-o` | `preview` | 输出目录 |
| `-scale` | 4 | GIF 里一个屏幕像素的大小 |
| `-highlight` | 关 | 改录 README 开头那段精华 |

## 这个站点的动图

站点里的每张动图都由 `cmd/docgif` 录制，页面上的示例代码就是录动图的那份代码。每个页面对应一个 `cmd/docgif/demo_<页面名>.go`，在 `init` 里注册 demo：名字、构造界面的 `build`、按键脚本 `script`。这是[自定义组件](custom-widget.md)那张滑块动图的注册代码：

```go
--8<-- "cmd/docgif/demo_custom.go:drivers-docgif"
```

脚本按时间顺序写，单位是秒：`Wait` 推进时间，`Emit` 发一个事件，`Tap` 按下松开（按住 0.08 秒），`Taps` 连按几次，`Hold` 按住一段时间再松开（中间没有连发），`Shot` 标记一张静帧。`build` 返回的根节点上，`// --8<-- [start:名字]` 和 `// --8<-- [end:名字]` 之间的代码会原样嵌进页面。

录制前后 `docgif` 会保存并恢复 `ui.CurrentTheme` 和 `anim` 的弹簧预设，所以 demo 可以随意改主题。回调在自己的 goroutine 上按真实时间 sleep，脚本末尾要留够时间让它们返回，否则会作用到下一个 demo 上。

```bash
go run ./cmd/docgif -list                     # 列出所有 demo
go run ./cmd/docgif 'custom-*'                # 只录匹配的，写到 docs/assets/demos
go run ./cmd/docgif -shots /tmp/shots button-click   # 同时把静帧拼成 PNG
```

`scripts/docs.sh` 把录制和 MkDocs 串起来：

| 子命令 | 作用 |
| --- | --- |
| `gif [模式...]` | 录 GIF 到 `docs/assets/demos` |
| `shots [模式...]` | 录 GIF，并把静帧写到 `.docshots/` 供检查 |
| `serve` | 缺 GIF 时先补录，再在 `http://127.0.0.1:8000/lilt/` 实时预览 |
| `check` | `go vet` demo，补录缺失的 GIF，严格模式构建站点 |
| `build` | 重录全部 GIF，严格模式构建，再把 wasm 模拟器放进 `site/sim/` |

GIF 和 `site/` 不进 git。推到 `main` 后，GitHub Actions 运行 `scripts/docs.sh build` 重新录制、构建并发布。
