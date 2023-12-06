package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"syscall"

	"github.com/pwiecz/go-fltk"
)

var (
	bitColorMap = map[string]fltk.Color{
		"0": fltk.WHITE,
		"1": fltk.BACKGROUND_COLOR,
	}
	headerColorMap = map[int]fltk.Color{
		11: fltk.BLACK,
		14: fltk.RED,
	}
	headerFontMap = map[int]fltk.Font{
		11: fltk.HELVETICA,
		14: fltk.HELVETICA_BOLD,
	}
	textMap = map[string]string{
		"0": "1",
		"1": "0",
	}
	MLmap = map[string]string{
		"MSB": "LSB",
		"LSB": "MSB",
	}
	pad                     = 2
	bitW                    = 16
	bitH                    = 18
	dataWidth               = 64
	maxRow                  = 5
	Row                     = 1
	DisplayNumW             = bitW * 6 * int(dataWidth/32)
	ShiftNumW               = bitW
	ButtonW                 = bitW * 2
	WIDTH                   = dataWidth*bitW + ButtonW*5 + ShiftNumW + DisplayNumW + pad*(dataWidth/4*5+8)
	HEIGHT                  = bitW + Row*bitH + pad*(3+Row) + 28
	maxHeight               = bitW + maxRow*bitH + pad*(3+maxRow) + 42 + bitH
	MaxNum, _               = big.NewInt(0).SetString(strings.Repeat("1", dataWidth), 2)
	user32DLL               = syscall.NewLazyDLL("User32.dll")
	procGetSystemMetrics    = user32DLL.NewProc("GetSystemMetrics")
	procGetSystemMenu       = user32DLL.NewProc("GetSystemMenu")
	procDeleteMenu          = user32DLL.NewProc("DeleteMenu")
	procDrawMenuBar         = user32DLL.NewProc("DrawMenuBar")
	procGetForegroundWindow = user32DLL.NewProc("GetForegroundWindow")
	procSetWindowPos        = user32DLL.NewProc("SetWindowPos")
	MonitorX, _, _          = procGetSystemMetrics.Call(uintptr(0))
	MonitorY, _, _          = procGetSystemMetrics.Call(uintptr(1))
	StartX                  = int(MonitorX)/2 - WIDTH/2
	StartY                  = int(MonitorY)/2 - HEIGHT/2
	svg                     = `<svg version="1.1" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m3.5 1.5c-1.11 0-2 .89-2 2v11c0 1.11.89 2 2 2h2v-2h-2v-11h11v2h2v-2c0-1.11-.89-2-2-2h-11m6 6c-1.11 0-2 .89-2 2v2h2v-2h2v-2h-2m4 0v2h1v1h2v-3h-3m5 0v2h2v11h-11v-2h-2v2c0 1.11.89 2 2 2h11c1.11 0 2-.89 2-2v-11c0-1.11-.89-2-2-2h-2m-4 5v2h-2v2h2c1.11 0 2-.89 2-2v-2h-2m-7 1v3h3v-2h-1v-1z" style="fill:#42a5f5"/></svg>`
)

func NewButton(x, y, w, h int, label string, call func()) *fltk.Button {
	button := fltk.NewButton(x, y, w, h, label)
	button.SetBox(fltk.GLEAM_UP_BOX)
	button.ClearVisibleFocus()
	button.SetLabelSize(12)
	button.SetLabelFont(fltk.HELVETICA)
	button.SetDownBox(fltk.GLEAM_DOWN_BOX)
	button.SetCallback(call)
	return button
}

func NewToggleButton(x, y, w, h int, label string) *fltk.ToggleButton {
	button := fltk.NewToggleButton(x, y, w, h, label)
	button.SetBox(fltk.GLEAM_UP_BOX)
	button.ClearVisibleFocus()
	button.SetLabelSize(12)
	button.SetLabelFont(fltk.HELVETICA)
	button.SetDownBox(fltk.GLEAM_DOWN_BOX)
	return button
}

func NewInput(x, y, w, h int, label string) *fltk.Input {
	input := fltk.NewInput(x, y, w, h)
	input.SetValue(label)
	input.SetBox(fltk.BORDER_BOX)
	return input
}

