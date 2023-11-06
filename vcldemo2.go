package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

const (
	dataWidth = 32
	pad       = 3
	bdW       = 22
	bdH       = 28
	MaxRow    = 5
	winX      = int32(dataWidth*bdW) + int32(dataWidth/4)*pad*2 + 335 + int32(2.5*dataWidth)
)

var (
	Row     = 1
	MaxNum  = int64(math.Pow(2, float64(dataWidth)) - 1)
	winY    = int32(22 + (Row+1)*bdH + (Row+1)*pad)
	textMap = map[string]string{
		"0": "1",
		"1": "0",
	}
	MLmap = map[string]string{
		"MSB": "LSB",
		"LSB": "MSB",
	}
	headerColor = map[int32]types.TColor{
		10: types.TColor(0),
		12: types.TColor(0x0000ff),
	}
	headerStyle = map[int32]types.TFontStyles{
		10: types.FsBold,
		12: types.FsItalic,
	}
	bitColor = map[string]types.TColor{
		"0": types.TColor(0xffffff),
		"1": types.TColor(0xffff88),
	}
)

func NewButton(parent vcl.IWinControl, x, y, w, h int32, caption string) *vcl.TButton {
	button := vcl.NewButton(parent)
	button.SetParent(parent)
	button.SetBounds(x, y, w, h)
	button.SetCaption(caption)
	return button
}

func NewColorButton(parent vcl.IWinControl, x, y, w, h int32, caption string, fn vcl.TNotifyEvent, color types.TColor) *vcl.TColorButton {
	button := vcl.NewColorButton(parent)
	button.SetParent(parent)
	button.SetBounds(x, y, w, h)
	button.SetCaption(caption)
	button.SetOnColorChanged(fn)
	button.SetButtonColor(color)
	return button
}

func NewToggleBox(parent vcl.IWinControl, x, y, w, h int32, caption string, fn vcl.TNotifyEvent) *vcl.TToggleBox {
	button := vcl.NewToggleBox(parent)
	button.SetParent(parent)
	button.SetBounds(x, y, w, h)
	button.SetCaption(caption)
	button.SetOnClick(fn)
	return button
}

func NewRadios(parent vcl.IWinControl, fn vcl.TNotifyEvent, x, y, w, h int32) *vcl.TRadioGroup {
	checkgroup := vcl.NewRadioGroup(parent)
	checkgroup.SetParent(parent)
	checkgroup.SetCaption("进制")
	checkgroup.SetBounds(x, y, w, h)
	checkgroup.SetColumns(3)
	caps := []string{"16", "10", "8"}
	for _, cap := range(caps){
		button := vcl.NewRadioButton(checkgroup)
		button.SetParent(checkgroup)
		button.SetCaption(cap)
		if cap == "16" {
			button.SetChecked(true)
		}
		button.SetOnClick(fn)
	}
	return checkgroup
}

func GetRowIndex(sender vcl.IWinControl) int64 {
	cname := sender.Name()
	name := string(cname[len(cname)-1])
	rowIx, _ := strconv.ParseInt(name, 10, 0)
	return rowIx
}

func GetColIndex(sender vcl.IWinControl) int64 {
	cname := sender.Name()
	reg := regexp.MustCompile(`^p\d+`)
	name := reg.FindAllString(cname, -1)[0][1:]
	colIx, _ := strconv.ParseInt(name, 10, 0)
	return colIx
}

type Bit struct {
	*vcl.TMemo
}

func (b *Bit) Clicked(sender vcl.IObject) {
	var str string
	bit := vcl.AsMemo(sender)
	bit.GetTextBuf(&str, 2)
	bit.SetTextBuf(textMap[str])
	bit.SetColor(bitColor[textMap[str]])
}

func NewBit(parent vcl.IWinControl, x, y, w, h int32, s, name string) *Bit {
	bit := vcl.NewMemo(parent)
	menu := vcl.NewPopupMenu(bit)
	bit.SetParent(parent)
	bit.SetPopupMenu(menu)
	bit.SetBounds(x, y, w, h)
	bit.SetBorderStyle(types.BsSingle)
	bit.SetTextBuf(s)
	bit.Font().SetSize(12)
	bit.SetColor(bitColor["0"])
	bit.SetName(name)
	bit.SetAlignment(types.TaCenter)
	bit.SetReadOnly(true)
	return &Bit{bit}
}

