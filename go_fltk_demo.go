package main

import (
    "fmt"
    "math"
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
    pad                  = 2
    bitW                 = 18
    bitH                 = 22
    dataWidth            = 32
    maxRow               = 5
    Row                  = 1
    WIDTH                = dataWidth*bitW + dataWidth/4*pad*2 + (dataWidth+1)*pad + pad*7 + bitW*13 + 50
    HEIGHT               = bitW + Row*bitH + pad*(3+Row) + 30
    maxHeight            = bitW + maxRow*bitH + pad*(3+maxRow) + 30
    MaxNum               = int64(math.Pow(2, float64(dataWidth)) - 1)
    user32DLL            = syscall.NewLazyDLL("User32.dll")
    procGetSystemMetrics = user32DLL.NewProc("GetSystemMetrics")
    MonitorX, _, _       = procGetSystemMetrics.Call(uintptr(0))
    MonitorY, _, _       = procGetSystemMetrics.Call(uintptr(1))
    StartX               = int(MonitorX)/2 - WIDTH/2
    StartY               = int(MonitorY)/2 - HEIGHT/2
)

func NewButton(x, y, w, h int, label string) *fltk.Button {
    button := fltk.NewButton(x, y, w, h, label)
    button.SetBox(fltk.GLEAM_UP_BOX)
    button.ClearVisibleFocus()
    button.SetLabelSize(12)
    button.SetLabelFont(fltk.HELVETICA)
    button.SetDownBox(fltk.GLEAM_DOWN_BOX)
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
    procSetWindowPos := user32DLL.NewProc("SetWindowPos")
    procGetForegroundWindow := user32DLL.NewProc("GetForegroundWindow")
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
    lastNum         int64
    lastShiftNum    int64
    shiftNumDisplay *fltk.Box
}

func (b *BitRow) GetBitString() []string {
    bitList := make([]string, dataWidth)
    for c := 0; c < dataWidth; c++ {
        bitList[c] = b.bitLocs[c].Label()
    }
    return bitList
}

func (b *BitRow) SetNum(num int64) {
    switch b.base {
    case 16:
        b.num.SetValue(fmt.Sprintf("%x", num))
    case 10:
        b.num.SetValue(fmt.Sprint(num))
    case 8:
        b.num.SetValue(fmt.Sprintf("%o", num))
    }
}

func (b *BitRow) UpdateNum() {
    bitList := b.GetBitString()
    binStr := strings.Join(bitList, "")
    bin, _ := strconv.ParseInt(binStr, 2, dataWidth*2)
    b.SetNum(bin)
}

func (b *BitRow) GetCurrentNum() (int64, int64) {
    curNum := b.num.Value()
    if curNum == "" {
        b.lastNum = 0
    }
    num, err := strconv.ParseInt(b.num.Value(), b.base, dataWidth*2)
    if err == nil {
        b.lastNum = num
    } else {
        b.SetNum(b.lastNum)
    }
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
    return b.lastNum, b.lastShiftNum
}

func (b *BitRow) UpdateBit(num int64) {
    str := fmt.Sprintf("%0*b", dataWidth, num)
    for c := 0; c < dataWidth; c++ {
        s := string(str[c])
        b.bitLocs[c].SetLabel(s)
        b.bitLocs[c].SetColor(bitColorMap[s])
    }
}

func (b *BitRow) UpdateBitNum(num int64) {
    b.SetNum(num)
    b.UpdateBit(num)
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
        for i, j := 0, len(b.bitLocs)-1; i < j; i, j = i+1, j-1 {
            h := b.bitLocs[i].Label()
            e := b.bitLocs[j].Label()
            b.bitLocs[i].SetLabel(e)
            b.bitLocs[i].SetColor(bitColorMap[e])
            b.bitLocs[j].SetLabel(h)
            b.bitLocs[j].SetColor(bitColorMap[h])
        }
        b.UpdateNum()
        fn()
        b.Display()
    }
}

func (b *BitRow) KeyType(fn func()) func(fltk.Event) bool {
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
    if e == fltk.Event(fltk.LeftMouse) && b.shiftNum.HasFocus() {
        b.shiftNumDisplay.Show()
        b.shiftNum.SetValue(fmt.Sprint(b.lastShiftNum))
        b.shiftNum.Hide()
    }
    if e == fltk.KEYUP {
        _, shiftNum := b.GetCurrentNum()
        b.shiftNumDisplay.SetLabel(fmt.Sprint(shiftNum))
        return true
    }
    return false
}

func (b *BitRow) Click(fn func(), fnc func()) func(fltk.Event) bool {
    return func(e fltk.Event) bool {
        if e == fltk.Event(fltk.LeftMouse) {
            fn()
            fnc()
            b.UpdateNum()
            b.Display()
            return true
        }
        return false
    }
}