func NewOutput(x, y, w, h int, label string) *fltk.Output {
	output := fltk.NewOutput(x, y, w, h)
	output.SetBox(fltk.BORDER_BOX)
	output.ClearVisibleFocus()
	output.SetLabelSize(12)
	output.SetLabel(label)
	output.SetAlign(fltk.ALIGN_TOP)
	return output
}

func NewBox(boxType fltk.BoxType, x, y, w, h, labelSize int, label string, bgColor fltk.Color) *fltk.Box {
	box := fltk.NewBox(boxType, x, y, w, h, label)
	box.SetAlign(fltk.ALIGN_CENTER)
	box.SetColor(bgColor)
	box.SetLabelSize(labelSize)
	box.SetLabelFont(fltk.HELVETICA)
	return box
}

func NewGroup(x, y, w, h int) *fltk.Group {
	group := fltk.NewGroup(x, y, w, h)
	return group
}

func NewRadioRoundButton(x, y, w, h, base int, label string, f func(int) func()) *fltk.RadioRoundButton {
	button := fltk.NewRadioRoundButton(x, y, w, h, label)
	button.ClearVisibleFocus()
	button.SetCallback(f(base))
	return button
}

func ParseHeight(row int) int {
	if row == 1 {
		return pad*2 + 28 + bitW
	} else {
		return (row-1)*bitH + pad*(row+1) + 28 + bitW
	}
}

func SetOntop(ontop bool) {
	swpNoSize := 0x1
	swpNoMove := 0x2
	flag := swpNoSize | swpNoMove
	hwnd, _, _ := procGetForegroundWindow.Call()
	if ontop {
		topMost := -1
		procSetWindowPos.Call(hwnd, uintptr(topMost), uintptr(0), uintptr(0), uintptr(0), uintptr(0), uintptr(flag))
	} else {
		bottom := 1
		top := 0
		procSetWindowPos.Call(hwnd, uintptr(bottom), uintptr(0), uintptr(0), uintptr(0), uintptr(0), uintptr(flag))
		procSetWindowPos.Call(hwnd, uintptr(top), uintptr(0), uintptr(0), uintptr(0), uintptr(0), uintptr(flag))
	}
}

func DisableMenuAndFullScreen() {
	hwnd, _, _ := procGetForegroundWindow.Call()
	hmenu, _, _ := procGetSystemMenu.Call(hwnd, uintptr(0))
	procDeleteMenu.Call(hmenu, uintptr(0xF030), uintptr(0))
	procDrawMenuBar.Call(hwnd)
}

type Bit struct {
	*fltk.Box
}

func (b *Bit) Click() {
	val := b.Label()
	str := textMap[val]
	b.SetLabel(str)
	b.SetColor(bitColorMap[str])
}

func NewBit(x, y, w, h int) *Bit {
	bit := NewBox(fltk.BORDER_BOX, x, y, w, h, 14, "0", fltk.WHITE)
	return &Bit{bit}
}

type BitRow struct {
	group           *fltk.Group
	bitLocs         []*Bit
	num             *fltk.Input
	lShift          *fltk.Button
	shiftNum        *fltk.Input
	rShift          *fltk.Button
	reverse         *fltk.Button
	invert          *fltk.Button
	clear           *fltk.Button
	base            int
	lastShiftNum    int64
	shiftNumDisplay *fltk.Box
	bigInt          big.Int
}

func (b *BitRow) GetBitString() []string {
	bitList := make([]string, dataWidth)
	for c := 0; c < dataWidth; c++ {
		bitList[c] = b.bitLocs[c].Label()
	}
	return bitList
}

func (b *BitRow) SetNum() {
	b.num.SetValue(b.bigInt.Text(b.base))
}

func (b *BitRow) GetCurrentNum() int64 {
	shiftNum, err := strconv.ParseInt(b.shiftNum.Value(), 10, dataWidth*2)
	if err == nil {
		b.lastShiftNum = shiftNum
	} else {
		if b.lastShiftNum != 0 {
			b.shiftNum.SetValue(fmt.Sprint(b.lastShiftNum))
		} else {
			b.shiftNum.SetValue("")
		}
	}
	return b.lastShiftNum
}

