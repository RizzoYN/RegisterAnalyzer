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
	MLmap = map[string]string{
		"MSB": "LSB",
		"LSB": "MSB",
	}
	pad       = 2
	bitW      = 18
	bitH      = 22
	dataWidth = 32
	maxRow    = 2
	Row       = 1
	WIDTH     = dataWidth*(bitW+pad) + pad*8 + bitW*13 + 50
	HEIGHT    = bitW + maxRow*bitH + pad*(3+maxRow) + 30
	MaxNum    = int64(math.Pow(2, float64(dataWidth)) - 1)
)

func ParseHeight(row int) int {
	if row == 1 {
		return pad*2 + 28 + bitW
	} else {
		return (row-1)*bitH + pad*(row+1) + 28 + bitW
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
	header.SetLabelSize(11)
	header.SetLabelFont(fltk.HELVETICA)
	return &Header{header}
}

type Headers []*Header

func (h Headers) UpdateHeader(bitMap map[string]int, c int) {
	if len(bitMap) == 1 {
		h[c].SetLabelColor(headerColorMap["same"])
		h[c].SetLabelSize(11)
		h[c].SetLabelFont(fltk.HELVETICA)
	} else {
		h[c].SetLabelColor(headerColorMap["diff"])
		h[c].SetLabelSize(14)
		h[c].SetLabelFont(fltk.HELVETICA_BOLD)
	}
	h[c].Redraw()
}

func NewHeaders() Headers {
	headers := make([]*Header, dataWidth)
	for c := 0; c < dataWidth; c++ {
		head := NewHeader(c*(bitW+pad)+pad, pad*2+28, bitW, bitW, dataWidth-1-c)
		headers[c] = head
	}
	return headers
}

type BitRow struct {
	BitLocs         []*Bit
	Num             *fltk.Input
	LShift          *fltk.Button
	ShiftNum        *fltk.Input
	RShift          *fltk.Button
	Reverse         *fltk.Button
	Invert          *fltk.Button
	Clear           *fltk.Button
	base            int
	lastNum         int64
	lastShiftNum    int64
	ShiftNumDisplay *fltk.Box
}

func (b *BitRow) GetBitString() []string {
	bitList := make([]string, dataWidth)
	for c := 0; c < dataWidth; c++ {
		bitList[c] = b.BitLocs[c].Label()
	}
	return bitList
}

func (b *BitRow) SetNum(num int64) {
	switch b.base {
	case 16:
		b.Num.SetValue(fmt.Sprintf("%x", num))
	case 10:
		b.Num.SetValue(fmt.Sprint(num))
	case 8:
		b.Num.SetValue(fmt.Sprintf("%o", num))
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
		b.Num.SetValue("0")
		if fn != nil {
			fn()
		}
		b.Display()
	}
}

func (b *BitRow) ClickInvert(fn func()) func() {
	return func() {
		for c := 0; c < dataWidth; c++ {
			b.BitLocs[c].Click(fltk.Event(fltk.LeftMouse))
		}
		b.UpdateNum()
		fn()
		b.Display()
	}
}

func (b *BitRow) GetCurrentNum() (int64, int64) {
	curNum := b.Num.Value()
	if curNum == "" {
		b.lastNum = 0
	}
	num, err := strconv.ParseInt(b.Num.Value(), b.base, dataWidth*2)
	if err == nil {
		b.lastNum = num
	} else {
		b.SetNum(b.lastNum)
	}
	shiftNum, err := strconv.ParseInt(b.ShiftNum.Value(), 10, dataWidth*2)
	if err == nil {
		b.lastShiftNum = shiftNum
	} else {
		if b.lastShiftNum != 0 {
			b.ShiftNum.SetValue(fmt.Sprint(b.lastShiftNum))
		} else {
			b.ShiftNum.SetValue("")
		}
	}
	return b.lastNum, b.lastShiftNum
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
		b.Display()
		fn()
	}
}

