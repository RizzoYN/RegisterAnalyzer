package main

import (
	"fmt"

	"github.com/pwiecz/go-fltk"
)

var (
	WIDTH    = 33 * 20
	HEIGHT   = 42
	colorMap = map[string]fltk.Color{
		"0":    fltk.WHITE,
		"1":    fltk.BACKGROUND_COLOR,
	}
	textMap = map[string]string{
		"0": "1",
		"1": "0",
	}
	padX, padY = 2, 2
	bitHeight, bitWidth = 18, 18
	dataWidth = 32
	maxRow = 3
	Row = 1
)

type Bit struct {
	Bit *fltk.TextDisplay
}

func (b *Bit) Click(e fltk.Event) bool {
	val := b.Bit.Buffer().Text()
	if e == fltk.Event(fltk.LeftMouse) {
		b.Bit.Buffer().SetText(textMap[val])
		return true
	}
	return false
}

func (b *Bit) DrawHandler(f func()) {
	val := b.Bit.Buffer().Text()
	fltk.SetDrawFont(fltk.HELVETICA, 14)
	fltk.DrawBox(fltk.BORDER_BOX, b.Bit.X(), b.Bit.Y(), b.Bit.W(), b.Bit.H(), colorMap[val])
	fltk.SetDrawColor(fltk.BLACK)
	fltk.Draw(textMap[textMap[val]], b.Bit.X(), b.Bit.Y(), b.Bit.W(), b.Bit.H(), fltk.ALIGN_CENTER)
}

func NewBit(x, y, w, h int) *Bit {
	b := new(Bit)
	bit := fltk.NewTextDisplay(x, y, w, h)
	buf := fltk.NewTextBuffer()
	buf.SetText("0")
	bit.SetBuffer(buf)
	bit.SetDrawHandler(b.DrawHandler)
	b.Bit = bit
	b.Bit.SetEventHandler(b.Click)
	return b
}

type Header struct {
	Header *fltk.TextDisplay
}

func (h *Header) DrawHandler(color fltk.Color) func(func()) {
	return func(fn func()) {
		fltk.SetDrawFont(fltk.HELVETICA, 14)
		fltk.SetDrawColor(color)
		fltk.Draw(h.Header.Buffer().Text(), h.Header.X(), h.Header.Y(), h.Header.W(), h.Header.H(), fltk.ALIGN_CENTER)
	}
}

func NewHeader(x, y, w, h int, ix string) *Header {
	header := new(Header)
	head := fltk.NewTextDisplay(x, y, w, h)
	buf := fltk.NewTextBuffer()
	buf.SetText(ix)
	head.SetBuffer(buf)
	head.SetTextColor(fltk.BLACK)
	head.SetDrawHandler(header.DrawHandler(fltk.BLACK))
	header.Header = head
	return header
}

func main() {
	fltk.InitStyles()
	win := fltk.NewWindow(WIDTH, HEIGHT)
	win.SetLabel("寄存器工具")
	win.SetColor(fltk.WHITE)
	for r := 0; r < 2; r++ {
		for i := 0; i < 32; i++ {
			if r == 0 {
				NewHeader(i*20+2, 2, 18, 18, fmt.Sprint(31-i))
			} else {
				NewBit(i*20+2, 20, 18, 18)
			}
		}
	}
	win.End()
	win.Show()
	fltk.Run()
}
