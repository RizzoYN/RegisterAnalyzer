package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/pwiecz/go-fltk"
	"github.com/lxn/win"
)

func SetConsoleTitle(title string) {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)
	_SetConsoleTitle, _ := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	syscall.Syscall(_SetConsoleTitle, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
}

func FindWindow(str string) win.HWND {
    hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(str))
    return hwnd
}

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}

func SetWindowPos(hWnd win.HWND, hWndInsertAfter, x, y, Width, Height, flags int) {
	syscall.NewLazyDLL(`User32.dll`).NewProc(`SetWindowPos`).Call(uintptr(hWnd), uintptr(hWndInsertAfter), uintptr(x), uintptr(y), uintptr(Width), uintptr(Height), uintptr(flags))
}

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
	Row       = 1
	WIDTH     = dataWidth*(bitW+pad) + pad*8 + bitW*13 + 50
	HEIGHT    = bitW + Row*bitH + pad*(3+Row) + 30
	MaxNum    = int64(math.Pow(2, float64(dataWidth)) - 1)
	StartX    = GetSystemMetrics(0)/2 - WIDTH/2
	StartY    = GetSystemMetrics(1)/2 - HEIGHT/2
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
	BitLocs      []*Bit
	Num          *fltk.TextEditor
	LShift       *fltk.Button
	ShiftNum     *fltk.TextEditor
	RShift       *fltk.Button
	Reverse      *fltk.Button
	Invert       *fltk.Button
	Clear        *fltk.Button
	base         int
	lastNum      int64
	lastShiftNum int64
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
		b.Num.Buffer().SetText(fmt.Sprintf("%x", num))
	case 10:
		b.Num.Buffer().SetText(fmt.Sprint(num))
	case 8:
		b.Num.Buffer().SetText(fmt.Sprintf("%o", num))
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
		if fn != nil {
			fn()
		}
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
	curNum := b.Num.Buffer().Text()
	if curNum == "" {
		b.lastNum = 0
	}
	curShiftNum := b.ShiftNum.Buffer().Text()
	if curShiftNum == "" {
		b.lastShiftNum = 0
	}
	num, err := strconv.ParseInt(b.Num.Buffer().Text(), b.base, dataWidth*2)
	if err == nil {
		b.lastNum = num
	} else {
		b.SetNum(b.lastNum)
		b.Num.SetInsertPosition(len(fmt.Sprint(b.lastNum)))
	}
	shiftNum, err := strconv.ParseInt(b.ShiftNum.Buffer().Text(), 10, dataWidth*2)
	if err == nil {
		b.lastShiftNum = shiftNum
	} else {
		if b.lastShiftNum != 0 {
			b.ShiftNum.Buffer().SetText(fmt.Sprint(b.lastShiftNum))
		} else {
			b.ShiftNum.Buffer().SetText("")
		}
	}
	return b.lastNum, int(b.lastShiftNum)
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

func (b *BitRow) KeyTyped(fn func()) func(fltk.Event) bool {
	return func(e fltk.Event) bool {
		if e == fltk.KEYUP {
			num, _ := b.GetCurrentNum()
			b.UpdateBit(num)
			fn()
			return true
		}
		return false
	}
}

