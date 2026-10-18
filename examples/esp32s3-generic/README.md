# esp32s3-generic

在 ESP32-S3 上跑 lilt 的 gallery：一块 128x64 SH1106 OLED，一个五向摇杆（四个方向加按下），外加 SET、RST 两个键。照这份文档接好线就能复刻出同样的设备。

## 材料

| 部件 | 说明 |
| --- | --- |
| ESP32-S3 开发板 | 任意引出了下表 GPIO 的板子；TinyGo target 用 `esp32s3-generic` |
| 1.3 寸 SH1106 OLED | 128x64，**7 针 SPI 版**（GND VCC CLK SDA RES DC CS）。4 针的是 I²C 版，用不了 |
| 五向导航按键模块 | 引脚一般标作 COM UP DWN LFT RHT MID SET RST；摇杆向下按压就是 MID |
| 杜邦线 | 14 根左右 |

## 接线

| 模块引脚 | ESP32-S3 | 常量 |
| --- | --- | --- |
| OLED GND | GND | |
| OLED VCC | 3V3 | |
| OLED CLK | GPIO4 | `SCREEN_CLK_PIN` |
| OLED SDA | GPIO5 | `SCREEN_SDA_PIN` |
| OLED RES | GPIO6 | `SCREEN_RES_PIN` |
| OLED DC | GPIO7 | `SCREEN_DC_PIN` |
| OLED CS | GPIO15 | `SCREEN_CS_PIN` |
| 按键 COM | 3V3 | |
| 按键 UP | GPIO16 | `KBD_UP_PIN` |
| 按键 DWN | GPIO17 | `KBD_DOWN_PIN` |
| 按键 LFT | GPIO18 | `KBD_LEFT_PIN` |
| 按键 RHT | GPIO8 | `KBD_RIGHT_PIN` |
| 按键 MID | GPIO3 | `KBD_MID_PIN` |
| 按键 SET | GPIO46 | `KBD_SET_PIN` |
| 按键 RST | GPIO9 | `KBD_RST_PIN` |

要点：

- OLED 的 SDA 在 SPI 版上就是 MOSI。屏幕走 SPI0，时钟 10 MHz，线太长花屏时先把线缩短。
- 按键模块的 COM 必须接 **3V3**，不能接 GND。固件把按键引脚配成下拉输入，按下时读到高电平。
- 模块上的 RST 只是一个普通按键，接在 GPIO9 上，和开发板的复位脚无关。
- GPIO3 和 GPIO46 是 ESP32-S3 的 strapping 脚，上电时会被采样。按键平时悬空、下拉为低，不影响启动，但上电或复位的瞬间不要按住 MID 或 SET。
- 这些引脚避开了 USB（GPIO19/20）、Flash 与 PSRAM（GPIO26–32）和 UART0（GPIO43/44）。

## 按键映射

| 按键 | 框架里的键 | 作用 |
| --- | --- | --- |
| 上 / 下 / 左 / 右 | `KeyUp` / `KeyDown` / `KeyLeft` / `KeyRight` | 移动光标 |
| MID（摇杆按下） | `KeyOK` | 确认 |
| SET | `KeyBack` | 返回上一屏 |
| RST | `KeyMenu` | 菜单 |

按下时光标收缩，**松手**才触发，和 iOS 一样。按住方向键会加速连发。

## 演示内容

gallery 的首页在通用演示前多了两页：

- **Keys**：引脚指示。左边是摇杆的五个触点，右边是 SET 和 RST，每个触点里写着它的 GPIO 编号。读到高电平的引脚会亮起来，松开就熄灭。这里显示的是原始电平，没被页面用到的键也会亮，适合用来查接线。按 SET 会离开这一页，离开时能看到它亮了一下。
- **Pins**：列出全部 12 根信号线和对应的 GPIO，按上下键滚动。

## 构建与烧录

需要 TinyGo 0.43。在仓库根目录执行：

```sh
tinygo flash -target esp32s3-generic ./examples/esp32s3-generic
```

只编译、不烧录：

```sh
tinygo build -target esp32s3-generic -size short -o /tmp/gallery.bin ./examples/esp32s3-generic
```

当前大小约为 Flash 111 KB、RAM 31 KB。