func (b *BitRow) ClickRShift(fn func()) func() {
	return func() {
		num, shiftNum := b.GetCurrentNum()
		num = (num >> shiftNum) & MaxNum
		b.UpdateBitNum(num)
		fn()
		b.Display()
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
		b.Display()
	}
}

func (b *BitRow) KeyTyped(fn func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			b.Display()
		}
		if e == fltk.KEYUP {
			num, _ := b.GetCurrentNum()
			b.UpdateBit(num)
			fn()
			b.Display()
			return true
		}
		return false
	}
}

func (b *BitRow) ShiftNumEvent(e fltk.Event) bool {
	if e == fltk.Event(fltk.LeftMouse) && b.ShiftNum.HasFocus() {
		b.ShiftNumDisplay.Show()
		b.ShiftNum.SetValue(fmt.Sprint(b.lastShiftNum))
		b.ShiftNum.Hide()
	}
	if e == fltk.KEYUP {
		_, shiftNum := b.GetCurrentNum()
		b.ShiftNumDisplay.SetLabel(fmt.Sprint(shiftNum))
		return true
	}
	return false
}

func (b *BitRow) Click(fn func(fltk.Event) bool, fnc func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			fn(e)
			fnc()
			b.UpdateNum()
			b.Display()
			return true
		}
		return false
	}
}

func (b *BitRow) DisplayClick(e fltk.Event) bool {
	if e == fltk.Event(fltk.LeftMouse) {
		b.ShiftNumDisplay.Hide()
		b.ShiftNum.SetValue("")
		b.ShiftNum.Show()
		b.ShiftNum.TakeFocus()
		return true
	}
	return false
}

func (b *BitRow) Display() {
	b.ShiftNum.Hide()
	b.ShiftNumDisplay.Show()
}

func (b *BitRow) Hide() {
	for _, obj := range b.BitLocs {
		obj.Hide()
	}
	b.Num.Hide()
	b.LShift.Hide()
	b.ShiftNum.Hide()
	b.RShift.Hide()
	b.Reverse.Hide()
	b.Invert.Hide()
	b.Clear.Hide()
	b.ShiftNumDisplay.Hide()
}

