package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pwiecz/go-fltk"
)

var (
	colorMap = map[string]fltk.Color{
		"0": fltk.WHITE,
		"1": fltk.BACKGROUND_COLOR,
	}
	textMap = map[string]string{
		"0": "1",
		"1": "0",
	}
	pad       = 2
	bitW      = 18
	bitH      = 22
	dataWidth = 32
	maxRow    = 2
	WIDTH     = dataWidth*(bitW+pad) + pad*8 + bitW*13 + 50
	HEIGHT    = bitW + maxRow*bitH + pad*4
)

type obj interface {
	Label() string
}

func ParseLoc(o obj) (int, int) {
	label := strings.Split(o.Label(), " ")
	r, _ := strconv.ParseInt(label[0], 10, 16)
	c, _ := strconv.ParseInt(label[1], 10, 16)
	return int(r), int(c)
}

func ParseHeight(row int) int {
	if row == 1 {
		return row * (pad + bitW)
	} else {
		return row * bitH
	}
}

type ClickHandler func(fltk.Event) bool

type Bit struct {
	Bit *fltk.TextDisplay
}

func (b *Bit) Click() ClickHandler {
	return func(e fltk.Event) bool {
		val := b.Bit.Buffer().Text()
		if e == fltk.Event(fltk.LeftMouse) {
			b.Bit.Buffer().SetText(textMap[val])
			return true
		}
		return false
	}
}

func (b *Bit) DrawHandler(f func()) {
	val := b.Bit.Buffer().Text()
	fltk.SetDrawFont(fltk.HELVETICA, 14)
	fltk.DrawBox(fltk.BORDER_BOX, b.Bit.X(), b.Bit.Y(), b.Bit.W(), b.Bit.H(), colorMap[val])
	fltk.SetDrawColor(fltk.BLACK)
	fltk.Draw(textMap[textMap[val]], b.Bit.X(), b.Bit.Y(), b.Bit.W(), b.Bit.H(), fltk.ALIGN_CENTER)
}

func NewBit(x, y, w, h int, ix string, fns []func(...interface{})) *Bit {
	b := new(Bit)
	bit := fltk.NewTextDisplay(x, y, w, h)
	buf := fltk.NewTextBuffer()
	buf.SetText("0")
	bit.SetBuffer(buf)
	bit.SetDrawHandler(b.DrawHandler)
	bit.SetLabel(ix)
	bit.SetLabelType(fltk.NO_LABEL)
	b.Bit = bit
	b.Bit.SetEventHandler(b.Click())
	return b
}

type Bits struct {
	Bits []*Bit
}

func NewBits(row int, fns []func(...interface{})) *Bits {
	bits := make([]*Bit, dataWidth)
	h := ParseHeight(row)
	for c := 0; c < dataWidth; c++ {
		bit := NewBit(c*(bitW+pad)+pad, h, bitW, bitH, fmt.Sprintf("%d %d", row-1, c), fns)
		bits[c] = bit
	}
	return &Bits{Bits: bits}
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

func (h *Header) UpdataHeader(len int) {
	if len == 1 {
		h.Header.SetDrawHandler(h.DrawHandler(fltk.BLACK))
	} else {
		h.Header.SetDrawHandler(h.DrawHandler(fltk.RED))
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

type Headers struct {
	Headers []*Header
}

func NewHeaders() *Headers {
	headers := make([]*Header, dataWidth)
	for c := 0; c < dataWidth; c++ {
		head := NewHeader(c*(bitW+pad)+pad, pad, bitW, bitW, fmt.Sprint(dataWidth-1-c))
		headers[c] = head
	}
	return &Headers{Headers: headers}
}

type NumShift struct {
	LShift   *fltk.Button
	ShiftNum *fltk.TextEditor
	RShift   *fltk.Button
}

func NewNumShift(row, base int, fns ...func(...interface{})) *NumShift {
	numShift := new(NumShift)
	h := ParseHeight(row)
	lShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*2+bitW*6, h, 25, bitH, "<<")
	lShift.SetBox(fltk.BORDER_BOX)
	lShift.SetDownBox(fltk.THIN_DOWN_BOX)
	lShift.SetType(uint8(row*(dataWidth+1)))
	// lShift.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+1))
	// lShift.SetLabelType(fltk.NO_LABEL)
	
	lShift.ClearVisibleFocus()
	numShift.LShift = lShift
	shiftNum := fltk.NewTextEditor(dataWidth*(bitW+pad)+pad*3+bitW*6+25, h, bitW, bitH)
	shiftBuf := fltk.NewTextBuffer()
	shiftBuf.SetText("1")
	shiftNum.SetBox(fltk.BORDER_BOX)
	shiftNum.SetBuffer(shiftBuf)
	shiftNum.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+2))
	shiftNum.SetLabelType(fltk.NO_LABEL)
	numShift.ShiftNum = shiftNum
	rShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*4+bitW*7+25, h, 25, bitH, ">>")
	rShift.SetBox(fltk.BORDER_BOX)
	rShift.SetDownBox(fltk.THIN_DOWN_BOX)
	// rShift.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+3))
	// rShift.SetLabelType(fltk.NO_LABEL)
	rShift.ClearVisibleFocus()
	numShift.RShift = rShift
	return numShift
}