func (b *BitRow) UpdateBit() {
	bins := fmt.Sprintf("%0*s", dataWidth, b.bigInt.Text(2))
	for c := 0; c < dataWidth; c++ {
		s := string(bins[c])
		b.bitLocs[c].SetLabel(s)
		b.bitLocs[c].SetColor(bitColorMap[s])
	}
}

func (b *BitRow) UpdateNum() {
	binStr := b.GetBitString()
	str := strings.Join(binStr, "")
	binInt, _ := b.bigInt.SetString(str, 2)
	b.bigInt = *binInt
}

func (b *BitRow) UpdateBitNum() {
	b.SetNum()
	b.UpdateBit()
}

func (b *BitRow) Display() {
	b.shiftNum.Hide()
	b.shiftNumDisplay.Show()
}

func (b *BitRow) Hide() {
	b.group.Hide()
}

func (b *BitRow) Show() {
	b.group.Show()
}

func (b *BitRow) DisplayClick(e fltk.Event) bool {
	if e == fltk.Event(fltk.LeftMouse) {
		b.shiftNumDisplay.Hide()
		b.shiftNum.SetValue("")
		b.shiftNum.Show()
		b.shiftNum.TakeFocus()
		return true
	}
	return false
}

func (b *BitRow) ClickLShift(fn, fnc func()) func() {
	return func() {
		shiftNum := b.GetCurrentNum()
		b.bigInt.Lsh(&b.bigInt, uint(shiftNum))
		b.bigInt.And(&b.bigInt, MaxNum)
		b.UpdateBitNum()
		b.Display()
		fn()
		fnc()
	}
}

func (b *BitRow) ClickRShift(fn, fnc func()) func() {
	return func() {
		shiftNum := b.GetCurrentNum()
		b.bigInt.Rsh(&b.bigInt, uint(shiftNum))
		b.bigInt.And(&b.bigInt, MaxNum)
		b.UpdateBitNum()
		fn()
		fnc()
		b.Display()
	}
}

func (b *BitRow) ClickReverse(fn, fnc func()) func() {
	return func() {
		binStr := make([]string, dataWidth)
		for i, j := 0, dataWidth-1; i < j; i, j = i+1, j-1 {
			h := b.bitLocs[i].Label()
			e := b.bitLocs[j].Label()
			b.bitLocs[i].SetLabel(e)
			b.bitLocs[i].SetColor(bitColorMap[e])
			b.bitLocs[j].SetLabel(h)
			b.bitLocs[j].SetColor(bitColorMap[h])
			binStr[i] = b.bitLocs[i].Label()
		}
		b.UpdateNum()
		b.SetNum()
		fn()
		fnc()
		b.Display()
	}
}

func (b *BitRow) KeyType(fn, fnc func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			b.Display()
		}
		if e == fltk.KEYUP {
			old := b.bigInt.Text(b.base)
			str := b.num.Value()
			if str == "" {
				b.bigInt.SetString("0", b.base)
			} else {
				bigInt, ok := b.bigInt.SetString(str, b.base)
				if !ok || bigInt.Cmp(MaxNum) == 1 {
					b.bigInt.SetString(old, b.base)
				}
			}
			b.UpdateBit()
			b.SetNum()
			fn()
			fnc()
			b.Display()
			b.SetNum()
			return true
		}
		return false
	}
}

func (b *BitRow) ShiftNumEvent(e fltk.Event) bool {
	if e == fltk.Event(fltk.LeftMouse) && b.shiftNum.HasFocus() {
		b.shiftNumDisplay.Show()
		b.shiftNum.SetValue(fmt.Sprint(b.lastShiftNum))
		b.shiftNum.Hide()
	}
	if e == fltk.KEYUP {
		shiftNum := b.GetCurrentNum()
		b.shiftNumDisplay.SetLabel(fmt.Sprint(shiftNum))
		return true
	}
	return false
}

