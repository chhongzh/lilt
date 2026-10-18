// Command esp32s3-sh1106 runs the widget gallery on an ESP32-S3 with a
// 128x64 SH1106 panel over SPI and seven buttons wired to 3V3.
package main

import (
	"machine"
	"time"

	"osu-mc/driver/sh1106"
	"osu-mc/examples/gallery"
	"osu-mc/gfx"
	"osu-mc/ui"
)

// GPIO19/20 are USB, 26-32 flash and PSRAM, 43/44 UART0; these avoid them.
const (
	pinSCK = machine.GPIO4
	pinSDO = machine.GPIO5
	pinRES = machine.GPIO6
	pinDC  = machine.GPIO7
	pinCS  = machine.GPIO15
)

var keys = [...]struct {
	pin machine.Pin
	key ui.Key
}{
	{machine.GPIO16, ui.KeyUp},
	{machine.GPIO17, ui.KeyDown},
	{machine.GPIO18, ui.KeyLeft},
	{machine.GPIO8, ui.KeyRight},
	{machine.GPIO3, ui.KeyOK},
	{machine.GPIO46, ui.KeyBack},
	{machine.GPIO9, ui.KeyMenu},
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
	err := machine.SPI0.Configure(machine.SPIConfig{Frequency: 10_000_000, SCK: pinSCK, SDO: pinSDO})
	if err != nil {
		panic(err)
	}
	dev := sh1106.New(machine.SPI0, pinDC, pinRES, pinCS)
	if err := dev.Configure(sh1106.Config{}); err != nil {
		panic(err)
	}
	for _, k := range keys {
		k.pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	app := ui.NewApp(oled{dev})
	gallery.Build()
	poller := ui.NewKeyPoller()

	// Keys are sampled once per frame, so this loop replaces app.Run.
	period := time.Second / maxFPS
	last := time.Now()
	for {
		now := time.Now()
		dt := min(float32(now.Sub(last).Seconds()), 0.1)
		last = now
		for _, k := range keys {
			poller.Poll(k.key, k.pin.Get(), dt)
		}
		if _, err := app.Frame(dt); err != nil {
			panic(err)
		}
		if spent := time.Since(now); spent < period {
			time.Sleep(period - spent)
		}
	}
}
