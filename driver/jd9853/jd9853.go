// Package jd9853 drives the JD9853 colour LCD controller over SPI, as on the
// 1.47" 172x320 panels. It follows labplus's esp_lcd_panel_jd9853.c and
// who_lcd.c in github.com/labplus-cn/esp-who.
//
// The Device keeps the frame in RAM as RGB565 big-endian, the controller's
// wire format, and implements gfx.Screen and gfx.Inverter on it. Flush sends
// only the rows drawn since the last Flush, in one transfer: drawing through
// the controller directly would cost one SPI transaction per pixel.
//
// CS and RESET are left to the caller, since boards tie CS active or route
// RESET elsewhere. Reset the panel before Configure.
package jd9853

import (
	"machine"
	"time"

	"tinygo.org/x/drivers"

	"github.com/chhongzh/lilt/gfx"
)

const (
	cmdSWRESET = 0x01
	cmdSLPOUT  = 0x11
	cmdINVOFF  = 0x20
	cmdINVON   = 0x21
	cmdDISPOFF = 0x28
	cmdDISPON  = 0x29
	cmdCASET   = 0x2A
	cmdRASET   = 0x2B
	cmdRAMWR   = 0x2C
	cmdMADCTL  = 0x36
	cmdCOLMOD  = 0x3A
	cmdRAMCTRL = 0xB0

	colmodRGB565 = 0x55
)

// MADCTL bits, for Config.Madctl.
const (
	MY  = 0x80 // mirror rows
	MX  = 0x40 // mirror columns
	MV  = 0x20 // swap rows and columns
	BGR = 0x08 // blue first
)

// Landscape is the orientation of labplus's board config: BGR, rows and
// columns swapped, both mirrored.
const Landscape = BGR | MV | MX | MY

// Config describes the panel. The zero value is a 172x320 panel in
// landscape, as on the Labplus Ledong v2.
type Config struct {
	// Width and Height are the screen after Madctl, 320x172 by default.
	Width, Height int16
	// ColOffset and RowOffset place the screen in the controller's RAM,
	// after Madctl. The 172 pixel side of the default panel starts 34
	// lines in, so RowOffset defaults to 34 when Width and Height are left
	// at their defaults.
	ColOffset, RowOffset int16
	// Madctl sets the orientation and colour order, Landscape by default.
	Madctl uint8
	// Invert turns on display inversion, for panels that need it to show
	// true colours.
	Invert bool
}

// Device is a JD9853 with a frame buffer.
type Device struct {
	bus drivers.SPI
	dc  machine.Pin

	w, h             int16
	colOff, rowOff   int16
	buf              []byte
	dirtyLo, dirtyHi int16 // rows [dirtyLo, dirtyHi) changed since Flush
	cmd              [1]byte
	arg              [4]byte
}

// New creates a JD9853 connection. The SPI bus must already be configured.
func New(bus drivers.SPI, dc machine.Pin) *Device {
	dc.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dc.High()
	return &Device{bus: bus, dc: dc}
}

// Configure sends the initialisation sequence and allocates the frame
// buffer. The display stays off until On, so clear and Flush first or the
// panel shows whatever its RAM held.
func (d *Device) Configure(cfg Config) error {
	if cfg.Width == 0 && cfg.Height == 0 {
		cfg.Width, cfg.Height = 320, 172
		if cfg.RowOffset == 0 && cfg.ColOffset == 0 {
			cfg.RowOffset = 34
		}
	}
	if cfg.Madctl == 0 {
		cfg.Madctl = Landscape
	}
	d.w, d.h = cfg.Width, cfg.Height
	d.colOff, d.rowOff = cfg.ColOffset, cfg.RowOffset
	d.buf = make([]byte, int(d.w)*int(d.h)*2)
	d.dirtyLo, d.dirtyHi = 0, d.h

	// There is no RESET pin here, so a software reset stands in for
	// esp_lcd_panel_reset.
	if err := d.command(cmdSWRESET, nil); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)
	d.command(cmdSLPOUT, nil)
	time.Sleep(120 * time.Millisecond)

	d.command1(cmdMADCTL, BGR)
	d.command1(cmdCOLMOD, colmodRGB565)
	d.arg[0], d.arg[1] = 0x00, 0xF0 // big-endian RGB565
	d.command(cmdRAMCTRL, d.arg[:2])
	inv := byte(cmdINVOFF)
	if cfg.Invert {
		inv = cmdINVON
	}
	d.command(inv, nil)
	return d.command1(cmdMADCTL, cfg.Madctl)
}