func (b *BitRow) Click(fn, fnc, fun func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.Event(fltk.LeftMouse) {
			fn()
			fnc()
			fun()
			b.UpdateNum()
			b.SetNum()
			b.Display()
			return true
		}
		return false
	}
}

func (b *BitRow) ClickClear(fn, fnc func()) func() {
	return func() {
		for c := 0; c < dataWidth; c++ {
			b.bitLocs[c].SetLabel("0")
			b.bitLocs[c].SetColor(bitColorMap["0"])
		}
		b.num.SetValue("0")
		if fn != nil {
			fn()
		}
		fnc()
		b.Display()
	}
}

func (b *BitRow) ClickInvert(fn, fnc func()) func() {
	return func() {
		for c := 0; c < dataWidth; c++ {
			b.bitLocs[c].Click()
		}
		b.UpdateNum()
		b.SetNum()
		fn()
		fnc()
		b.Display()
	}
}

func NewBitRow(row int, fn, fnc func()) *BitRow {
	bitRow := new(BitRow)
	h := ParseHeight(row)
	group := fltk.NewGroup(0, h, WIDTH, bitH)
	bitRow.group = group
	bitLocs := make([]*Bit, dataWidth)
	for c := 0; c < dataWidth; c++ {
		n := c*bitW + c/4*pad*2 + (c+1)*pad
		if n == 1 {
			n = 4
		}
		bit := NewBit(n, h, bitW, bitH)
		bit.SetEventHandler(bitRow.Click(bit.Click, fn, fnc))
		bitLocs[c] = bit
	}
	bitsWidth := dataWidth*bitW + dataWidth/4*pad*2 + (dataWidth+1)*pad
	bitRow.bitLocs = bitLocs
	num := NewInput(bitsWidth, h, DisplayNumW, bitH, "0")
	num.SetEventHandler(bitRow.KeyType(fn, fnc))
	bitRow.num = num
	lShift := NewButton(bitsWidth+pad+DisplayNumW, h, 25, bitH, "<<", bitRow.ClickLShift(fn, fnc))
	bitRow.lShift = lShift
	shiftNum := NewInput(bitsWidth+pad*2+DisplayNumW+25, h, bitW, bitH, "1")
	shiftNum.SetEventHandler(bitRow.ShiftNumEvent)
	shiftNum.Hide()
	bitRow.shiftNum = shiftNum
	rShift := NewButton(bitsWidth+pad*3+DisplayNumW+bitW+25, h, 25, bitH, ">>", bitRow.ClickRShift(fn, fnc))
	bitRow.rShift = rShift
	reverse := NewButton(bitsWidth+pad*4+DisplayNumW+bitW+50, h, bitW*2, bitH, "倒序", bitRow.ClickReverse(fn, fnc))
	bitRow.reverse = reverse
	invert := NewButton(bitsWidth+pad*5+DisplayNumW+bitW*2+50, h, bitW*2, bitH, "转换", bitRow.ClickInvert(fn, fnc))
	bitRow.invert = invert
	clear := NewButton(bitsWidth+pad*6+DisplayNumW+bitW*4+50, h, bitW*2, bitH, "清空", bitRow.ClickClear(fn, fnc))
	bitRow.clear = clear
	bitRow.base = 16
	bitRow.lastShiftNum = 1
	shiftDisplay := NewBox(fltk.BORDER_BOX, bitsWidth+pad*2+bitW*6+25, h, bitW, bitH, 14, fmt.Sprint(bitRow.lastShiftNum), fltk.WHITE)
	shiftDisplay.SetEventHandler(bitRow.DisplayClick)
	bitRow.shiftNumDisplay = shiftDisplay
	bitRow.bigInt = *big.NewInt(0)
	group.End()
	return bitRow
}

type Header struct {
	*fltk.Box
}

func NewHeader(x, y, w, h, ix int) *Header {
	header := NewBox(fltk.FLAT_BOX, x, y, w, h, 11, fmt.Sprint(ix), fltk.WHITE)
	return &Header{header}
}

type Headers []*Header

func (h Headers) UpdateHeader(c, size int) {
	h[c].SetLabelColor(headerColorMap[size])
	h[c].SetLabelSize(size)
	h[c].SetLabelFont(headerFontMap[size])
	h[c].Redraw()
}