func (b *BitRow) Show() {
	for _, obj := range b.BitLocs {
		obj.Show()
	}
	b.Num.Show()
	b.LShift.Show()
	b.ShiftNum.Show()
	b.RShift.Show()
	b.Reverse.Show()
	b.Invert.Show()
	b.Clear.Show()
	b.ShiftNumDisplay.Show()
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
	num := fltk.NewInput(dataWidth*(bitW+pad)+pad, h, bitW*6, bitH)
	num.SetValue("0")
	num.SetBox(fltk.BORDER_BOX)
	num.SetEventHandler(bitRow.KeyTyped(fn))
	bitRow.Num = num
	lShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*2+bitW*6, h, 25, bitH, "<<")
	lShift.SetBox(fltk.GLEAM_UP_BOX)
	lShift.ClearVisibleFocus()
	lShift.SetLabelSize(12)
	lShift.SetLabelFont(fltk.HELVETICA)
	lShift.SetDownBox(fltk.GLEAM_DOWN_BOX)
	lShift.SetCallback(bitRow.ClickLShift(fn))
	bitRow.LShift = lShift
	shiftNum := fltk.NewInput(dataWidth*(bitW+pad)+pad*3+bitW*6+25, h, bitW, bitH)
	shiftNum.SetValue("1")
	shiftNum.SetBox(fltk.BORDER_BOX)
	shiftNum.SetEventHandler(bitRow.ShiftNumEvent)
	shiftNum.Hide()
	bitRow.ShiftNum = shiftNum
	rShift := fltk.NewButton(dataWidth*(bitW+pad)+pad*4+bitW*7+25, h, 25, bitH, ">>")
	rShift.SetBox(fltk.GLEAM_UP_BOX)
	rShift.SetLabelSize(12)
	rShift.SetLabelFont(fltk.HELVETICA)
	rShift.SetDownBox(fltk.GLEAM_DOWN_BOX)
	rShift.ClearVisibleFocus()
	rShift.SetCallback(bitRow.ClickRShift(fn))
	bitRow.RShift = rShift
	reverse := fltk.NewButton(dataWidth*(bitW+pad)+pad*5+bitW*7+50, h, bitW*2, bitH, "倒序")
	reverse.SetBox(fltk.GLEAM_UP_BOX)
	reverse.SetLabelSize(12)
	reverse.SetLabelFont(fltk.HELVETICA)
	reverse.SetDownBox(fltk.GLEAM_DOWN_BOX)
	reverse.ClearVisibleFocus()
	reverse.SetCallback(bitRow.ClickReverse(fn))
	bitRow.Reverse = reverse
	invert := fltk.NewButton(dataWidth*(bitW+pad)+pad*6+bitW*9+50, h, bitW*2, bitH, "转换")
	invert.SetBox(fltk.GLEAM_UP_BOX)
	invert.SetLabelSize(12)
	invert.SetLabelFont(fltk.HELVETICA)
	invert.ClearVisibleFocus()
	invert.SetDownBox(fltk.GLEAM_DOWN_BOX)
	invert.SetCallback(bitRow.ClickInvert(fn))
	bitRow.Invert = invert
	clear := fltk.NewButton(dataWidth*(bitW+pad)+pad*7+bitW*11+50, h, bitW*2, bitH, "清空")
	clear.SetBox(fltk.GLEAM_UP_BOX)
	clear.SetLabelSize(12)
	clear.SetLabelFont(fltk.HELVETICA)
	clear.ClearVisibleFocus()
	clear.SetDownBox(fltk.GLEAM_DOWN_BOX)
	clear.SetCallback(bitRow.ClickClear(fn))
	bitRow.Clear = clear
	bitRow.base = 16
	bitRow.lastNum = 0
	bitRow.lastShiftNum = 1
	shiftDisplay := fltk.NewBox(fltk.BORDER_BOX, dataWidth*(bitW+pad)+pad*3+bitW*6+25, h, bitW, bitH, fmt.Sprint(bitRow.lastShiftNum))
	shiftDisplay.SetAlign(fltk.ALIGN_CENTER)
	shiftDisplay.SetColor(fltk.WHITE)
	shiftDisplay.SetLabelFont(fltk.HELVETICA)
	shiftDisplay.SetEventHandler(bitRow.DisplayClick)
	bitRow.ShiftNumDisplay = shiftDisplay
	return bitRow
}

