package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pwiecz/go-fltk"
)

var (
	bitColorMap = map[string]fltk.Color{
		"0": fltk.WHITE,
		"1": fltk.BACKGROUND_COLOR,
	}
	headerColorMap = map[string]fltk.Color{
		"same": fltk.BLACK,
		"diff": fltk.RED,
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
	Row       = 2
	WIDTH     = dataWidth*(bitW+pad) + pad*8 + bitW*13 + 50
	HEIGHT    = bitW + maxRow*bitH + pad*4
	MaxNum    = int64(math.Pow(2, float64(dataWidth)) - 1)
)

func ParseHeight(row int) int {
	if row == 1 {
		return row * (pad + bitW)
	} else {
		return row * bitH
	}
}

type Bit struct {
	*fltk.Box
}

func (b *Bit) Click(e fltk.Event) bool {
	if e == fltk.Event(fltk.LeftMouse) {
		val := b.Label()
		str := textMap[val]
		b.SetLabel(str)
		b.SetColor(bitColorMap[str])
		return true
	}
	return false
}

func NewBit(x, y, w, h int) *Bit {
	bit := fltk.NewBox(fltk.FLAT_BOX, x, y, w, h, "0")
	bit.SetAlign(fltk.ALIGN_CENTER)
	bit.SetColor(bitColorMap["0"])
	bit.SetLabelSize(14)
	bit.SetBox(fltk.BORDER_BOX)

	return &Bit{bit}
}

type Header struct {
	*fltk.Box
}

func NewHeader(x, y, w, h, ix int) *Header {
	header := fltk.NewBox(fltk.FLAT_BOX, x, y, w, h, fmt.Sprint(ix))
	header.SetAlign(fltk.ALIGN_CENTER)
	header.SetColor(fltk.WHITE)
	header.SetLabelColor(headerColorMap["same"])
	header.SetLabelSize(14)
	return &Header{header}
}

type Headers []*Header

func (h Headers) UpdateHeader(bitMap map[string]int, c int) {
	if len(bitMap) == 1 {
		h[c].SetColor(headerColorMap["same"])
	} else {
		h[c].SetColor(headerColorMap["diff"])
	}
}

func NewHeaders() Headers {
	headers := make([]*Header, dataWidth)
	for c := 0; c < dataWidth; c++ {
		head := NewHeader(c*(bitW+pad)+pad, pad, bitW, bitW, dataWidth-1-c)
		headers[c] = head
	}
	return headers
}

type BitRow struct {
	BitLocs  []*Bit
	Num      *fltk.TextEditor
	LShift   *fltk.Button
	ShiftNum *fltk.TextEditor
	RShift   *fltk.Button
	Reverse  *fltk.Button
	Invert   *fltk.Button
	Clear    *fltk.Button
	base     int
}

func (b *BitRow) GetBitString() []string {
	bitList := make([]string, dataWidth)
	for c := 0; c < dataWidth; c++ {
		bitList[c] = b.BitLocs[c].Label()
	}
	return bitList
}

func (b *BitRow) SetNum(bin int64) {
	switch b.base {
	case 16:
		b.Num.Buffer().SetText(fmt.Sprintf("%x", bin))
	case 10:
		b.Num.Buffer().SetText(fmt.Sprint(bin))
	case 8:
		b.Num.Buffer().SetText(fmt.Sprintf("%o", bin))
	}
}

func (b *BitRow) UpdateNum() {
	bitList := b.GetBitString()
	binStr := strings.Join(bitList, "")
	bin, _ := strconv.ParseInt(binStr, 2, dataWidth*2)
	b.SetNum(bin)
}

func (b *BitRow) ClickClear(fn func()) func() {
	return func() {
		for c := 0; c < dataWidth; c++ {
			b.BitLocs[c].SetLabel("0")
			b.BitLocs[c].SetColor(bitColorMap["0"])
		}
		b.Num.Buffer().SetText("0")
		fn()
	}
}

func (b *BitRow) ClickInvert(fn func()) func() {
	return func() {
		for c := 0; c < dataWidth; c++ {
			b.BitLocs[c].Click(fltk.Event(fltk.LeftMouse))
		}
		b.UpdateNum()
		fn()
	}
}

func (b *BitRow) GetCurrentNum() (int64, int) {
	num, _ := strconv.ParseInt(b.Num.Buffer().Text(), b.base, dataWidth*2)
	shiftNum, _ := strconv.ParseInt(b.ShiftNum.Buffer().Text(), 10, dataWidth*2)
	return num, int(shiftNum)
}

func (b *BitRow) UpdateBit(num int64) {
	binStr := strconv.FormatInt(num, 2)
	n := len(binStr)
	sum := 0
	for c := dataWidth - 1; c >= 0; c-- {
		if sum < n {
			s := string(binStr[n-sum-1])
			b.BitLocs[c].SetLabel(s)
			b.BitLocs[c].SetColor(bitColorMap[s])
		} else {
			b.BitLocs[c].SetLabel("0")
			b.BitLocs[c].SetColor(bitColorMap["0"])
		}
		sum++
	}
}

func (b *BitRow) UpdateBitNum(num int64) {
	b.SetNum(num)
	b.UpdateBit(num)
}

func (b *BitRow) ClickLShift(fn func()) func() {
	return func() {
		num, shiftNum := b.GetCurrentNum()
		num = (num << shiftNum) & MaxNum
		b.UpdateBitNum(num)
		fn()
	}
}

func (b *BitRow) ClickRShift(fn func()) func() {
	return func() {
		num, shiftNum := b.GetCurrentNum()
		num = (num >> shiftNum) & MaxNum
		b.UpdateBitNum(num)
		fn()
	}
}

func (b *BitRow) ClickReverse(fn func()) func() {
	return func() {
		for i, j := 0, len(b.BitLocs)-1; i < j; i, j = i+1, j-1 {
			h := b.BitLocs[i].Label()
			e := b.BitLocs[j].Label()
			b.BitLocs[i].SetLabel(e)
			b.BitLocs[i].SetColor(bitColorMap[e])
			b.BitLocs[j].SetLabel(h)
			b.BitLocs[j].SetColor(bitColorMap[h])
		}
		b.UpdateNum()
		fn()
	}
}

func (b *BitRow) KeyType(fn func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.KEYUP {
			num, _ := b.GetCurrentNum()
			b.UpdateBit(num)
			fn()
		}
		return false
	}
}