func NewHeaders() Headers {
	headers := make([]*Header, dataWidth)
	for c := 0; c < dataWidth; c++ {
		n := c*bitW + (c/4)*pad*2 + (c+1)*pad
		if n == 1 {
			n = 4
		}
		head := NewHeader(n, pad*2+28, bitW, bitW, dataWidth-1-c)
		headers[c] = head
	}
	return headers
}

type ColorSelect struct {
	group  *fltk.Group
	colors []*fltk.Button
	index  int
}

func (c *ColorSelect) Click(m *MainForm) func() {
	return func() {
		var color fltk.Color
		for _, box := range c.colors {
			if box.HasFocus() {
				color = box.Color()
				break
			}
		}
		if c.index == 0 {
			m.BitColorBox.SetColor(color)
			m.BitColorBox.Redraw()
			bitColorMap["1"] = color
			for r := 0; r < Row; r++ {
				bitRow := m.BitRows[r]
				bitRow.UpdateBit()
			}
		} else {
			headerColorMap[14] = color
			m.HeaderColorBox.SetColor(color)
			m.HeaderColorBox.Redraw()
			m.Updateheaders()
		}
		c.group.Hide()
	}
}

func NewColorSelect(m *MainForm) *ColorSelect {
	group := fltk.NewGroup(pad*10+307, 0, WIDTH-pad*10-504, 30)
	colorSelect := new(ColorSelect)
	colors := make([]*fltk.Button, 52)
	colorCode := []fltk.Color{
		fltk.BACKGROUND_COLOR, fltk.INACTIVE_COLOR, fltk.DARK_BLUE,
		fltk.DARK_CYAN, fltk.DARK_GREEN, fltk.DARK_MAGENTA,
		fltk.DARK_RED, fltk.DARK_YELLOW, fltk.LIGHT2,
		fltk.SELECTION_COLOR, fltk.BLUE, fltk.CYAN,
		fltk.GREEN, fltk.MAGENTA, fltk.RED, fltk.Color(0xFFA50000),
		fltk.YELLOW, fltk.Color(0x008B8B00), fltk.Color(0x00CED100),
		fltk.Color(0x00FFFF00), fltk.Color(0xFFE4E100), fltk.Color(0xFF8C0000),
		fltk.Color(0x1E90FF00), fltk.Color(0x20B2AA00), fltk.Color(0x2F4F4F00),
		fltk.Color(0x483D8B00), fltk.Color(0x4B008200), fltk.Color(0xBA55D300),
		fltk.Color(0x6B8E2300), fltk.Color(0xADFF2F00), fltk.Color(0x80008000),
		fltk.Color(0x8A2BE200), fltk.Color(0x8B451300), fltk.Color(0x8FBC8F00),
		fltk.Color(0x90EE9000), fltk.Color(0x9932CC00), fltk.Color(0x9400D300),
		fltk.Color(0x9ACD3200), fltk.Color(0xA0522D00), fltk.Color(0xA52A2A00),
		fltk.Color(0xADD8E600), fltk.Color(0xFFD70000),
		// fltk.Color(0xB8860B00),
		// fltk.Color(0x00800000), fltk.Color(0xC7158500), fltk.Color(0x00BFFF00),
		// fltk.Color(0x00Fe7e00), fltk.Color(0xEEE8AA00), fltk.Color(0x69696900),
		// fltk.Color(0xFFFFF000), fltk.Color(0xFFEFD500), fltk.Color(0xB2222200),
	}
	idx := 0
	for r := 0; r < 2; r++ {
		for c := 0; c < 21; c++ {
			box := fltk.NewButton((pad+15)*c+pad*13+304, pad*r+r*15, 15, 15)
			box.SetColor(colorCode[idx])
			box.SetBox(fltk.BORDER_BOX)
			box.SetDownBox(fltk.BORDER_BOX)
			box.SetCallback(colorSelect.Click(m))
			colors[idx] = box
			idx++
		}
	}
	colorSelect.colors = colors
	colorSelect.group = group
	group.End()
	group.Hide()
	return colorSelect
}