type BitRow struct {
	BitLocs      []*Bit
	Num          *vcl.TMemo
	LShift       *vcl.TButton
	ShiftNum     *vcl.TMemo
	RShift       *vcl.TButton
	Reverse      *vcl.TButton
	Invert       *vcl.TButton
	Clear        *vcl.TButton
	base         int
	lastNum      int64
	lastShiftNum int64
}

func (b *BitRow) GetBitString() []string {
	bitList := make([]string, dataWidth)
	for c := 0; c < dataWidth; c++ {
		var str string
		b.BitLocs[c].GetTextBuf(&str, 2)
		bitList[c] = str
	}
	return bitList
}

func (b *BitRow) SetNum(num int64) {
	switch b.base {
	case 16:
		b.Num.SetTextBuf(fmt.Sprintf("%x", num))
	case 10:
		b.Num.SetTextBuf(fmt.Sprint(num))
	case 8:
		b.Num.SetTextBuf(fmt.Sprintf("%o", num))
	}
}

func (b *BitRow) GetNum() int64 {
	bitList := b.GetBitString()
	binStr := strings.Join(bitList, "")
	bin, _ := strconv.ParseInt(binStr, 2, dataWidth*2)
	return bin
}

func (b *BitRow) UpdateNum() {
	b.SetNum(b.GetNum())
}

func (b *BitRow) GetCurrentNum() (int64, int64) {
	var curNum string
	b.Num.GetTextBuf(&curNum, 256)
	if curNum == "" {
		b.lastNum = 0
	}
	num, err := strconv.ParseInt(curNum, b.base, dataWidth*2)
	if err == nil {
		b.lastNum = num
	} else {
		b.SetNum(b.lastNum)
	}
	var curShiftNum string
	b.ShiftNum.GetTextBuf(&curShiftNum, 64)
	shiftNum, err := strconv.ParseInt(curShiftNum, 10, dataWidth*2)
	if err == nil {
		b.lastShiftNum = shiftNum
	} else {
		b.ShiftNum.SetTextBuf(fmt.Sprint(b.lastShiftNum))
	}
	b.Num.SetSelStart(256)
	return b.lastNum, b.lastShiftNum
}

func (b *BitRow) UpdateBit(num int64) {
	binStr := strconv.FormatInt(num, 2)
	n := len(binStr)
	sum := 0
	for c := dataWidth - 1; c >= 0; c-- {
		if sum < n {
			s := string(binStr[n-sum-1])
			b.BitLocs[c].SetTextBuf(s)
			b.BitLocs[c].SetColor(bitColor[s])
		} else {
			b.BitLocs[c].SetTextBuf("0")
			b.BitLocs[c].SetColor(bitColor["0"])
		}
		sum++
	}
}

func (b *BitRow) UpdateBitNum(num int64) {
	b.SetNum(num)
	b.UpdateBit(num)
}

func (b *BitRow) SetEnable(show bool) {
	slice := []vcl.IWinControl{
		b.Num,
		b.LShift,
		b.ShiftNum,
		b.RShift,
		b.Reverse,
		b.Invert,
		b.Clear,
	}
	for _, obj := range b.BitLocs {
		tmp := obj
		go func(j *Bit) {
			vcl.ThreadSync(func() {
				if show {
					j.Show()
				} else {
					j.Hide()
				}
			})
		}(tmp)
	}
	for _, ic := range slice {
		tmp := ic
		go func(j vcl.IWinControl) {
			vcl.ThreadSync(func() {
				if show {
					j.Show()
				} else {
					j.Hide()
				}
			})
		}(tmp)
	}
}