type MainForm struct {
	Headers        Headers
	BitRows        []*BitRow
	Compare        *fltk.ToggleButton
	Base16         *fltk.RadioRoundButton
	Base10         *fltk.RadioRoundButton
	Base8          *fltk.RadioRoundButton
	Bar            *fltk.Box
	MLSwitchButton *fltk.ToggleButton
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

func (m *MainForm) Add() func() {
	return func() {
		Row++
		for r := 0; r < Row-1; r++ {
			m.BitRows[r].ShiftNum.Hide()
			m.BitRows[r].ShiftNumDisplay.Show()
		}
		m.Bar.Hide()
		bitRow := m.BitRows[Row-1]
		bitRow.Show()
		m.Updateheaders()
	}
}

func (m *MainForm) Remove() func() {
	return func() {
		Row--
		for r := 0; r < Row; r++ {
			m.BitRows[r].ShiftNum.Hide()
			m.BitRows[r].ShiftNumDisplay.Show()
		}
		bitRow := m.BitRows[Row]
		bitRow.ClickClear(nil)()
		bitRow.Hide()
		m.Bar.Show()
		m.Updateheaders()
	}
}

func (m *MainForm) BaseChoise(base int) func() {
	return func() {
		for r := 0; r < maxRow; r++ {
			num, _ := m.BitRows[r].GetCurrentNum()
			m.BitRows[r].base = base
			m.BitRows[r].SetNum(num)
			if r < Row {
				m.BitRows[r].ShiftNum.Hide()
				m.BitRows[r].ShiftNumDisplay.Show()
			}
		}
	}
}

func (m *MainForm) MLSwitch() {
	t := m.MLSwitchButton.Label()
	for c := 0; c < dataWidth; c++ {
		var label string
		if t == "LSB" {
			label = fmt.Sprint(dataWidth - 1 - c)
		} else {
			label = fmt.Sprint(c)
		}
		m.Headers[c].SetLabel(label)
	}
	m.MLSwitchButton.SetLabel(MLmap[t])
}

func (m *MainForm) CompareReg() {
	if Row == 1 {
		m.Add()()
	} else {
		m.Remove()()
	}
}

func NewMainForm() {
	mainForm := new(MainForm)
	bitRows := make([]*BitRow, maxRow)
	for r := 0; r <= maxRow; r++ {
		if r == 0 {
			mainForm.Headers = NewHeaders()
		} else {
			bitRow := NewBitRow(r, mainForm.Updateheaders)
			if r > Row {
				bitRow.Hide()
			}
			bitRows[r-1] = bitRow
		}
	}
	mainForm.BitRows = bitRows
	box := fltk.NewBox(fltk.BORDER_BOX, pad*5+30, pad*4, 118, 20, "进制")
	box.SetLabelSize(12)
	box.SetColor(fltk.WHITE)
	box.SetAlign(fltk.ALIGN_LEFT)
	base16 := fltk.NewRadioRoundButton(pad*8+30, pad*5, 16, 16, "16")
	base16.ClearVisibleFocus()
	base16.SetValue(true)
	base16.SetCallback(mainForm.BaseChoise(16))
	mainForm.Base16 = base16
	base10 := fltk.NewRadioRoundButton(pad*8+70, pad*5, 16, 16, "10")
	base10.ClearVisibleFocus()
	base10.SetCallback(mainForm.BaseChoise(10))
	mainForm.Base10 = base10
	base8 := fltk.NewRadioRoundButton(pad*8+110, pad*5, 16, 16, "8")
	base8.ClearVisibleFocus()
	base8.SetCallback(mainForm.BaseChoise(8))
	mainForm.Base8 = base8
	compare := fltk.NewToggleButton(pad*5+150, pad*4, 75, 20, "寄存器对比")
	compare.SetBox(fltk.GLEAM_UP_BOX)
	compare.SetLabelSize(12)
	compare.SetLabelFont(fltk.HELVETICA)
	compare.ClearVisibleFocus()
	compare.SetDownBox(fltk.GLEAM_DOWN_BOX)
	compare.SetCallback(mainForm.CompareReg)
	bar := fltk.NewBox(fltk.FLAT_BOX, pad, HEIGHT-bitH-pad*3, WIDTH-pad*2, bitH, "增加对比")
	bar.SetColor(fltk.LIGHT1)
	bar.SetLabelSize(12)
	bar.SetLabelFont(fltk.HELVETICA)
	bar.ClearVisibleFocus()
	bar.SetEventHandler(func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			compare.SetValue(true)
			mainForm.CompareReg()
			return true
		}
		return false
	})
	mainForm.Compare = compare
	mainForm.Bar = bar
	mlSwitch := fltk.NewToggleButton(pad*6+225, pad*4, 35, 20, "MSB")
	mlSwitch.SetBox(fltk.GLEAM_UP_BOX)
	mlSwitch.SetLabelSize(12)
	mlSwitch.SetLabelFont(fltk.HELVETICA)
	mlSwitch.ClearVisibleFocus()
	mlSwitch.SetDownBox(fltk.GLEAM_DOWN_BOX)
	mlSwitch.SetCallback(mainForm.MLSwitch)
	mainForm.MLSwitchButton = mlSwitch
}

func main() {
	fltk.InitStyles()
	win := fltk.NewWindowWithPosition(450, 450, WIDTH, HEIGHT, "寄存器工具")
	win.SetColor(fltk.WHITE)
	NewMainForm()
	win.End()
	win.Show()
	fltk.Run()
}