type BitAnalyze struct {
	group *fltk.Group
	res   []*fltk.Output
	input *fltk.Input
}

func NewBitAnalyze() *BitAnalyze {
	width := (WIDTH-pad*2)/(maxRow+1) - pad
	bitAnalyze := new(BitAnalyze)
	group := fltk.NewGroup(0, HEIGHT, WIDTH, bitH+12)
	res := make([]*fltk.Output, maxRow)
	for c := 0; c <= maxRow; c++ {
		if c == 0 {
			input := NewInput(c*(width+pad)+pad*2, HEIGHT, width, bitH, "")
			input.SetLabel("Bit位域")
			input.SetAlign(fltk.ALIGN_TOP)
			input.SetLabelSize(12)
			bitAnalyze.input = input
		} else {
			output := NewOutput(c*(width+pad)+pad*2, HEIGHT, width, bitH, fmt.Sprintf("第%d行", c))
			res[c-1] = output
			if c > Row {
				output.Hide()
			}
		}
	}
	bitAnalyze.group = group
	bitAnalyze.res = res
	group.End()
	group.Hide()
	return bitAnalyze
}

type MainForm struct {
	Group          *fltk.Group
	Headers        Headers
	BitRows        []*BitRow
	AddRow         *fltk.Button
	RmRow          *fltk.Button
	Base16         *fltk.RadioRoundButton
	Base10         *fltk.RadioRoundButton
	Base8          *fltk.RadioRoundButton
	ontop          *fltk.ToggleButton
	base           int
	MLSwitchButton *fltk.ToggleButton
	BitColorSel    *fltk.Button
	BitColorBox    *fltk.Box
	HeaderColorSel *fltk.Button
	HeaderColorBox *fltk.Box
	BitRangeParse  *fltk.ToggleButton
	AnalyzeArea    *BitAnalyze
	ColorSelArea   *ColorSelect
}

func (m *MainForm) Updateheaders() {
	for c := 0; c < dataWidth; c++ {
		bitMap := make(map[string]int, Row)
		for r := 0; r < Row; r++ {
			val := m.BitRows[r].bitLocs[c].Label()
			bitMap[val] = 0
		}
		if len(bitMap) == 1 {
			m.Headers.UpdateHeader(c, 11)
		} else {
			m.Headers.UpdateHeader(c, 14)
		}
	}
}

func (m *MainForm) AnalyzeAreaChange() {
	if m.BitRangeParse.Value() {
		m.AnalyzeArea.group.SetPosition(0, HEIGHT-bitH-pad)
		m.AnalyzeArea.input.SetPosition(m.AnalyzeArea.input.X(), HEIGHT-bitH-pad)
		for idx, output := range m.AnalyzeArea.res {
			output.SetPosition(output.X(), HEIGHT-bitH-pad)
			if idx+1 <= Row {
				output.Show()
			} else {
				output.Hide()
			}
		}
	}
}

func (m *MainForm) Add() {
	Row++
	m.RmRow.Activate()
	if Row == maxRow {
		m.AddRow.Deactivate()
	}
	HEIGHT += bitH + pad
	m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
	m.AnalyzeAreaChange()
	bitRow := m.BitRows[Row-1]
	bitRow.Show()
	m.Updateheaders()
}

func (m *MainForm) Remove() {
	Row--
	m.AddRow.Activate()
	if Row == 1 {
		m.RmRow.Deactivate()
	}
	HEIGHT -= bitH + pad
	m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
	m.AnalyzeAreaChange()
	bitRow := m.BitRows[Row]
	bitRow.ClickClear(nil, m.UpdateAnalyzeArea)()
	bitRow.Hide()
	m.Updateheaders()
}

func (m *MainForm) BaseChoise(base int) func() {
	return func() {
		m.base = base
		for r := 0; r < maxRow; r++ {
			m.BitRows[r].base = base
			m.BitRows[r].UpdateNum()
			m.UpdateAnalyzeRes(r)
			if r < Row {
				m.BitRows[r].Display()
			}
		}
	}
}

func (m *MainForm) SetOnTop() {
	status := m.ontop.Value()
	SetOntop(status)
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
	m.UpdateAnalyzeArea()
}