func (b *BitRow) ShiftNumTyped(e fltk.Event) bool {
	if e == fltk.KEYUP {
		b.GetCurrentNum()
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
	num.SetWrapMode(fltk.WRAP_NONE)
	num.SetEventHandler(bitRow.KeyTyped(fn))
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
	shiftNum.SetWrapMode(fltk.WRAP_NONE)
	shiftNum.SetEventHandler(bitRow.ShiftNumTyped)
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
	bitRow.Invert = invert
	clear := fltk.NewButton(dataWidth*(bitW+pad)+pad*7+bitW*11+50, h, bitW*2, bitH, "清空")
	clear.SetBox(fltk.BORDER_BOX)
	clear.SetLabelSize(12)
	clear.SetLabelFont(fltk.HELVETICA)
	clear.ClearVisibleFocus()
	clear.SetDownBox(fltk.FLAT_BOX)
	clear.SetCallback(bitRow.ClickClear(fn))
	bitRow.Clear = clear
	bitRow.base = 16
	bitRow.lastNum = 0
	bitRow.lastShiftNum = 1
	return bitRow
}

type MainForm struct {
	Headers Headers
	BitRows []*BitRow
	AddRow  *fltk.Button
	RmRow   *fltk.Button
	Base16  *fltk.RadioRoundButton
	Base10  *fltk.RadioRoundButton
	Base8   *fltk.RadioRoundButton
	OnTop   *fltk.CheckButton
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

func (m *MainForm) Add(w *fltk.Window) func() {
	return func() {
		Row++
		m.RmRow.Activate()
		if Row == maxRow {
			m.AddRow.Deactivate()
		}
		HEIGHT = bitW + Row*bitH + pad*(3+Row) + 30
		w.Resize(w.X(), w.Y(), WIDTH, HEIGHT)
		bitRow := m.BitRows[Row-1]
		for _, obj := range bitRow.BitLocs {
			obj.Show()
		}
		bitRow.Num.Show()
		bitRow.LShift.Show()
		bitRow.ShiftNum.Show()
		bitRow.RShift.Show()
		bitRow.Reverse.Show()
		bitRow.Invert.Show()
		bitRow.Clear.Show()
		m.Updateheaders()
	}
}

func (m *MainForm) Remove(w *fltk.Window) func() {
	return func() {
		Row--
		m.AddRow.Activate()
		if Row == 1 {
			m.RmRow.Deactivate()
		}
		HEIGHT = bitW + Row*bitH + pad*(3+Row) + 30
		w.Resize(w.X(), w.Y(), WIDTH, HEIGHT)
		bitRow := m.BitRows[Row]
		bitRow.ClickClear(nil)()
		for _, obj := range bitRow.BitLocs {
			obj.Hide()
		}
		bitRow.Num.Hide()
		bitRow.LShift.Hide()
		bitRow.ShiftNum.Hide()
		bitRow.RShift.Hide()
		bitRow.Reverse.Hide()
		bitRow.Invert.Hide()
		bitRow.Clear.Hide()
		m.Updateheaders()
	}
}

func (m *MainForm) BaseChoise(base int) func() {
	return func() {
		for r := 0; r < maxRow; r++ {
			num, _ := m.BitRows[r].GetCurrentNum()
			m.BitRows[r].base = base
			m.BitRows[r].SetNum(num)
		}
	}
}

func (m *MainForm) SetOntop(w *fltk.Window) func() {
	return func() {
		p := FindWindow("寄存器工具use go-fltk")
		if m.OnTop.Value() {
			SetWindowPos(p, 0, StartX, StartY, WIDTH, HEIGHT, 2)
		} else {
			SetWindowPos(p, -2, StartX, StartY, WIDTH, HEIGHT, 2)
		}
	}
}

func NewMainForm(w *fltk.Window) {
	mainForm := new(MainForm)
	bitRows := make([]*BitRow, maxRow)
	for r := 0; r <= maxRow; r++ {
		if r == 0 {
			mainForm.Headers = NewHeaders()
		} else {
			bitRow := NewBitRow(r, mainForm.Updateheaders)
			if r > Row {
				for _, obj := range bitRow.BitLocs {
					obj.Hide()
				}
				bitRow.Num.Hide()
				bitRow.LShift.Hide()
				bitRow.ShiftNum.Hide()
				bitRow.RShift.Hide()
				bitRow.Reverse.Hide()
				bitRow.Invert.Hide()
				bitRow.Clear.Hide()
			}
			bitRows[r-1] = bitRow
		}
	}
	mainForm.BitRows = bitRows
	box := fltk.NewBox(fltk.BORDER_BOX, pad*5+30, pad*4, 118, 20, "进制")
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
	addR := fltk.NewButton(pad*5+150, pad*4, 60, 20, "增加一行")
	addR.SetBox(fltk.BORDER_BOX)
	addR.SetLabelSize(12)
	addR.SetLabelFont(fltk.HELVETICA)
	addR.ClearVisibleFocus()
	addR.SetDownBox(fltk.FLAT_BOX)
	addR.SetCallback(mainForm.Add(w))
	rmR := fltk.NewButton(pad*6+210, pad*4, 60, 20, "删除一行")
	rmR.SetBox(fltk.BORDER_BOX)
	rmR.SetLabelSize(12)
	rmR.SetLabelFont(fltk.HELVETICA)
	rmR.ClearVisibleFocus()
	rmR.SetDownBox(fltk.FLAT_BOX)
	rmR.Deactivate()
	rmR.SetCallback(mainForm.Remove(w))
	mainForm.AddRow = addR
	mainForm.RmRow = rmR
	onTop := fltk.NewCheckButton(WIDTH-pad-20, pad*4, 1206, 20, "置顶")
	onTop.ClearVisibleFocus()
	onTop.SetAlign(fltk.ALIGN_LEFT)
	onTop.SetCallback(mainForm.SetOntop(w))
	onTop.Deactivate()
	mainForm.OnTop = onTop
}

func main() {
	SetConsoleTitle("寄存器工具use go-fltk")
	fltk.InitStyles()
	win := fltk.NewWindowWithPosition(StartX, StartY, WIDTH, HEIGHT)
	win.SetLabel("寄存器工具")
	win.SetColor(fltk.WHITE)
	NewMainForm(win)
	win.End()
	win.Show()
	fltk.Run()
}