func NewBitRow(parent vcl.IWinControl, row int, y int32) *BitRow {
	bitRow := new(BitRow)
	bitLocs := make([]*Bit, dataWidth)
	for c := 0; c < dataWidth; c++ {
		n := int32(c*bdW) + int32(c/4)*pad*2 + int32(c+1)*pad
		if n == 1 {
			n = 4
		}
		bit := NewBit(parent, n, y, bdW, bdH, "0", fmt.Sprintf("p%dhead%d", c, row))
		bitLocs[c] = bit

	}
	bitsWidth := int32(dataWidth*bdW) + int32(dataWidth/4)*pad*2 + int32(dataWidth+1)*pad
	bitRow.BitLocs = bitLocs
	num := vcl.NewMemo(parent)
	num.SetParent(parent)
	num.SetBounds(bitsWidth, y, int32(2.5*dataWidth)+40, bdH)
	num.SetTextBuf("0")
	num.Font().SetSize(12)
	num.SetName(fmt.Sprintf("numEdit%d", row))
	bitRow.Num = num
	lShift := NewButton(parent, bitsWidth+42+int32(2.5*dataWidth), y, 30, bdH, "<<")
	lShift.SetName(fmt.Sprintf("lshift%d", row))
	bitRow.LShift = lShift
	shiftNum := vcl.NewMemo(parent)
	shiftNum.SetParent(parent)
	shiftNum.SetBounds(bitsWidth+74+int32(2.5*dataWidth), y, 30, bdH)
	shiftNum.SetTextBuf("1")
	shiftNum.SetAlignment(types.TaCenter)
	shiftNum.Font().SetSize(12)
	bitRow.ShiftNum = shiftNum
	rShift := NewButton(parent, bitsWidth+106+int32(2.5*dataWidth), y, 30, bdH, ">>")
	rShift.SetName(fmt.Sprintf("rshift%d", row))
	bitRow.RShift = rShift
	reverse := NewButton(parent, bitsWidth+138+int32(2.5*dataWidth), y, 30, bdH, "倒序")
	reverse.SetName(fmt.Sprintf("reverse%d", row))
	bitRow.Reverse = reverse
	invert := NewButton(parent, bitsWidth+170+int32(2.5*dataWidth), y, 30, bdH, "转换")
	invert.SetName(fmt.Sprintf("invert%d", row))
	bitRow.Invert = invert
	clear := NewButton(parent, bitsWidth+202+int32(2.5*dataWidth), y, 30, bdH, "清空")
	clear.SetName(fmt.Sprintf("clear%d", row))
	bitRow.Clear = clear
	bitRow.base = 16
	bitRow.lastNum = 0
	bitRow.lastShiftNum = 1
	return bitRow
}

type Header struct {
	label *vcl.TLabel
	frame *vcl.TFrame
}

func (h *Header) ChangeLabelColor(size int32) {
	h.label.Font().SetSize(size)
	h.label.Font().SetStyle(headerStyle[size])
	h.label.Font().SetColor(headerColor[size])
}

func NewHeader(parent vcl.IWinControl, x, y, w, h int32, s string) *Header {
	frame := vcl.NewFrame(parent)
	frame.SetParent(parent)
	frame.SetBounds(x, y, w, h)
	label := vcl.NewLabel(frame)
	label.SetParent(frame)
	label.AnchorHorizontalCenterTo(frame)
	label.AnchorVerticalCenterTo(frame)
	label.SetCaption(s)
	label.Font().SetSize(10)
	return &Header{label: label, frame: frame}
}

type Headers []*Header

func (h Headers) UpdateHeader(bitMap map[string]int, c int) {
	if len(bitMap) == 1 {
		h[c].ChangeLabelColor(10)
	} else {
		h[c].ChangeLabelColor(12)
	}
}

func NewHeaders(parent vcl.IWinControl, y int32) Headers {
	headers := make([]*Header, dataWidth)
	for c := 0; c < dataWidth; c++ {
		n := int32(c*bdW) + int32(c/4)*pad*2 + int32(c+1)*pad
		if n == 1 {
			n = 4
		}
		head := NewHeader(parent, n, y, bdW, bdW, fmt.Sprint(dataWidth-1-c))
		headers[c] = head
	}
	return headers
}

type BitAnalyze struct {
	frame  *vcl.TFrame
	labels []*vcl.TLabel
	res    []*vcl.TEdit
}