func (m *MainForm) Analyze() {
	boolean := m.BitRangeParse.Value()
	if boolean {
		m.AnalyzeArea.group.SetPosition(0, HEIGHT)
		m.AnalyzeArea.input.SetPosition(m.AnalyzeArea.input.X(), HEIGHT)
		for _, output := range m.AnalyzeArea.res {
			output.SetPosition(output.X(), HEIGHT)
		}
		HEIGHT += bitH + 14
		m.AnalyzeArea.group.Show()
		m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
	} else {
		m.AnalyzeArea.group.Hide()
		HEIGHT -= bitH + 14
		m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
	}
	m.AnalyzeAreaChange()
}

func (m *MainForm) SetNum(num int64, r int32) {
	switch m.base {
	case 16:
		m.AnalyzeArea.res[r].SetValue(fmt.Sprintf("%x", num))
	case 10:
		m.AnalyzeArea.res[r].SetValue(fmt.Sprint(num))
	case 8:
		m.AnalyzeArea.res[r].SetValue(fmt.Sprintf("%o", num))
	}
}

func (m *MainForm) ParseBitRange(nums []string, r int32) (int64, error) {
	var res []string
	if len(nums) == 1 {
		left := strings.Trim(nums[0], "\r\n")
		num, err := strconv.ParseInt(left, 10, 0)
		if err != nil || num >= int64(dataWidth) {
			return 0, fmt.Errorf("无效输入")
		}
		for c := 0; c < dataWidth; c++ {
			if m.Headers[c].Label() == left {
				str := m.BitRows[r].bitLocs[c].Label()
				res = append(res, str)
			}
		}
		bin, _ := strconv.ParseInt(strings.Join(res, ""), 2, dataWidth*2)
		return bin, nil
	} else if len(nums) == 2 {
		left := strings.Trim(nums[0], "\r\n")
		right := strings.Trim(nums[1], "\r\n")
		numL, errL := strconv.ParseInt(left, 10, 0)
		if errL != nil || numL >= int64(dataWidth) {
			return 0, fmt.Errorf("无效输入")
		}
		numR, errR := strconv.ParseInt(right, 10, 0)
		if errR != nil || numR >= int64(dataWidth) {
			return 0, fmt.Errorf("无效输入")
		}
		if numR == numL {
			return 0, fmt.Errorf("无效输入")
		}
		flag := false
		for c := 0; c < dataWidth; c++ {
			var str string
			cap := m.Headers[c].Label()
			if (cap == left && !flag) || (cap == right && !flag) {
				flag = true
			} else if (cap == right && flag) || (cap == left && flag) {
				flag = false
				str = m.BitRows[r].bitLocs[c].Label()
				res = append(res, str)
			}
			if flag {
				str = m.BitRows[r].bitLocs[c].Label()
				res = append(res, str)
			}
		}
		bin, _ := strconv.ParseInt(strings.Join(res, ""), 2, dataWidth*2)
		return bin, nil
	} else {
		return 0, fmt.Errorf("无效输入")
	}
}

func (m *MainForm) UpdateAnalyzeRes(r int) {
	str := m.AnalyzeArea.input.Value()
	num, err := m.ParseBitRange(strings.Split(str, ":"), int32(r))
	output := m.AnalyzeArea.res[r]
	if err != nil {
		if str != "" {
			output.SetValue("无效输入")
			output.SetColor(fltk.RED)
		} else {
			output.SetValue("")
			output.SetColor(fltk.WHITE)
		}
	} else {
		m.SetNum(num, int32(r))
		output.SetColor(fltk.WHITE)
	}
	output.Redraw()
}

func (m *MainForm) Edit(e fltk.Event) bool {
	if e == fltk.KEYUP {
		for r := 0; r < maxRow; r++ {
			m.UpdateAnalyzeRes(r)
		}
		return true
	}
	return false
}

func (m *MainForm) UpdateAnalyzeArea() {
	m.Edit(fltk.KEYUP)
}