func (b *BitRow) Click(fn func(fltk.Event) bool, fnc func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			fn(e)
			fnc()
			b.UpdateNum()
			return true
		}
		return false
	}
}

func NewBitRow(row int, fn func()) *BitRow {
	bitRow := new(BitRow)
	bitLocs := make([]*Bit, dataWidth)
	h := ParseHeight(row)
	for c := 0; c < dataWidth; c++ {
		bit := NewBit(c*(bitW+pad)+pad, h, bitW, bitH)
		bit.SetEventHandler(bitRow.Click(bit.Click, fn))
		bitLocs[c] = bit
	}
	bitRow.BitLocs = bitLocs
	num := fltk.NewTextEditor(dataWidth*(bitW+pad)+pad, h, bitW*6, bitH)
	numBuf := fltk.NewTextBuffer()
	numBuf.SetText("0")
	num.SetBox(fltk.BORDER_BOX)
	num.SetBuffer(numBuf)
	num.SetInsertPosition(1)
	num.SetEventHandler(bitRow.KeyType(fn))
	bitRow.Num = num
	lShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*2+bitW*6, h, 25, bitH, "<<")
	lShift.SetBox(fltk.BORDER_BOX)
	lShift.ClearVisibleFocus()
	lShift.SetLabelSize(12)
	lShift.SetLabelFont(fltk.HELVETICA)
	lShift.SetDownBox(fltk.FLAT_BOX)
	lShift.SetCallback(bitRow.ClickLShift(fn))
	bitRow.LShift = lShift
	shiftNum := fltk.NewTextEditor(dataWidth*(bitW+pad)+pad*3+bitW*6+25, h, bitW, bitH)
	shiftBuf := fltk.NewTextBuffer()
	shiftBuf.SetText("1")
	shiftNum.SetBox(fltk.BORDER_BOX)
	shiftNum.SetBuffer(shiftBuf)
	shiftNum.SetInsertPosition(1)
	bitRow.ShiftNum = shiftNum
	rShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*4+bitW*7+25, h, 25, bitH, ">>")
	rShift.SetBox(fltk.BORDER_BOX)
	rShift.SetLabelSize(12)
	rShift.SetLabelFont(fltk.HELVETICA)
	rShift.SetDownBox(fltk.FLAT_BOX)
	rShift.ClearVisibleFocus()
	rShift.SetCallback(bitRow.ClickRShift(fn))
	bitRow.RShift = rShift
	reverse := fltk.NewButton(dataWidth*(bitW+pad)+pad*5+bitW*7+50, h, bitW*2, bitH, "倒序")
	reverse.SetBox(fltk.BORDER_BOX)
	reverse.SetLabelSize(12)
	reverse.SetLabelFont(fltk.HELVETICA)
	reverse.SetDownBox(fltk.FLAT_BOX)
	reverse.ClearVisibleFocus()
	reverse.SetCallback(bitRow.ClickReverse(fn))
	bitRow.Reverse = reverse
	invert := fltk.NewButton(dataWidth*(bitW+pad)+pad*6+bitW*9+50, h, bitW*2, bitH, "转换")
	invert.SetBox(fltk.BORDER_BOX)
	invert.SetLabelSize(12)
	invert.SetLabelFont(fltk.HELVETICA)
	invert.ClearVisibleFocus()
	invert.SetDownBox(fltk.FLAT_BOX)
	invert.SetCallback(bitRow.ClickInvert(fn))
	bitRow.Invert = lShift
	clear := fltk.NewButton(dataWidth*(bitW+pad)+pad*7+bitW*11+50, h, bitW*2, bitH, "清空")
	clear.SetBox(fltk.BORDER_BOX)
	clear.SetLabelSize(12)
	clear.SetLabelFont(fltk.HELVETICA)
	clear.ClearVisibleFocus()
	clear.SetDownBox(fltk.FLAT_BOX)
	clear.SetCallback(bitRow.ClickClear(fn))
	bitRow.Clear = clear
	bitRow.base = 16
	return bitRow
}

type MainForm struct {
	Headers Headers
	BitRows []*BitRow
}

func (m *MainForm) Updateheaders() {
	for c := 0; c < dataWidth; c++ {
		bitMap := make(map[string]int, Row)
		for r := 0; r < Row; r++ {
			val := m.BitRows[r].BitLocs[c].Label()
			bitMap[val] = 0
		}
		m.Headers.UpdateHeader(bitMap, c)
	}
}

func NewMainForm() {
	mainForm := new(MainForm)
	bitRows := make([]*BitRow, maxRow)
	for r := 0; r < maxRow+1; r++ {
		if r == 0 {
			mainForm.Headers = NewHeaders()
		} else {
			bitRows[r-1] = NewBitRow(r, mainForm.Updateheaders)
		}
	}
	mainForm.BitRows = bitRows
}

func main() {
	fltk.InitStyles()
	win := fltk.NewWindow(WIDTH, HEIGHT)
	win.SetLabel("寄存器工具")
	win.SetColor(fltk.WHITE)
	NewMainForm()
	win.End()
	win.Show()
	fltk.Run()
}
