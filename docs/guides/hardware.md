# 接入硬件

lilt 本身不读引脚、不碰总线：除了 `driver/sh1106`，所有包都不 import `machine`。板子上的代码负责两件事：给出一块实现了 `gfx.Screen` 的屏幕，把按键（或触摸、键盘）事件交给 `ui`。这一页讲这两件事怎么做。

## 屏幕接口 { #screen }

```go
type Screen interface {
	Size() (w, h int16)
	SetPixel(x, y int16, c Color)
	FillRect(x, y, w, h int16, c Color)
	// Flush pushes the finished frame to the panel.
	Flush() error
}
```

| 方法 | 框架怎么用它 |
| --- | --- |
| `Size` | `lilt.Start` 时读一次，作为根节点和每个页面的尺寸 |
| `SetPixel` | 画文字、斜线、圆角的边缘 |
| `FillRect` | 画其余的一切：背景、横竖线、圆角矩形的每一行 |
| `Flush` | 一帧画完后调用一次，把画面推到屏上；返回的错误会结束帧循环，由 `lilt.Stop` 返回 |

要点：

- **坐标超出屏幕时忽略，不要报错**。`Canvas` 已经裁剪过，但接口约定驱动自己也要容忍。
- **有变化才画**。没有动画、没有输入、没有改树的帧不画也不调用 `Flush`，静止的界面不占总线。
- **每次都整屏重画**。每帧先用主题背景色 `FillRect` 整屏，再画所有可见的层和光标，不做脏矩形。所以 `FillRect` 要快，`Flush` 最好一次推整块缓冲。
- **颜色是 `gfx.Color`，也就是 `color.RGBA`**。单色屏用 `gfx.IsLit(c)` 判断是否点亮（按 3:6:1 加权的亮度不低于一半），这样主题的 `Muted` 灰色在单色屏上也有确定的结果。

## tinygo 驱动：`gfx.FromDisplayer`

`gfx.Displayer` 的方法签名和 `tinygo.org/x/drivers` 的 `Displayer` 一样，实现了它的驱动可以直接包一层：

```go
dev := ssd1306.NewI2C(machine.I2C0)
dev.Configure(ssd1306.Config{Width: 128, Height: 64})
lilt.Start(gfx.FromDisplayer(dev), home)
```

`FromDisplayer` 把 `Flush` 接到驱动的 `Display`。驱动有 `FillRectangle` 方法时用它填充，没有就逐点 `SetPixel`。单色驱动自己决定什么颜色算亮，例如 ssd1306 把任何非纯黑的颜色都当作亮。