func NewBitAnalyze(parent vcl.IWinControl) *BitAnalyze {
	width := (winX-pad*2)/(MaxRow+1) - 3
	bitAnalyze := new(BitAnalyze)
	frame := vcl.NewFrame(parent)
	frame.SetParent(parent)
	frame.SetWidth(winX - pad)
	frame.SetHeight(bdH*2 - pad)
	res := make([]*vcl.TEdit, MaxRow+1)
	labels := make([]*vcl.TLabel, MaxRow+1)
	for c := 0; c <= MaxRow; c++ {
		memo := vcl.NewEdit(frame)
		memo.SetParent(frame)
		memo.SetBounds(int32(c)*(width+pad)+3, bdH-pad*2-2, width, bdH)
		memo.SetMaxLength(256)
		memo.Font().SetSize(12)
		label := vcl.NewLabel(frame)
		label.SetParent(frame)
		label.SetBounds(int32(c)*(width+pad)+3, pad, width, bdH)
		var title string
		if c == 0 {
			title = "\tBit位域"
			memo.SetHint("e.g. 多个域段 31: 0\n      单个域段 31")
			memo.SetShowHint(true)
		} else {
			title = fmt.Sprintf("\t第%d行", c)
			memo.SetReadOnly(true)
		}
		if c > Row {
			memo.Hide()
			label.Hide()
		}
		label.SetCaption(title)
		labels[c] = label
		res[c] = memo
	}
	frame.Hide()
	bitAnalyze.frame = frame
	bitAnalyze.labels = labels
	bitAnalyze.res = res
	return bitAnalyze
}

type TMainForm struct {
	*vcl.TForm
	Headers            Headers
	BitRows            []*BitRow
	BaseChoise         *vcl.TRadioGroup
	base               int
	AddRow             *vcl.TButton
	RmRow              *vcl.TButton
	HeaderSwitch       *vcl.TToggleBox
	OnTop              *vcl.TCheckBox
	ColorSetting       *vcl.TColorButton
	HeaderColorSetting *vcl.TColorButton
	AnalyzeButton      *vcl.TToggleBox
	AnalyzeArea        *BitAnalyze
}

func (f *TMainForm) initComponents(cols, rows int) {
	f.base = 16
	addrow := NewButton(f, winX-pad-70, pad+5, 60, 18, "增加一行")
	addrow.SetOnClick(f.AddR)
	rmrow := NewButton(f, winX-pad-70, pad+25, 60, 18, "删除一行")
	rmrow.SetEnabled(false)
	rmrow.SetOnClick(f.RemoveR)
	checkgroup := NewRadios(f, f.BaseChange, winX-195, pad, 120, 40)
	f.BaseChoise = checkgroup
	cb := vcl.NewCheckBox(f)
	cb.SetParent(f)
	cb.SetCaption("置顶")
	cb.SetBounds(winX-243, 16, 10, 10)
	cb.SetOnClick(f.ClickOnTop)
	f.OnTop = cb
	hSwitch := NewToggleBox(f, 16, pad*2, 60, 22, "MSB", f.MLSwitch)
	f.HeaderSwitch = hSwitch
	bits := make([]*BitRow, MaxRow)
	f.Headers = NewHeaders(f, 32)
	for r := 0; r < MaxRow; r++ {
		bits[r] = NewBitRow(f, r, int32(22+(r+1)*bdH+(r+1)*pad))
		for c := 0; c < dataWidth; c++ {
			bits[r].BitLocs[c].SetOnMouseDown(f.Clicked)
		}
		bits[r].Num.SetOnKeyUp(f.KeyTyped)
		bits[r].LShift.SetOnClick(f.ClickLShift)
		bits[r].RShift.SetOnClick(f.ClickRShift)
		bits[r].Reverse.SetOnClick(f.ClickReverse)
		bits[r].Invert.SetOnClick(f.ClickInvert)
		bits[r].Clear.SetOnClick(f.ClickClear)
		if r > Row-1 {
			bits[r].SetEnable(false)
		}
	}
	f.BitRows = bits
	f.BaseChoise = checkgroup
	f.AddRow = addrow
	f.RmRow = rmrow
	colorSetting := NewColorButton(f, 80, pad*2, 75, 22, "颜色选择", f.SelectColor, bitColor["1"])
	headColorSetting := NewColorButton(f, 158, pad*2, 99, 22, "对比颜色选择", f.SelectHeaderColor, headerColor[12])
	f.ColorSetting = colorSetting
	f.HeaderColorSetting = headColorSetting
	analyzeArea := NewBitAnalyze(f)
	analyzeArea.res[0].SetOnKeyUp(f.Edit)
	f.AnalyzeArea = analyzeArea
	analyzeButton := NewToggleBox(f, winX/2-35, pad*2, 70, 22, "位域解析", f.Analyze)
	f.AnalyzeButton = analyzeButton
}