// On turns the display on.
func (d *Device) On() error { return d.command(cmdDISPON, nil) }

// Off turns the display off; the frame memory keeps its contents.
func (d *Device) Off() error { return d.command(cmdDISPOFF, nil) }

func (d *Device) Size() (w, h int16) { return d.w, d.h }

func (d *Device) SetPixel(x, y int16, c gfx.Color) {
	if x < 0 || y < 0 || x >= d.w || y >= d.h {
		return
	}
	hi, lo := rgb565(c)
	i := (int(y)*int(d.w) + int(x)) * 2
	d.buf[i], d.buf[i+1] = hi, lo
	d.dirty(y, y+1)
}

func (d *Device) FillRect(x, y, w, h int16, c gfx.Color) {
	x0, y0, x1, y1, ok := d.clip(x, y, w, h)
	if !ok {
		return
	}
	hi, lo := rgb565(c)
	stride := int(d.w) * 2
	first := d.buf[int(y0)*stride+int(x0)*2 : int(y0)*stride+int(x1)*2]
	for i := 0; i < len(first); i += 2 {
		first[i], first[i+1] = hi, lo
	}
	for y := int(y0) + 1; y < int(y1); y++ {
		copy(d.buf[y*stride+int(x0)*2:], first)
	}
	d.dirty(y0, y1)
}

// InvertRect inverts every colour channel in the rectangle.
func (d *Device) InvertRect(x, y, w, h int16) {
	x0, y0, x1, y1, ok := d.clip(x, y, w, h)
	if !ok {
		return
	}
	stride := int(d.w) * 2
	for y := int(y0); y < int(y1); y++ {
		row := d.buf[y*stride+int(x0)*2 : y*stride+int(x1)*2]
		for i := range row {
			row[i] = ^row[i]
		}
	}
	d.dirty(y0, y1)
}

// Flush sends the rows changed since the last Flush.
func (d *Device) Flush() error {
	if d.dirtyLo >= d.dirtyHi {
		return nil
	}
	lo, hi := d.dirtyLo, d.dirtyHi
	d.dirtyLo, d.dirtyHi = d.h, 0
	d.setRange(cmdCASET, d.colOff, d.w)
	d.setRange(cmdRASET, d.rowOff+lo, hi-lo)
	stride := int(d.w) * 2
	return d.command(cmdRAMWR, d.buf[int(lo)*stride:int(hi)*stride])
}

func (d *Device) dirty(lo, hi int16) {
	d.dirtyLo, d.dirtyHi = min(d.dirtyLo, lo), max(d.dirtyHi, hi)
}

func (d *Device) clip(x, y, w, h int16) (x0, y0, x1, y1 int16, ok bool) {
	x0, y0 = max(x, 0), max(y, 0)
	x1, y1 = min(x+w, d.w), min(y+h, d.h)
	return x0, y0, x1, y1, x0 < x1 && y0 < y1
}

// command sends one command byte with DC low, then its parameters.
func (d *Device) command(c byte, data []byte) error {
	d.dc.Low()
	d.cmd[0] = c
	err := d.bus.Tx(d.cmd[:], nil)
	d.dc.High()
	if err == nil && len(data) > 0 {
		err = d.bus.Tx(data, nil)
	}
	return err
}

func (d *Device) command1(c, v byte) error {
	d.arg[0] = v
	return d.command(c, d.arg[:1])
}

func (d *Device) setRange(cmd byte, start, n int16) {
	end := start + n - 1
	d.arg = [4]byte{byte(start >> 8), byte(start), byte(end >> 8), byte(end)}
	d.command(cmd, d.arg[:])
}

func rgb565(c gfx.Color) (hi, lo byte) {
	v := uint16(c.R>>3)<<11 | uint16(c.G>>2)<<5 | uint16(c.B>>3)
	return byte(v >> 8), byte(v)
}