> **以下均未测试，但理论上可运行。** 实际测过的只有 [examples/esp32s3-generic](https://github.com/chhongzh/lilt/tree/main/examples/esp32s3-generic) 里的 SH1106，用的是本仓库自己的驱动。

**带帧缓冲的单色屏**：最合适。框架画进驱动的缓冲区，`Display` 一次推整屏。

| 驱动 | 屏幕 | 原生填充 |
| --- | --- | --- |
| `ssd1306` | SSD1306 OLED，128x64 / 128x32，I²C 或 SPI | 有 |
| `sh1106` | SH1106 OLED，128x64 | |
| `adafruit4650` | Adafruit FeatherWing OLED 128x64（SH1107） | |
| `pcd8544` | Nokia 5110 LCD，84x48 | |
| `sharpmem` | Sharp Memory LCD | |

**彩色屏**：主题默认黑白，改 `ui.CurrentTheme` 的颜色即可。这些驱动的 `SetPixel` 直接写屏、`Display` 什么也不做。框架每次整屏重绘，所以画面会闪，速度也受总线限制。要流畅，建议自己在内存里放一块帧缓冲，再写一个适配器整块推送，见[彩色屏与帧缓冲](#color)。

| 驱动 | 屏幕 | 原生填充 |
| --- | --- | --- |
| `st7789` | ST7789 TFT，240x240 / 240x320 等 | 有 |
| `st7735` | ST7735 TFT，128x160 / 80x160 | 有 |
| `ili9341` | ILI9341 TFT，240x320 | 有 |
| `gc9a01` | GC9A01 圆形 TFT，240x240 | 有 |
| `ssd1331` | SSD1331 彩色 OLED，96x64 | 有 |
| `ssd1351` | SSD1351 彩色 OLED，128x128 | 有 |
| `ssd1289` | SSD1289 TFT，240x320，并口 | |
| `hub75` | HUB75 RGB LED 点阵 | |

**电子墨水屏**：带缓冲，但每次 `Display` 都要刷新整屏，耗时可达数秒，弹簧动画基本看不到。适合静态菜单，见[帧率与电子墨水屏](#fps)。

| 驱动 | 屏幕 | 原生填充 |
| --- | --- | --- |
| `uc8151` | UC8151 墨水屏（如 Pimoroni Badger 2040） | 有 |
| `gdew0154m09` | GoodDisplay 1.54 寸墨水屏 | |
| `waveshare-epd/epd1in54` | 微雪 1.54 寸 | |
| `waveshare-epd/epd2in13` | 微雪 2.13 寸 | |
| `waveshare-epd/epd2in13x` | 微雪 2.13 寸三色 | |
| `waveshare-epd/epd2in66b` | 微雪 2.66 寸三色 | |
| `waveshare-epd/epd2in9` | 微雪 2.9 寸 | |
| `waveshare-epd/epd2in9v2` | 微雪 2.9 寸 V2 | |
| `waveshare-epd/epd4in2` | 微雪 4.2 寸 | |

`microbitmatrix`（micro:bit 5x5）和 `unoqmatrix`（Arduino UNO Q 点阵）也实现了这个接口，但分辨率连一行字都放不下，不适合。

这些驱动都没有 `InvertRect`，经 `FromDisplayer` 包装后光标不会反色，见下一节。

!!! warning "ssd1306 的原生填充会多推一帧"
    ssd1306 驱动的 `FillRectangle` 在被要求用黑色填满整屏时，会调用 `ClearDisplay`，清空缓冲区**并立即推到屏上**。lilt 每帧开头正是用背景色（默认主题是黑色）填满整屏，所以经 `FromDisplayer` 包装的 ssd1306 每帧会先推一次空白画面，可能看到闪烁，总线流量也翻倍。下一节的适配器只写缓冲区，可以绕开它。

## 反色光标：自己写适配器 { #inverter }

```go
type Inverter interface {
	InvertRect(x, y, w, h int16)
}
```

`gfx.Inverter` 是可选的。屏幕实现了它，光标就把焦点节点所在的区域反色：亮的变暗、暗的变亮，像 iPod 的菜单那样，文字和图标都留在原处。没实现时，光标画成主题 `Accent` 色的圆角轮廓。光标在所有页面和弹窗之后、Toast 之前绘制，所以反色作用在已经画好的内容上，Toast 不会被反色。

驱动自带帧缓冲时，反色就是在缓冲区里翻转像素。仓库里的 ESP32-S3 固件就是这样包装 SH1106 驱动的，它的 `Color` 是一位的：

```go
// oled adapts the driver to gfx.Screen; its own Color type is one bit.
type oled struct{ d *sh1106.Device }

func (o oled) Size() (w, h int16) { return o.d.Size() }

func (o oled) SetPixel(x, y int16, c gfx.Color) { o.d.Pixel(x, y, sh1106.Color(gfx.IsLit(c))) }

func (o oled) FillRect(x, y, w, h int16, c gfx.Color) {
	o.d.FillRect(x, y, w, h, sh1106.Color(gfx.IsLit(c)))
}

func (o oled) InvertRect(x, y, w, h int16) { o.d.InvertRect(x, y, w, h) }

func (o oled) Flush() error { return o.d.Display() }
```

tinygo 的 ssd1306 驱动没有 `InvertRect`，但有 `GetPixel`，可以逐点翻转。嵌入 `*ssd1306.Device` 直接得到 `Size` 和 `SetPixel`，只补另外三个方法：

```go
// oled adds inversion to the ssd1306 driver, which keeps a frame buffer.
type oled struct{ *ssd1306.Device }

// FillRect writes the buffer only. The driver's FillRectangle sends a blank
// frame to the panel whenever it is asked to clear the whole screen, which
// lilt does at the start of every frame.
func (o oled) FillRect(x, y, w, h int16, c gfx.Color) {
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			o.SetPixel(i, j, c)
		}
	}
}

func (o oled) InvertRect(x, y, w, h int16) {
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			c := gfx.White
			if o.GetPixel(i, j) {
				c = gfx.Black
			}
			o.SetPixel(i, j, c)
		}
	}
}

func (o oled) Flush() error { return o.Display() }
```

```go
dev := ssd1306.NewI2C(machine.I2C0)
dev.Configure(ssd1306.Config{Width: 128, Height: 64})
lilt.Start(oled{dev}, home)
```

光标每帧只反色一个圆角矩形，逐点翻转的开销和画一行字相当。

## 彩色屏与帧缓冲 { #color }

直接写屏的彩色驱动每画一个矩形就走一次总线，整屏重画时能看到画面一层层出来。在内存里放一块 RGB565 缓冲，画完一帧再一次推出去，就没有这个问题，顺带还能反色：

```go
// panel is what st7789, st7735 and ili9341 have in common.
type panel interface {
	Size() (w, h int16)
	DrawRGBBitmap8(x, y int16, data []uint8, w, h int16) error
}

// rgb565 draws into memory and sends the finished frame in one transfer.
type rgb565 struct {
	p    panel
	w, h int16
	buf  []byte // two bytes per pixel, high byte first
}

func newRGB565(p panel) *rgb565 {
	w, h := p.Size()
	return &rgb565{p: p, w: w, h: h, buf: make([]byte, int(w)*int(h)*2)}
}

func (s *rgb565) Size() (w, h int16) { return s.w, s.h }

func (s *rgb565) SetPixel(x, y int16, c gfx.Color) { s.FillRect(x, y, 1, 1, c) }

func (s *rgb565) FillRect(x, y, w, h int16, c gfx.Color) {
	v := uint16(c.R>>3)<<11 | uint16(c.G>>2)<<5 | uint16(c.B>>3)
	x0, y0, x1, y1 := s.clip(x, y, w, h)
	for j := y0; j < y1; j++ {
		for i := x0; i < x1; i++ {
			k := (j*int(s.w) + i) * 2
			s.buf[k], s.buf[k+1] = byte(v>>8), byte(v)
		}
	}
}

func (s *rgb565) InvertRect(x, y, w, h int16) {
	x0, y0, x1, y1 := s.clip(x, y, w, h)
	for j := y0; j < y1; j++ {
		for i := x0; i < x1; i++ {
			k := (j*int(s.w) + i) * 2
			s.buf[k] ^= 0xFF
			s.buf[k+1] ^= 0xFF
		}
	}
}

func (s *rgb565) Flush() error { return s.p.DrawRGBBitmap8(0, 0, s.buf, s.w, s.h) }

// clip keeps the rectangle on the screen, as the Screen contract asks.
func (s *rgb565) clip(x, y, w, h int16) (x0, y0, x1, y1 int) {
	return max(int(x), 0), max(int(y), 0), min(int(x)+int(w), int(s.w)), min(int(y)+int(h), int(s.h))
}
```

接到 st7789 上是 `lilt.Start(newRGB565(&dev), home)`。缓冲区的大小是宽 × 高 × 2 字节，240x240 要 115 KB，先确认芯片的 RAM 放得下。

彩色屏还要改主题。组件在构造时读取主题，所以在建树之前改：

```go
t := ui.CurrentTheme
t.Background = gfx.Color{R: 16, G: 18, B: 24, A: 255}
t.Foreground = gfx.White
t.Accent = gfx.Color{R: 80, G: 160, B: 255, A: 255} // 不能反色时光标的颜色
t.Font = &freesans.Regular18pt7b                   // 大屏换大字体

home := hello() // 改完主题再建树
```

换成行高更大的字体后，`ui.Unit()` 随之变大，所有以单位计的内边距、间距和圆角一起放大，见[主题与单位](../concepts/theme.md#unit)。默认字体是 12pt 的 Notosans；`freesans` 是 `tinygo.org/x/tinyfont/freesans`，任何 tinyfont 字体都可以。

## 按键 { #keys }

按键由板子上的代码自己读，再告诉框架。框架只认逻辑键：

| 键 | 默认行为 |
| --- | --- |
| `ui.KeyUp`、`KeyDown`、`KeyLeft`、`KeyRight` | 没有组件认领时移动焦点，滚动容器在边缘时滚动 |
| `ui.KeyOK` | 激活焦点组件；按下时光标收缩 |
| `ui.KeyBack` | 关闭当前弹窗或页面；在根页面上只轻轻顶一下 |
| `ui.KeyMenu` | 框架不处理，留给应用的 `KeyHandler` |
| `ui.KeyBackspace` | `Input` 和软键盘删除一个字符 |
| `ui.KeyChar` | 由 `ui.EmitRune` 产生，字符在 `KeyEvent.Rune` 里 |

每次按键要发 `Press` 和 `Release` 两个事件，按住时中间可以有若干 `Repeat`。动作在 `Release` 或 `Repeat` 时发生，`Press` 只给反馈，见[激活](../concepts/input.md#activation)。

### 用 `KeyPoller` 轮询引脚

最省事的办法是在一个循环里轮询所有按键，交给 `ui.KeyPoller`，它负责消抖和按住连发。下面是 ESP32-S3 固件的主循环，`lilt.Start` 之后 UI 在自己的 goroutine 上跑，`main` 接着读按键：

```go
poller := ui.NewKeyPoller()
period := time.Second / maxFPS
last := time.Now()
for {
	now := time.Now()
	dt := min(float32(now.Sub(last).Seconds()), 0.1)
	last = now
	for i, k := range keys {
		down := k.pin.Get()
		held[i].Store(down)
		poller.Poll(k.key, down, dt)
	}
	time.Sleep(period)
}
```

`Poll(k, down, dt)` 每个键每轮调用一次，`down` 是这一刻是否按着，`dt` 是距上一轮的秒数。这块板子的按键接 3.3V、引脚下拉，所以高电平是按下；按键接地、引脚上拉时要写 `!pin.Get()`，见[快速开始](../getting-started.md#接到板子上)的例子。

| 字段 | 默认值 | 作用 |
| --- | --- | --- |
| `Debounce` | 2 | 连续几次采样一致才算状态改变 |
| `RepeatDelay` | 0.2 | 按住多少秒后开始连发 |
| `RepeatInterval` | 0.07 | 第一次连发的间隔，秒 |
| `MinInterval` | 0.03 | 连发间隔的下限 |
| `Accel` | 0.9 | 每连发一次，间隔乘以它，越按越快 |

消抖按采样次数计，所以和轮询周期有关：周期 5ms、`Debounce` 为 2 时，电平要稳定约 10ms 才被接受。用轮询而不用中断，是因为轮询天然不受按键抖动产生的一串中断影响。

### 直接 `ui.Emit`

已经有自己的消抖逻辑，或者按键来自别的芯片（IO 扩展、红外遥控、串口），就直接发事件：

```go
ui.Emit(ui.KeyOK, ui.Press)
// ...
ui.Emit(ui.KeyOK, ui.Release)
```

`Emit` 可以在任何 goroutine 上调用。事件进一个容量 32 的队列，在下一帧开头统一分发；队列满了新事件直接丢弃，UI 卡住时不会反过来阻塞读按键的 goroutine。`lilt.Pause` 期间按键也在队列里等着。

### 触摸屏

触摸芯片报的是坐标，用 `ui.EmitPointer` 发出去，坐标是屏幕像素：

```go
// readTouch stands for the touch controller's driver.
var readTouch func() (x, y int16, touching bool)

func pollTouch() {
	var was bool
	var lx, ly int16
	for {
		x, y, touching := readTouch()
		switch {
		case touching && !was:
			ui.EmitPointer(x, y, ui.PointerDown)
		case touching:
			ui.EmitPointer(x, y, ui.PointerMove)
		case was:
			ui.EmitPointer(lx, ly, ui.PointerUp)
		}
		was, lx, ly = touching, x, y
		time.Sleep(10 * time.Millisecond)
	}
}
```

框架把触摸翻译成按键：按下时把焦点移到指下的节点（它本身不可聚焦时取最近的可聚焦祖先），并发出 `KeyOK` 的 `Press`；抬起时发 `Release`，所以按钮、列表行不用为触摸写任何代码。手指移动超过 4 像素就算拖动，滚动容器接管拖动，那次 `Release` 带上 `Canceled`，不会误触发。没有按下时的 `PointerMove`（鼠标悬停）只移动焦点。

### 物理键盘

有真正键盘的板子，字符用 `ui.EmitRune(r)` 发，退格用 `KeyBackspace`：

```go
ui.EmitRune('a')
ui.Emit(ui.KeyBackspace, ui.Press)
ui.Emit(ui.KeyBackspace, ui.Release)
```

`EmitRune` 只发一个 `KeyChar` 的 `Press`，不需要配对的 `Release`。焦点在 `Input` 上时字符直接打进去；软键盘打开时也接收这两种事件。退格在 `Press` 和每次 `Repeat` 时删除，所以按住会连删。如果不想在按 OK 时弹出软键盘，把 `ui.TextInput` 换成自己的输入法，见[文字输入](../widgets/input.md#text-input)。

## 帧率与电子墨水屏 { #fps }

`lilt.FPS` 是帧率上限，默认 60，每帧都会读，运行中也能改。帧循环画完一帧后睡到这一帧的时间用完；一帧本身就超时时不睡，只让出一次 CPU，让读按键的 goroutine 能跑。每帧的 `dt` 最大取 0.1 秒，所以卡顿之后动画不会一下跳到终点。

刷新慢的屏幕要注意：`Flush` 在 UI goroutine 上同步执行，它不返回，下一帧就不开始。

- **OLED 和带缓冲的 LCD**：瓶颈一般是总线。SPI 推一帧 128x64 单色画面只要 1KB，ESP32-S3 固件把 `lilt.FPS` 设到 120，按键也按这个频率采样。
- **直接写屏的彩色驱动**：每帧整屏重画，帧率受总线限制，还会闪，建议用[帧缓冲](#color)。
- **电子墨水屏**：一次 `Display` 要刷新整屏，耗时可达数秒。把 `lilt.FPS` 调到 1 到 2，并且让动画在一两帧内结束，例如把 `anim.TimeScale` 调大，否则一段弹簧动画会触发很多次整屏刷新。另外，任何还在动的东西都会让框架每帧重画：放不下的 `Label` 获得焦点时会滚动（跑马灯），回调运行期间组件会显示加载点，在墨水屏上都意味着连续刷新，设计界面时避开它们。