func (f *TMainForm) UpdateHeader(bitMap map[string]int, c int) {
	if len(bitMap) == 1 {
		f.Headers[c].ChangeLabelColor(10)
	} else {
		f.Headers[c].ChangeLabelColor(12)
	}
}

func (f *TMainForm) UpdateHeaders() {
	for c := 0; c < dataWidth; c++ {
		bitMap := make(map[string]int, Row)
		for r := 0; r < Row; r++ {
			var str string
			f.BitRows[r].BitLocs[c].GetTextBuf(&str, 2)
			bitMap[str] = 0
		}
		f.UpdateHeader(bitMap, c)
	}
}

func (f *TMainForm) KeyTyped(sender vcl.IObject, key *types.Char, shift types.TShiftState) {
	numEdit := vcl.AsMemo(sender)
	rowIx := GetRowIndex(numEdit)
	num, _ := f.BitRows[rowIx].GetCurrentNum()
	f.BitRows[rowIx].UpdateBit(num)
	f.UpdateHeaders()
	f.Edit(f.AnalyzeArea.res[0], key, shift)
}

func (f *TMainForm) Clicked(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	bit := vcl.AsMemo(sender)
	colIx := GetColIndex(bit)
	rowIx := GetRowIndex(bit)
	f.BitRows[rowIx].Num.SetFocus()
	f.BitRows[rowIx].BitLocs[colIx].Clicked(bit)
	f.BitRows[rowIx].UpdateNum()
	f.UpdateHeaders()
	f.BitRows[rowIx].Num.SetSelStart(256)
	var key uint16 = 0xff
	f.Edit(f.AnalyzeArea.res[0], &key, shift)
}

func (f *TMainForm) UpdateAnalyzeArea() {
	var key uint16 = 0xff
	f.Edit(f.AnalyzeArea.res[0], &key, 0)
}

