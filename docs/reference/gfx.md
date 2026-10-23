# gfx

[pkg.go.dev](https://pkg.go.dev/github.com/chhongzh/lilt/gfx) · `import "github.com/chhongzh/lilt/gfx"`

绘图层：驱动要实现的 `Screen` 接口，以及在它之上加了裁剪、平移和图形的 `Canvas`。

## 屏幕

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Screen interface { Size; SetPixel; FillRect; Flush }` | 驱动要提供的四个方法；超出屏幕的坐标要忽略 | [屏幕接口](../guides/hardware.md#screen) |
| `Size() (w, h int16)` | 屏幕像素尺寸 | [屏幕接口](../guides/hardware.md#screen) |
| `SetPixel(x, y int16, c Color)` | 画一个点 | [屏幕接口](../guides/hardware.md#screen) |
| `FillRect(x, y, w, h int16, c Color)` | 填充矩形 | [屏幕接口](../guides/hardware.md#screen) |
| `Flush() error` | 一帧画完后把画面推到屏上 | [屏幕接口](../guides/hardware.md#screen) |
| `type Inverter interface { InvertRect(x, y, w, h int16) }` | 可选：能反色的屏幕上，光标反色绘制 | [反色光标](../guides/hardware.md#inverter) |
| `type Displayer interface { Size; SetPixel; Display }` | 与 `tinygo.org/x/drivers` 的 `Displayer` 同签名，不 import 它 | [tinygo 驱动](../guides/hardware.md#screen) |
| `func FromDisplayer(d Displayer) Screen` | 包装任意 tinygo 驱动；有 `FillRectangle` 时用它填充，否则逐点画 | [接入硬件](../guides/hardware.md) |

## 颜色

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Color = color.RGBA` | 普通 RGBA，单色屏和彩色屏共用一套 API | [屏幕接口](../guides/hardware.md#screen) |
| `var Black`、`var White` | 不透明的黑和白 | [主题](../concepts/theme.md#theme) |
| `func IsLit(c Color) bool` | 单色屏上 `c` 是否点亮：按 3:6:1 加权的亮度不低于一半 | [屏幕接口](../guides/hardware.md#screen) |

## 画布

`Canvas` 是值类型，`Translate`、`Clip`、`Absolute` 返回新画布。坐标都是 `int16` 像素。

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Canvas struct` | 屏幕的一个平移、裁剪过的视图 | [画布](../guides/custom-widget.md#canvas) |
| `func NewCanvas(s Screen) Canvas` | 覆盖整块屏幕的画布；帧循环自己创建，组件在 `Draw` 里收到 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Translate(dx, dy int16) Canvas` | 原点移动 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Clip(r Rect) Canvas` | 只允许画在 `r`（本地坐标）以内 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Absolute() Canvas` | 原点回到屏幕左上角，保留裁剪 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Origin() (x, y int16)` | 本地原点的屏幕坐标 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Bounds() Rect` | 裁剪区，本地坐标 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Visible(r Rect) bool` | `r` 是否有一部分在裁剪区内 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) CanInvert() bool` | 屏幕是否实现了 `Inverter` | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Pixel(x, y int16, col Color)` | 一个点 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) HLine(x, y, w int16, col Color)` | 横线 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) VLine(x, y, h int16, col Color)` | 竖线 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) Line(x0, y0, x1, y1 int16, col Color)` | 任意方向的直线（Bresenham） | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) FillRect(r Rect, col Color)` | 实心矩形 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) StrokeRect(r Rect, col Color)` | 一像素边框 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) FillRoundRect(r Rect, radius int16, col Color)` | 实心圆角矩形，按行填充 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) StrokeRoundRect(r Rect, radius int16, col Color)` | 圆角边框 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) FillCircle(cx, cy, r int16, col Color)` | 实心圆 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) StrokeCircle(cx, cy, r int16, col Color)` | 圆周 | [画布](../guides/custom-widget.md#canvas) |
| `func (c Canvas) InvertRect(r Rect)` | 区域反色，屏幕不能反色时什么也不做 | [反色光标](../guides/hardware.md#inverter) |
| `func (c Canvas) InvertRoundRect(r Rect, radius int16)` | 圆角区域反色 | [反色光标](../guides/hardware.md#inverter) |
| `func (c Canvas) Text(f *Font, x, y int16, s string, col Color) int16` | 画一行字，`y` 是基线，返回末尾的 `x` | [画布](../guides/custom-widget.md#canvas) |

## 矩形

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Rect struct { X, Y, W, H int16 }` | 像素矩形；布局用的是 float32 的 `ui.Rect`，用 `Pixels()` 转过来 | [测量与绘制](../guides/custom-widget.md#measure) |
| `func (r Rect) Empty() bool` | 宽或高不大于 0 | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Right() int16`、`Bottom() int16` | 右边、下边（不含） | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Contains(x, y int16) bool` | 点是否在内 | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Offset(dx, dy int16) Rect` | 平移 | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Inset(d int16) Rect` | 四边各缩进 `d`，负数为扩大 | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Intersect(o Rect) Rect` | 交集，不相交时为空矩形 | [画布](../guides/custom-widget.md#canvas) |
| `func (r Rect) Overlaps(o Rect) bool` | 是否相交 | [画布](../guides/custom-widget.md#canvas) |

## 字体

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `type Font = tinyfont.Font` | tinyfont 位图字体，字形需按 rune 排序（tinyfont 生成器就是这样输出的） | [主题](../concepts/theme.md#theme) |
| `func Ascent(f *Font) int16` | 最高字形顶部到基线的距离 | [画布](../guides/custom-widget.md#canvas) |
| `func LineHeight(f *Font) int16` | 一行字占的高度；`ui.Unit()` 由它算出 | [单位](../concepts/theme.md#unit) |
| `func TextWidth(f *Font, s string) int16` | 文字的前进宽度 | [画布](../guides/custom-widget.md#canvas) |
