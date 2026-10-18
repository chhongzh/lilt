// Command esp32s3-generic runs the lilt gallery on an ESP32-S3 with a 128x64
// SH1106 over SPI, a five-way joystick and SET and RST keys. README.md has
// the parts list and wiring.
package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/chhongzh/lilt/driver/sh1106"
	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

const (
	SCREEN_CLK_PIN = machine.GPIO4
	SCREEN_SDA_PIN = machine.GPIO5
	SCREEN_RES_PIN = machine.GPIO6
	SCREEN_DC_PIN  = machine.GPIO7
	SCREEN_CS_PIN  = machine.GPIO15

	KBD_UP_PIN    = machine.GPIO16
	KBD_DOWN_PIN  = machine.GPIO17
	KBD_LEFT_PIN  = machine.GPIO18
	KBD_RIGHT_PIN = machine.GPIO8
	KBD_MID_PIN   = machine.GPIO3
	KBD_SET_PIN   = machine.GPIO46
	KBD_RST_PIN   = machine.GPIO9
)

// RST is a plain key on GPIO9, not the chip's reset line, so it gets Menu.
var keys = [padCount]struct {
	pin machine.Pin
	key ui.Key
}{
	padUp:    {KBD_UP_PIN, ui.KeyUp},
	padDown:  {KBD_DOWN_PIN, ui.KeyDown},
	padLeft:  {KBD_LEFT_PIN, ui.KeyLeft},
	padRight: {KBD_RIGHT_PIN, ui.KeyRight},
	padMid:   {KBD_MID_PIN, ui.KeyOK},
	padSet:   {KBD_SET_PIN, ui.KeyBack},
	padRst:   {KBD_RST_PIN, ui.KeyMenu},
}

var wiring = [...]struct {
	name string
	pin  machine.Pin
}{
	{"Screen CLK", SCREEN_CLK_PIN},
	{"Screen SDA", SCREEN_SDA_PIN},
	{"Screen RES", SCREEN_RES_PIN},
	{"Screen DC", SCREEN_DC_PIN},
	{"Screen CS", SCREEN_CS_PIN},
	{"Joystick up", KBD_UP_PIN},
	{"Joystick down", KBD_DOWN_PIN},
	{"Joystick left", KBD_LEFT_PIN},
	{"Joystick right", KBD_RIGHT_PIN},
	{"Joystick press", KBD_MID_PIN},
	{"SET", KBD_SET_PIN},
	{"RST", KBD_RST_PIN},
}

const maxFPS = 120

// oled adapts the driver to gfx.Screen; its own Color type is one bit.
type oled struct{ d *sh1106.Device }

func (o oled) Size() (w, h int16) { return o.d.Size() }

func (o oled) SetPixel(x, y int16, c gfx.Color) { o.d.Pixel(x, y, sh1106.Color(gfx.IsLit(c))) }

func (o oled) FillRect(x, y, w, h int16, c gfx.Color) {
	o.d.FillRect(x, y, w, h, sh1106.Color(gfx.IsLit(c)))
}

func (o oled) InvertRect(x, y, w, h int16) { o.d.InvertRect(x, y, w, h) }

func (o oled) Flush() error { return o.d.Display() }

func main() {
	err := machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 10_000_000,
		SCK:       SCREEN_CLK_PIN,
		SDO:       SCREEN_SDA_PIN,
	})
	if err != nil {
		panic(err)
	}
	dev := sh1106.New(machine.SPI0, SCREEN_DC_PIN, SCREEN_RES_PIN, SCREEN_CS_PIN)
	if err := dev.Configure(sh1106.Config{}); err != nil {
		panic(err)
	}

	var held [padCount]bool
	var gpio [padCount]string
	for i, k := range keys {
		k.pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
		gpio[i] = strconv.Itoa(int(k.pin))
	}

	app := ui.NewApp(oled{dev})
	gallery.Build(
		gallery.Demo{Title: "Keys", Build: func() ui.Drawable {
			return gallery.Page("Keys", newKeypad(&held, &gpio))
		}},
		gallery.Demo{Title: "Pins", Build: pins},
	)
	poller := ui.NewKeyPoller()

	// Keys are sampled once per frame, so this loop replaces app.Run.
	period := time.Second / maxFPS
	last := time.Now()
	for {
		now := time.Now()
		dt := min(float32(now.Sub(last).Seconds()), 0.1)
		last = now
		for i, k := range keys {
			held[i] = k.pin.Get()
			poller.Poll(k.key, held[i], dt)
		}
		if _, err := app.Frame(dt); err != nil {
			panic(err)
		}
		if spent := time.Since(now); spent < period {
			time.Sleep(period - spent)
		}
	}
}

// pins lists every connection with its GPIO, for checking the wiring.
func pins() ui.Drawable {
	small := widget.Font(ui.CurrentTheme.SmallFont)
	l := widget.NewList()
	for _, w := range wiring {
		gpio := widget.NewLabel("GPIO"+strconv.Itoa(int(w.pin)), small)
		ui.Add(l, widget.NewListItem(w.name, widget.Trailing(gpio)))
	}
	return gallery.Page("Pins", l)
}