func (f *TMainForm) ClickClear(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	rowIx := GetRowIndex(button)
	f.BitRows[rowIx].UpdateBitNum(0)
	f.UpdateHeaders()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) ClickInvert(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	rowIx := GetRowIndex(button)
	for c := 0; c < dataWidth; c++ {
		f.BitRows[rowIx].BitLocs[c].Clicked(f.BitRows[rowIx].BitLocs[c])
	}
	f.BitRows[rowIx].UpdateNum()
	f.UpdateHeaders()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) ClickLShift(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	rowIx := GetRowIndex(button)
	num, shiftNum := f.BitRows[rowIx].GetCurrentNum()
	num = (num << shiftNum) & MaxNum
	f.BitRows[rowIx].UpdateBitNum(num)
	f.UpdateHeaders()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) ClickRShift(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	rowIx := GetRowIndex(button)
	num, shiftNum := f.BitRows[rowIx].GetCurrentNum()
	num = (num >> shiftNum) & MaxNum
	f.BitRows[rowIx].UpdateBitNum(num)
	f.UpdateHeaders()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) ClickReverse(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	rowIx := GetRowIndex(button)
	var left, right string
	for i, j := 0, len(f.BitRows[rowIx].BitLocs)-1; i < j; i, j = i+1, j-1 {
		f.BitRows[rowIx].BitLocs[i].GetTextBuf(&left, 2)
		f.BitRows[rowIx].BitLocs[j].GetTextBuf(&right, 2)
		f.BitRows[rowIx].BitLocs[i].SetTextBuf(right)
		f.BitRows[rowIx].BitLocs[i].SetColor(bitColor[right])
		f.BitRows[rowIx].BitLocs[j].SetTextBuf(left)
		f.BitRows[rowIx].BitLocs[j].SetColor(bitColor[left])
	}
	f.BitRows[rowIx].UpdateNum()
	f.UpdateHeaders()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) BaseChange(sender vcl.IObject) {
	var str string
	oldbase := f.base
	ra := vcl.AsRadioButton(sender)
	ra.GetTextBuf(&str, 4)
	base, _ := strconv.ParseInt(str, 10, 16)
	f.base = int(base)
	for i := 0; i < Row; i++ {
		var bitString string
		f.BitRows[i].Num.GetTextBuf(&bitString, dataWidth*2)
		f.BitRows[i].base = int(base)
		num, _ := strconv.ParseInt(bitString, oldbase, dataWidth*2)
		f.BitRows[i].UpdateBitNum(num)
	}
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) AddR(sender vcl.IObject) {
	Row++
	f.BitRows[Row-1].base = f.base
	f.RmRow.SetEnabled(true)
	if Row == MaxRow {
		f.AddRow.SetEnabled(false)
	}
	winY += bdH + pad
	f.SetHeight(winY)
	if f.AnalyzeButton.Checked() {
		f.AnalyzeArea.frame.SetLeft(pad)
		f.AnalyzeArea.frame.SetTop(int32(22 + (Row+1)*bdH + (Row+1)*pad))
	}
	f.BitRows[Row-1].SetEnable(true)
	f.BitRows[Row-1].UpdateNum()
	f.UpdateHeaders()
	f.AnalyzeArea.res[Row].Show()
	f.AnalyzeArea.labels[Row].Show()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) RemoveR(sender vcl.IObject) {
	Row--
	f.AddRow.SetEnabled(true)
	if Row == 1 {
		f.RmRow.SetEnabled(false)
	}
	winY -= bdH + pad
	f.SetHeight(winY)
	if f.AnalyzeButton.Checked() {
		f.AnalyzeArea.frame.SetLeft(pad)
		f.AnalyzeArea.frame.SetTop(int32(22 + (Row+1)*bdH + (Row+1)*pad))
	}
	f.BitRows[Row].Clear.Click()
	f.BitRows[Row].SetEnable(false)
	f.UpdateHeaders()
	f.AnalyzeArea.res[Row+1].Hide()
	f.AnalyzeArea.labels[Row+1].Hide()
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) ClickOnTop(sender vcl.IObject) {
	cb := vcl.AsCheckBox(sender)
	if cb.Checked() {
		f.SetFormStyle(types.FsSystemStayOnTop)
	} else {
		f.SetFormStyle(types.FsNormal)
	}
}

func (f *TMainForm) MLSwitch(sender vcl.IObject) {
	button := vcl.AsButton(sender)
	val := button.Caption()
	button.SetCaption(MLmap[val])
	for c := 0; c < dataWidth; c++ {
		var label string
		if val == "LSB" {
			label = fmt.Sprint(dataWidth - 1 - c)
		} else {
			label = fmt.Sprint(c)
		}
		f.Headers[c].label.SetCaption(label)
	}
	f.UpdateAnalyzeArea()
}

func (f *TMainForm) SelectColor(sender vcl.IObject) {
	colorButtom := vcl.AsColorButton(sender)
	bitColor["1"] = colorButtom.ButtonColor()
	for r := 0; r < MaxRow; r++ {
		num, _ := f.BitRows[r].GetCurrentNum()
		f.BitRows[r].UpdateBitNum(num)
	}
}

func (f *TMainForm) SelectHeaderColor(sender vcl.IObject) {
	colorButtom := vcl.AsColorButton(sender)
	headerColor[12] = colorButtom.ButtonColor()
	f.UpdateHeaders()
}