func NewMainForm(w *fltk.Window) {
	mainForm := new(MainForm)
	mainForm.base = 16
	bitRows := make([]*BitRow, maxRow)
	for r := 0; r <= maxRow; r++ {
		if r == 0 {
			mainForm.Headers = NewHeaders()
		} else {
			bitRow := NewBitRow(r, mainForm.Updateheaders, mainForm.UpdateAnalyzeArea)
			if r > Row {
				bitRow.Hide()
			}
			bitRows[r-1] = bitRow
		}
	}
	mainForm.BitRows = bitRows
	box := NewBox(fltk.GTK_UP_BOX, WIDTH-189, 18, 118, 25, 12, "进制", fltk.WHITE)
	box.SetAlign(fltk.ALIGN_TOP)
	base16 := NewRadioRoundButton(WIDTH-184, pad*11+1, 16, 16, 16, "16", mainForm.BaseChoise)
	mainForm.Base16 = base16
	base10 := NewRadioRoundButton(WIDTH-144, pad*11+1, 16, 16, 10, "10", mainForm.BaseChoise)
	mainForm.Base10 = base10
	base8 := NewRadioRoundButton(WIDTH-104, pad*11+1, 16, 16, 8, "8", mainForm.BaseChoise)
	mainForm.Base8 = base8
	addR := NewButton(WIDTH-66, pad-1, 60, 20, "增加一行", mainForm.Add)
	rmR := NewButton(WIDTH-66, pad+21, 60, 20, "删除一行", mainForm.Remove)
	rmR.Deactivate()
	mainForm.AddRow = addR
	mainForm.RmRow = rmR
	ontop := NewToggleButton(pad*3, pad*4, 35, 20, "置顶")
	ontop.SetCallback(mainForm.SetOnTop)
	base16.SetValue(true)
	mainForm.ontop = ontop
	mlSwitch := NewToggleButton(pad*4+35, pad*4, 35, 20, "MSB")
	mlSwitch.SetCallback(mainForm.MLSwitch)
	rangeParse := NewToggleButton(pad*5+70, pad*4, 60, 20, "位域解析")
	analyzeArea := NewBitAnalyze()
	analyzeArea.input.SetEventHandler(mainForm.Edit)
	rangeParse.SetCallback(mainForm.Analyze)
	bitColorBox := fltk.NewBox(fltk.BORDER_BOX, pad*9+130, pad*6, 12, 12)
	bitColorBox.SetColor(bitColorMap["1"])
	colorDia := NewColorSelect(mainForm)
	callBack := func(i int) func() {
		return func() {
			if colorDia.group.Visible() {
				colorDia.group.Hide()
				colorDia.index = i
			} else {
				colorDia.group.Show()
				colorDia.index = i
			}
		}
	}
	bitColorSel := NewButton(pad*6+150, pad*4, 60, 20, "颜色选择", callBack(0))
	headerColorBox := fltk.NewBox(fltk.BORDER_BOX, pad*10+210, pad*6, 12, 12)
	headerColorBox.SetColor(headerColorMap[14])
	headerColorSel := NewButton(pad*7+230, pad*4, 80, 20, "对比颜色选择", callBack(1))
	mainForm.BitColorSel = bitColorSel
	mainForm.BitColorBox = bitColorBox
	mainForm.HeaderColorSel = headerColorSel
	mainForm.HeaderColorBox = headerColorBox
	mainForm.MLSwitchButton = mlSwitch
	mainForm.BitRangeParse = rangeParse
	mainForm.ColorSelArea = colorDia
	mainForm.AnalyzeArea = analyzeArea
	mainForm.Group = &w.Group
}

func main() {
	fltk.InitStyles()
	icon, _ := fltk.NewSvgImageFromString(svg)
	win := fltk.NewWindowWithPosition(StartX, StartY, WIDTH, HEIGHT, "寄存器工具")
	win.SetColor(fltk.WHITE)
	win.SetSizeRange(WIDTH, HEIGHT, WIDTH, maxHeight, 0, 0, false)
	NewMainForm(win)
	win.SetIcons([]*fltk.RgbImage{&icon.RgbImage})
	win.End()
	win.Show()
	DisableMenuAndFullScreen()
	fltk.Run()
}
