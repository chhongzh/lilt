// Command ledong-v2 runs the lilt gallery on the Labplus Ledong v2: an
// ESP32-S3 with a 1.47" 172x320 JD9853 LCD over SPI and two keys, A and B.
// The LCD wiring follows the official firmware (who_lcd.h in labplus-cn/esp-who
// and the board in xiaozhi-esp32), which differs from 1_3_IO.rst: CS is not
// driven, GPIO34 is the backlight and the STM8 on I2C resets the panel.
package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/chhongzh/lilt"
	"github.com/chhongzh/lilt/driver/jd9853"
	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

const (
	LCD_BL_PIN  = machine.GPIO34
	LCD_DC_PIN  = machine.GPIO35
	LCD_SCK_PIN = machine.GPIO36
	LCD_SDA_PIN = machine.GPIO37

	// The I2C bus to the STM8 coprocessor, which holds the LCD's RESET line.
	I2C_SCL_PIN = machine.GPIO43
	I2C_SDA_PIN = machine.GPIO44

	KEY_A_PIN = machine.GPIO0
	KEY_B_PIN = machine.GPIO46
)

var wiring = [...]struct {
	name string
	pin  machine.Pin
}{
	{"LCD BL", LCD_BL_PIN},
	{"LCD DC", LCD_DC_PIN},
	{"LCD SCK", LCD_SCK_PIN},
	{"LCD SDA", LCD_SDA_PIN},
	{"I2C SCL", I2C_SCL_PIN},
	{"I2C SDA", I2C_SDA_PIN},
	{"Key A", KEY_A_PIN},
	{"Key B", KEY_B_PIN},
}

// The STM8 resets the LCD when it receives this byte over I2C, as
// BOARD_STM8_ADDR and BOARD_STM8_CMD in who_lcd.h.
const (
	stm8Addr     = 17
	stm8LCDReset = 4
)

const (
	maxFPS = 30
	// A key held this long sends its long-press key instead.
	longPress = 400 * time.Millisecond
)

// Two keys cover the four the gallery needs: a tap and a hold on each.
var keys = [...]struct {
	pin         machine.Pin
	short, long ui.Key
}{
	{KEY_A_PIN, ui.KeyDown, ui.KeyUp},
	{KEY_B_PIN, ui.KeyOK, ui.KeyBack},
}

func main() {
	// Keep the backlight off until the frame memory is cleared.
	LCD_BL_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LCD_BL_PIN.Low()

	err := machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 20_000_000,
		SCK:       LCD_SCK_PIN,
		SDO:       LCD_SDA_PIN,
		SDI:       machine.NoPin,
	})
	if err != nil {
		panic(err)
	}
	err = machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400_000,
		SCL:       I2C_SCL_PIN,
		SDA:       I2C_SDA_PIN,
	})
	if err != nil {
		panic(err)
	}
	if err := machine.I2C0.Tx(stm8Addr, []byte{stm8LCDReset}, nil); err != nil {
		panic(err)
	}
	time.Sleep(120 * time.Millisecond)

	lcd := jd9853.New(machine.SPI0, LCD_DC_PIN)
	if err := lcd.Configure(jd9853.Config{}); err != nil {
		panic(err)
	}
	w, h := lcd.Size()
	lcd.FillRect(0, 0, w, h, gfx.Black)
	lcd.Flush()
	LCD_BL_PIN.High()
	lcd.On()

	for _, k := range keys {
		k.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	gallery.FitFonts(w, h)
	lilt.FPS = maxFPS
	lilt.Start(lcd, gallery.Home(gallery.Demo{Title: "Pins", Build: pins}))

	// The UI runs on its own goroutine; this one samples the keys. A key's
	// short action fires on release, its long action as soon as the hold is
	// long enough.
	var (
		down   [len(keys)]bool
		stable [len(keys)]uint8
		since  [len(keys)]time.Time
		fired  [len(keys)]bool
	)
	for {
		now := time.Now()
		for i, k := range keys {
			// The keys pull the pin low when pressed.
			raw := !k.pin.Get()
			if raw == down[i] {
				stable[i] = 0
			} else if stable[i]++; stable[i] >= 2 {
				down[i], stable[i] = raw, 0
				if raw {
					since[i], fired[i] = now, false
				} else if !fired[i] {
					tap(k.short)
				}
			}
			if down[i] && !fired[i] && now.Sub(since[i]) >= longPress {
				fired[i] = true
				tap(k.long)
			}
		}
		time.Sleep(time.Second / 120)
	}
}

func tap(k ui.Key) {
	ui.Emit(k, ui.Press)
	ui.Emit(k, ui.Release)
}

// pins lists every connection with its GPIO, for checking the wiring.
func pins() ui.Node {
	small := widget.Font(ui.CurrentTheme.SmallFont)
	l := widget.NewList()
	for _, w := range wiring {
		gpio := widget.NewLabel("GPIO"+strconv.Itoa(int(w.pin)), small)
		ui.Add(l, widget.NewListItem(w.name, widget.Trailing(gpio)))
	}
	return gallery.Page("Pins", l)
}