func (f *TMainForm) Analyze(sender vcl.IObject) {
	cb := vcl.AsCheckBox(sender)
	if cb.Checked() {
		f.AnalyzeArea.frame.SetLeft(pad)
		f.AnalyzeArea.frame.SetTop(int32(22 + (Row+1)*bdH + (Row+1)*pad))
		f.AnalyzeArea.frame.Show()
		winY += bdH*2 - pad
	} else {
		f.AnalyzeArea.frame.Hide()
		winY -= bdH*2 - pad
	}
	f.SetClientHeight(winY)
}

func (f *TMainForm) SetNum(num int64, r int32) {
	for c := 0; c <= MaxRow; c++ {
		switch f.base {
		case 16:
			f.AnalyzeArea.res[r+1].SetTextBuf(fmt.Sprintf("%x", num))
		case 10:
			f.AnalyzeArea.res[r+1].SetTextBuf(fmt.Sprint(num))
		case 8:
			f.AnalyzeArea.res[r+1].SetTextBuf(fmt.Sprintf("%o", num))
		}
	}
}

func (f *TMainForm) ParseBitRange(nums []string, r int32) (int64, error) {
	var res []string
	if len(nums) == 1 {
		left := strings.Trim(nums[0], "\r\n")
		num, err := strconv.ParseInt(left, 10, 0)
		if err != nil || num >= dataWidth {
			return 0, fmt.Errorf("无效输入")
		}
		for c := 0; c < dataWidth; c++ {
			var str string
			if f.Headers[c].label.Caption() == left {
				f.BitRows[r].BitLocs[c].GetTextBuf(&str, 2)
				res = append(res, str)
			}
		}
		bin, _ := strconv.ParseInt(strings.Join(res, ""), 2, dataWidth*2)
		return bin, nil
	} else if len(nums) == 2 {
		left := strings.Trim(nums[0], "\r\n")
		right := strings.Trim(nums[1], "\r\n")
		numL, errL := strconv.ParseInt(left, 10, 0)
		if errL != nil || numL >= dataWidth {
			return 0, fmt.Errorf("无效输入")
		}
		numR, errR := strconv.ParseInt(right, 10, 0)
		if errR != nil || numR >= dataWidth {
			return 0, fmt.Errorf("无效输入")
		}
		if numR == numL {
			return 0, fmt.Errorf("无效输入")
		}
		flag := false
		for c := 0; c < dataWidth; c++ {
			var str string
			cap := f.Headers[c].label.Caption()
			if (cap == left && !flag) || (cap == right && !flag) {
				flag = true
			} else if (cap == right && flag) || (cap == left && flag) {
				flag = false
				f.BitRows[r].BitLocs[c].GetTextBuf(&str, 2)
				res = append(res, str)
			}
			if flag {
				f.BitRows[r].BitLocs[c].GetTextBuf(&str, 2)
				res = append(res, str)
			}
		}
		bin, _ := strconv.ParseInt(strings.Join(res, ""), 2, dataWidth*2)
		return bin, nil
	} else {
		return 0, fmt.Errorf("无效输入")
	}
}

func (f *TMainForm) Edit(sender vcl.IObject, key *types.Char, shift types.TShiftState) {
	edit := vcl.AsMemo(sender)
	var str string
	edit.GetTextBuf(&str, 256)
	for r := 0; r < MaxRow; r++ {
		num, err := f.ParseBitRange(strings.Split(str, ":"), int32(r))
		if err != nil {
			if str != "" {
				f.AnalyzeArea.res[r+1].SetTextBuf("无效输入")
				f.AnalyzeArea.res[r+1].SetColor(types.TColor(0x0000ff))
			} else {
				f.AnalyzeArea.res[r+1].SetTextBuf("")
				f.AnalyzeArea.res[r+1].SetColor(types.TColor(0xffffff))
			}
		} else {
			f.SetNum(num, int32(r))
			f.AnalyzeArea.res[r+1].SetColor(types.TColor(0xffffff))
		}
	}
}

var mainForm *TMainForm

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&mainForm)
	mainForm.EnabledMaximize(false)
	mainForm.WorkAreaCenter()
	vcl.Application.Run()
}

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("寄存器工具")
	f.SetClientHeight(winY)
	f.SetClientWidth(winX)
	f.SetColor(bitColor["0"])
	f.initComponents(dataWidth, Row)
}