type BitRow struct {
	Loc      *Bits
	Num      *fltk.TextEditor
	LShift   *fltk.TextDisplay
	ShiftNum *fltk.TextEditor
	RShift   *fltk.TextDisplay
	Reverse  *fltk.TextDisplay
	Invert   *fltk.TextDisplay
	Clear    *fltk.TextDisplay
}

func NewBitRow(row int, fns ...func(...interface{})) *BitRow {
	bitRow := new(BitRow)
	bitRow.Loc = NewBits(row, fns)
	h := ParseHeight(row)
	num := fltk.NewTextEditor(dataWidth*(bitW+pad)+pad, h, bitW*6, bitH)
	numBuf := fltk.NewTextBuffer()
	num.SetBox(fltk.BORDER_BOX)
	num.SetBuffer(numBuf)
	num.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth))
	num.SetLabelType(fltk.NO_LABEL)

	bitRow.Num = num
	lShift := fltk.NewTextDisplay(dataWidth*(bitW+pad)+pad*2+bitW*6, h, 25, bitH)
	lShift.SetBox(fltk.BORDER_BOX)
	lBuf := fltk.NewTextBuffer()
	lBuf.SetText("<<")
	lShift.SetBuffer(lBuf)
	lShift.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+1))
	lShift.SetLabelType(fltk.NO_LABEL)

	bitRow.LShift = lShift
	shiftNum := fltk.NewTextEditor(dataWidth*(bitW+pad)+pad*3+bitW*6+25, h, bitW, bitH)
	shiftBuf := fltk.NewTextBuffer()
	shiftBuf.SetText("1")
	shiftNum.SetBox(fltk.BORDER_BOX)
	shiftNum.SetBuffer(shiftBuf)
	shiftNum.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+2))
	shiftNum.SetLabelType(fltk.NO_LABEL)
	bitRow.ShiftNum = shiftNum
	rShift := fltk.NewTextDisplay(dataWidth*(bitW+pad)+pad*4+bitW*7+25, h, 25, bitH)
	rShift.SetBox(fltk.BORDER_BOX)
	rBuf := fltk.NewTextBuffer()
	rBuf.SetText(">>")
	rShift.SetBuffer(rBuf)
	rShift.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+3))
	rShift.SetLabelType(fltk.NO_LABEL)

	bitRow.RShift = rShift
	reverse := fltk.NewTextDisplay(dataWidth*(bitW+pad)+pad*5+bitW*7+50, h, bitW*2, bitH)
	reverse.SetBox(fltk.BORDER_BOX)
	reverseBuf := fltk.NewTextBuffer()
	reverseBuf.SetText("倒序")
	reverse.SetBuffer(reverseBuf)
	reverse.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+4))
	reverse.SetLabelType(fltk.NO_LABEL)

	bitRow.Reverse = reverse
	invert := fltk.NewTextDisplay(dataWidth*(bitW+pad)+pad*6+bitW*9+50, h, bitW*2, bitH)
	invert.SetBox(fltk.BORDER_BOX)
	invertBuf := fltk.NewTextBuffer()
	invertBuf.SetText("转换")
	invert.SetBuffer(invertBuf)
	invert.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+5))
	invert.SetLabelType(fltk.NO_LABEL)

	bitRow.Invert = lShift
	clear := fltk.NewTextDisplay(dataWidth*(bitW+pad)+pad*7+bitW*11+50, h, bitW*2, bitH)
	clear.SetBox(fltk.BORDER_BOX)
	clearBuf := fltk.NewTextBuffer()
	clearBuf.SetText("清空")
	clear.SetBuffer(clearBuf)
	clear.SetLabel(fmt.Sprintf("%d %d", row-1, dataWidth+6))
	clear.SetLabelType(fltk.NO_LABEL)

	bitRow.Clear = clear
	return bitRow
}



func main() {
	fltk.InitStyles()
	win := fltk.NewWindow(WIDTH, HEIGHT)
	win.SetLabel("寄存器工具")
	win.SetColor(fltk.WHITE)
	for r := 0; r < maxRow+1; r++ {
		if r == 0 {
			NewHeaders()
		} else {
			NewNumShift(r,r)
		}
	}
	win.End()
	win.Show()
	fltk.Run()
}