func (b *BitRow) ClickClear(fn func()) func() {
    return func() {
        for c := 0; c < dataWidth; c++ {
            b.bitLocs[c].SetLabel("0")
            b.bitLocs[c].SetColor(bitColorMap["0"])
        }
        b.num.SetValue("0")
        if fn != nil {
            fn()
        }
        b.Display()
    }
}

func (b *BitRow) ClickInvert(fn func()) func() {
    return func() {
        for c := 0; c < dataWidth; c++ {
            b.bitLocs[c].Click()
        }
        b.UpdateNum()
        fn()
        b.Display()
    }
}

func NewBitRow(row int, fn func()) *BitRow {
    bitRow := new(BitRow)
    h := ParseHeight(row)
    group := fltk.NewGroup(0, h, WIDTH, bitH, fmt.Sprint(row))
    group.SetLabelType(fltk.NO_LABEL)
    bitRow.group = group
    bitLocs := make([]*Bit, dataWidth)
    for c := 0; c < dataWidth; c++ {
        n := c*bitW + c/4*pad*2 + (c+1)*pad
        if n == 1 {
            n = 4
        }
        bit := NewBit(n, h, bitW, bitH)
        bit.SetEventHandler(bitRow.Click(bit.Click, fn))
        bitLocs[c] = bit
    }
    bitsWidth := dataWidth*bitW + dataWidth/4*pad*2 + (dataWidth+1)*pad
    bitRow.bitLocs = bitLocs
    num := NewInput(bitsWidth, h, bitW*6, bitH, "0")
    num.SetEventHandler(bitRow.KeyType(fn))
    bitRow.num = num
    lShift := NewButton(bitsWidth+pad+bitW*6, h, 25, bitH, "<<")
    lShift.SetCallback(bitRow.ClickLShift(fn))
    bitRow.lShift = lShift
    shiftNum := NewInput(bitsWidth+pad*2+bitW*6+25, h, bitW, bitH, "1")
    shiftNum.SetEventHandler(bitRow.ShiftNumEvent)
    shiftNum.Hide()
    bitRow.shiftNum = shiftNum
    rShift := NewButton(bitsWidth+pad*3+bitW*7+25, h, 25, bitH, ">>")
    rShift.SetCallback(bitRow.ClickRShift(fn))
    bitRow.rShift = rShift
    reverse := NewButton(bitsWidth+pad*4+bitW*7+50, h, bitW*2, bitH, "倒序")
    reverse.SetCallback(bitRow.ClickReverse(fn))
    bitRow.reverse = reverse
    invert := NewButton(bitsWidth+pad*5+bitW*9+50, h, bitW*2, bitH, "转换")
    invert.SetCallback(bitRow.ClickInvert(fn))
    bitRow.invert = invert
    clear := NewButton(bitsWidth+pad*6+bitW*11+50, h, bitW*2, bitH, "清空")
    clear.SetCallback(bitRow.ClickClear(fn))
    bitRow.clear = clear
    bitRow.base = 16
    bitRow.lastNum = 0
    bitRow.lastShiftNum = 1
    shiftDisplay := NewBox(fltk.BORDER_BOX, bitsWidth+pad*2+bitW*6+25, h, bitW, bitH, 14, fmt.Sprint(bitRow.lastShiftNum), fltk.WHITE)
    shiftDisplay.SetEventHandler(bitRow.DisplayClick)
    bitRow.shiftNumDisplay = shiftDisplay
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
            }
        }
        if c.index == 0 {
            m.BitColorBox.SetColor(color)
            m.BitColorBox.Redraw()
            bitColorMap["1"] = color
            for r := 0; r < Row; r++ {
                bitRow := m.BitRows[r]
                num, _ := bitRow.GetCurrentNum()
                bitRow.UpdateBit(num)
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
    group := fltk.NewGroup(pad*10+310, 0, WIDTH-pad*10-540, 30)
    group.SetLabelType(fltk.NO_LABEL)
    colorSelect := new(ColorSelect)
    colors := make([]*fltk.Button, 24)
    colorCode := []fltk.Color{
        fltk.BACKGROUND_COLOR, fltk.INACTIVE_COLOR, fltk.DARK_BLUE, fltk.DARK_CYAN, fltk.DARK_GREEN,
        fltk.DARK_MAGENTA, fltk.DARK_RED, fltk.DARK_YELLOW, fltk.LIGHT2, fltk.SELECTION_COLOR,
        fltk.BLUE, fltk.CYAN, fltk.GREEN, fltk.MAGENTA, fltk.RED, fltk.YELLOW,
        fltk.Color(0x00800000), fltk.Color(0x008B8B00), fltk.Color(0x00BFFF00),
        fltk.Color(0x00Fe7e00), fltk.Color(0x4B008200), fltk.Color(0x69696900),
        fltk.Color(0x77889900), fltk.Color(0x80808000),
    }
    idx := 0
    for r := 0; r < 2; r++ {
        for c := 0; c < 12; c++ {
            box := fltk.NewButton((pad+15)*c+pad*12+310, pad*r+r*15, 15, 15)
            box.SetColor(colorCode[idx])
            box.SetBox(fltk.GLEAM_THIN_UP_BOX)
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

func (m *MainForm) Add() {
    Row++
    m.RmRow.Activate()
    if Row == maxRow {
        m.AddRow.Deactivate()
    }
    HEIGHT = bitW + Row*bitH + pad*(3+Row) + 30
    m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
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
    HEIGHT = bitW + Row*bitH + pad*(3+Row) + 30
    m.Group.Resize(m.Group.X(), m.Group.Y(), WIDTH, HEIGHT)
    bitRow := m.BitRows[Row]
    bitRow.ClickClear(nil)()
    bitRow.Hide()
    m.Updateheaders()
}

func (m *MainForm) BaseChoise(base int) func() {
    return func() {
        for r := 0; r < maxRow; r++ {
            num, _ := m.BitRows[r].GetCurrentNum()
            m.BitRows[r].base = base
            m.BitRows[r].SetNum(num)
            if r < Row {
                m.BitRows[r].shiftNum.Hide()
                m.BitRows[r].shiftNumDisplay.Show()
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
}

func NewMainForm(w *fltk.Window) {
    mainForm := new(MainForm)
    mainForm.base = 16
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
    box := NewBox(fltk.GTK_UP_BOX, WIDTH-195, pad*5, 118, 25, 12, "进制", fltk.WHITE)
    box.SetAlign(fltk.ALIGN_LEFT)
    base16 := NewRadioRoundButton(WIDTH-190, pad*7+1, 16, 16, 16, "16", mainForm.BaseChoise)
    mainForm.Base16 = base16
    base10 := NewRadioRoundButton(WIDTH-150, pad*7+1, 16, 16, 10, "10", mainForm.BaseChoise)
    mainForm.Base10 = base10
    base8 := NewRadioRoundButton(WIDTH-110, pad*7+1, 16, 16, 8, "8", mainForm.BaseChoise)
    mainForm.Base8 = base8
    addR := NewButton(WIDTH-72, pad-1, 60, 20, "增加一行")
    addR.SetCallback(mainForm.Add)
    rmR := NewButton(WIDTH-72, pad+21, 60, 20, "删除一行")
    rmR.Deactivate()
    rmR.SetCallback(mainForm.Remove)
    mainForm.AddRow = addR
    mainForm.RmRow = rmR
    ontop := NewToggleButton(pad*6, pad*4, 35, 20, "置顶")
    ontop.SetCallback(mainForm.SetOnTop)
    base16.SetValue(true)
    mainForm.ontop = ontop
    mlSwitch := NewToggleButton(pad*7+35, pad*4, 35, 20, "MSB")
    mlSwitch.SetCallback(mainForm.MLSwitch)

    rangeParse := NewToggleButton(pad*8+70, pad*4, 60, 20, "位域解析")

    bitColorBox := fltk.NewBox(fltk.GLEAM_UP_BOX, pad*12+130, pad*6, 12, 12)
    bitColorBox.SetColor(bitColorMap["1"])
    bitColorSel := NewButton(pad*9+150, pad*4, 60, 20, "颜色选择")
    colorDia := NewColorSelect(mainForm)
    callBack := func(i int) func() {
        return func() {
            colorDia.group.Show()
            colorDia.index = i
        }
    }
    bitColorSel.SetCallback(callBack(0))
    headerColorBox := fltk.NewBox(fltk.GLEAM_UP_BOX, pad*13+210, pad*6, 12, 12)
    headerColorBox.SetColor(headerColorMap[14])
    headerColorSel := NewButton(pad*10+230, pad*4, 80, 20, "对比颜色选择")
    headerColorSel.SetCallback(callBack(1))
    mainForm.BitColorSel = bitColorSel
    mainForm.BitColorBox = bitColorBox
    mainForm.HeaderColorSel = headerColorSel
    mainForm.HeaderColorBox = headerColorBox
    mainForm.MLSwitchButton = mlSwitch
    mainForm.BitRangeParse = rangeParse
    mainForm.Group = &w.Group
}

func main() {
    fltk.InitStyles()
    win := fltk.NewWindowWithPosition(StartX, StartY, WIDTH, HEIGHT, "寄存器工具")
    win.SetSizeRange(WIDTH, HEIGHT, WIDTH, maxHeight, 0, 0, false)
    win.SetColor(fltk.WHITE)
    NewMainForm(win)
    win.End()
    win.Show()
    fltk.Run()
}
