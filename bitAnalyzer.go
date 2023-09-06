package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

const (
	padx      = 4
	pady      = 4
	bitBgX    = 20
	bitBgY    = 25
	bitWidth  = 32
	bitNumEdX = int32(2.5*bitWidth) + 10
	ButtonS   = 30
	winX      = (bitWidth+1)*bitBgX + 2*padx + bitNumEdX + 5*ButtonS
)

var (
	color = map[string]types.TColor{
		"0":    types.TColor(0xffffff),
		"1":    types.TColor(0xffff88),
		"diff": types.TColor(0xffaaff),
		"same": types.TColor(0xf0f0f0),
	}
	Row  = 1
	MaxRow = 3
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
)

type bit interface {
	GetTextBuf(Buffer *string, BufSize int32) int32
	SetTextBuf(Buffer string)
	SetColor(value types.TColor)
	As() vcl.TAs
	Is() vcl.TIs
	ClassName() string
	ClassType() types.TClass
	Equals(Obj vcl.IObject) bool
	Free()
	GetHashCode() int32
	InheritsFrom(AClass types.TClass) bool
	Instance() uintptr
	InstanceSize() int32
	IsValid() bool
	ToString() string
	UnsafeAddr() unsafe.Pointer
	Show()
	Hide()
}

type BitLoc []bit

func newMemo(parent vcl.IWinControl, x, y, w, h int32, ix, row, bitWidth int, color types.TColor, text string, show bool, fn ...vcl.TNotifyEvent) *vcl.TMemo {
	memo := vcl.NewMemo(parent)
	menu := vcl.NewPopupMenu(parent)
	memo.SetParent(parent)
	memo.SetPopupMenu(menu)
	memo.SetTextBuf(text)
	memo.SetBounds(x+int32(ix%bitWidth)*w, y+int32(row)*h, w, h)
	var maxLength int32
	if row < 1 {
		maxLength = 3
		memo.SetBorderStyle(types.BsNone)
		memo.SetHeight(17)
		memo.SetTop(y+7)
		memo.SetControlState(types.CsNoStdEvents)
		memo.SetName(fmt.Sprintf("m%dhead%d", ix, 0))
	} else {
		maxLength = 1
		memo.SetOnClick(fn[0])
		memo.SetName(fmt.Sprintf("m%dbit%d", ix, row-1))
	}
	memo.SetMaxLength(maxLength)
	memo.SetColor(color)
	memo.SetAlignment(types.TaCenter)
	memo.SetReadOnly(true)
	if show {
		memo.Show()
	} else {
		memo.Hide()
	}
	return memo
}

func newBitLoc(parent vcl.IWinControl, x, y, w, h int32, bitWidth, row int, color types.TColor, show bool,fnc vcl.TKeyEvent, fn ...vcl.TNotifyEvent) BitLoc {
	bit := make(BitLoc, bitWidth+7)
	for c := 0; c < bitWidth; c++ {
		bit[c] = newMemo(parent, x, y, w, h, c, row, bitWidth, color, "0", show, fn[0])	
	}
	numEdit := vcl.NewMemo(parent)
	numEdit.SetParent(parent)
	numEdit.SetBounds(int32(padx+bitWidth*bitBgX), padx+int32(row)*bitBgY+50, bitNumEdX, bitBgY)
	numEdit.SetOnKeyUp(fnc)
	numEdit.SetName(fmt.Sprintf("numEdit%d", row-1))
	numEdit.SetTextBuf("0")
	lshift := vcl.NewButton(parent)
	lshift.SetParent(parent)
	lshift.SetBounds(int32(padx+bitWidth*bitBgX)+bitNumEdX, padx+int32(row)*bitBgY+50, ButtonS, bitBgY)
	lshift.SetTextBuf("<<")
	lshift.SetOnClick(fn[1])
	lshift.SetName(fmt.Sprintf("lshift%d", row-1))
	shiftnum := vcl.NewEdit(parent)
	shiftnum.SetParent(parent)
	shiftnum.SetBounds(int32(padx+bitWidth*bitBgX)+bitNumEdX+ButtonS, padx+int32(row)*bitBgY+50, bitBgX, bitBgY)
	shiftnum.SetTextBuf("1")
	shiftnum.SetMaxLength(2)
	shiftnum.SetAlignment(types.TaCenter)
	rshift := vcl.NewButton(parent)
	rshift.SetParent(parent)
	rshift.SetBounds(int32(padx+(bitWidth+1)*bitBgX)+bitNumEdX+ButtonS, padx+int32(row)*bitBgY+50, ButtonS, bitBgY)
	rshift.SetTextBuf(">>")
	rshift.SetOnClick(fn[1])
	rshift.SetName(fmt.Sprintf("rshift%d", row-1))
	rev := vcl.NewButton(parent)
	rev.SetParent(parent)
	rev.SetBounds(int32(padx+(bitWidth+1)*bitBgX)+bitNumEdX+ButtonS*2, padx+int32(row)*bitBgY+50, ButtonS, bitBgY)
	rev.SetTextBuf("倒序")
	rev.SetOnClick(fn[2])
	rev.SetName(fmt.Sprintf("rev%d", row-1))
	invt := vcl.NewButton(parent)
	invt.SetParent(parent)
	invt.SetBounds(int32(padx+(bitWidth+1)*bitBgX)+bitNumEdX+ButtonS*3, padx+int32(row)*bitBgY+50, ButtonS, bitBgY)
	invt.SetTextBuf("转换")
	invt.SetOnClick(fn[3])
	invt.SetName(fmt.Sprintf("invt%d", row-1))
	cler := vcl.NewButton(parent)
	cler.SetParent(parent)
	cler.SetBounds(int32(padx+(bitWidth+1)*bitBgX)+bitNumEdX+ButtonS*4, padx+int32(row)*bitBgY+50, ButtonS, bitBgY)
	cler.SetTextBuf("清空")
	cler.SetOnClick(fn[4])
	cler.SetName(fmt.Sprintf("cler%d", row-1))
	if show {
		numEdit.Show()
		lshift.Show()
		shiftnum.Show()
		rshift.Show()
		rev.Show()
		invt.Show()
		cler.Show()
	} else {
		numEdit.Hide()
		lshift.Hide()
		shiftnum.Hide()
		rshift.Hide()
		rev.Hide()
		invt.Hide()
		cler.Hide()
	}
	bit[bitWidth] = numEdit
	bit[bitWidth+1] = lshift
	bit[bitWidth+2] = shiftnum
	bit[bitWidth+3] = rshift
	bit[bitWidth+4] = rev
	bit[bitWidth+5] = invt
	bit[bitWidth+6] = cler
	return bit
}

type TMainForm struct {
	*vcl.TForm
	BitLocs    []BitLoc
	BitHeader  []*vcl.TMemo
	BaseChoise *vcl.TRadioGroup
	base       int
	AddRow     *vcl.TButton
	RmRow      *vcl.TButton
	OnTop      *vcl.TCheckBox
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
	f.initComponents(f, bitWidth, Row, color)
}

func (f *TMainForm) initComponents(parent vcl.IWinControl, cols, rows int, color map[string]types.TColor) {
	f.base = 16
	addrow := vcl.NewButton(parent)
	addrow.SetParent(parent)
	addrow.SetBounds(winX-padx-ButtonS*2, pady+5, ButtonS*2, 20)
	addrow.SetTextBuf("增加一行")
	addrow.SetOnClick(f.AddR)
	rmrow := vcl.NewButton(parent)
	rmrow.SetParent(parent)
	rmrow.SetBounds(winX-padx-ButtonS*2, pady+25, ButtonS*2, 20)
	rmrow.SetEnabled(false)
	rmrow.SetTextBuf("删除一行")
	rmrow.SetOnClick(f.RemoveR)
	checkgroup := vcl.NewRadioGroup(parent)
	checkgroup.SetParent(parent)
	checkgroup.SetCaption("进制")
	checkgroup.SetBounds(winX-padx*2-120-ButtonS*2, pady, 120, 46)
	checkgroup.SetColumns(3)
	checkbutton16 := vcl.NewRadioButton(checkgroup)
	checkbutton16.SetParent(checkgroup)
	checkbutton16.SetCaption("16")
	checkbutton16.SetChecked(true)
	checkbutton16.SetOnClick(f.BaseChange)
	checkbutton10 := vcl.NewRadioButton(checkgroup)
	checkbutton10.SetParent(checkgroup)
	checkbutton10.SetCaption("10")
	checkbutton10.SetOnClick(f.BaseChange)
	checkbutton8 := vcl.NewRadioButton(checkgroup)
	checkbutton8.SetParent(checkgroup)
	checkbutton8.SetCaption("8")
	checkbutton8.SetOnClick(f.BaseChange)
	cb := vcl.NewCheckBox(parent)
	cb.SetParent(parent)
	cb.SetCaption("置顶")
	cb.SetBounds(winX-padx*5-160-ButtonS*2, 20, 10, 10)
	cb.SetOnClick(f.ClickOnTop)
	f.OnTop = cb
	headers := make([]*vcl.TMemo, cols)
	bits := make([]BitLoc, MaxRow)
	for r := 0; r <= MaxRow; r++ {
		show := r <= 1
		if r == 0 {
			for col := 0; col < cols; col++ {
				headers[col] = newMemo(parent, padx, pady+50, bitBgX, bitBgY, col, r, bitWidth, color["same"], fmt.Sprint(bitWidth-col-1), show)
			}
		} else {
			bits[r-1] = newBitLoc(parent, padx, pady+50, bitBgX, bitBgY, bitWidth, r, color["0"], show, f.Typed, f.Clicked, f.ClickShift, f.ClickReverse, f.ClickInvert, f.ClickClear)
		}
	}
	f.BitLocs = bits
	f.BaseChoise = checkgroup
	f.BitHeader = headers
	f.AddRow = addrow
	f.RmRow = rmrow
}

func (f *TMainForm) Typed(sender vcl.IObject, key *types.Char, shift types.TShiftState) {
	var str string
	num := vcl.AsMemo(sender)
	num.GetTextBuf(&str, bitWidth)
	rowIx := f.GetRowIndex(num)
	resNum, err := strconv.ParseInt(str, f.base, bitWidth*2)
	if err != nil && str != "" {
		bitList := f.GetBitString(int(rowIx))
		binStr := strings.Join(bitList, "")
		bin, _ := strconv.ParseInt(binStr, 2, bitWidth*2)
		f.UpdateBitNum(bin, rowIx)
		num.SetSelStart(int32(len(str)))
		var bitString string
		f.BitLocs[rowIx][bitWidth].GetTextBuf(&bitString, bitWidth*2)
		resNum, _ = strconv.ParseInt(bitString, f.base, bitWidth*2)
	}
	resNum &= 0xffffffff
	f.UpdateBit(rowIx, resNum)
}

func (f *TMainForm) Clicked(sender vcl.IObject) {
	var str string
	bitMap := make(map[string]int, Row)
	bit := vcl.AsMemo(sender)
	bit.GetTextBuf(&str, 2)
	bit.SetMaxLength(2)
	if str == "1" {
		bit.SetTextBuf("0")
		bit.SetColor(color["0"])
	} else {
		bit.SetTextBuf("1")
		bit.SetColor(color["1"])
	}
	bit.SetAlignment(types.TaCenter)
	rowIx := f.GetRowIndex(bit)
	colIx := f.GetColIndex(bit)
	for i := 0; i < Row; i++ {
		var bitString string
		f.BitLocs[i][colIx].GetTextBuf(&bitString, 2)
		bitMap[bitString] = 0
	}
	f.UpdateHeader(bitMap, int(colIx))
	bitList := f.GetBitString(int(rowIx))
	binStr := strings.Join(bitList, "")
	bin, _ := strconv.ParseInt(binStr, 2, bitWidth*2)
	f.UpdateBitNum(bin, int64(rowIx))
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
		f.BitLocs[i][bitWidth].GetTextBuf(&bitString, bitWidth)
		num, _ := strconv.ParseInt(bitString, oldbase, bitWidth*2)
		f.UpdateBitNum(num, int64(i))
	}
}

func (f *TMainForm) ClickClear(sender vcl.IObject) {
	cler := vcl.AsButton(sender)
	rowIx := f.GetRowIndex(cler)
	f.BitLocs[rowIx][bitWidth].SetTextBuf("0")
	for c := 0; c < bitWidth; c++ {
		bitMap := make(map[string]int, Row)
		for i := 0; i < MaxRow; i++ {
			var bitString string
			if int64(i) == rowIx {
				f.BitLocs[rowIx][c].SetTextBuf("0")
				f.BitLocs[rowIx][c].SetColor(color["0"])
			}
			f.BitLocs[i][c].GetTextBuf(&bitString, 2)
			bitMap[bitString] = 0
		}
		f.UpdateHeader(bitMap, c)
	}
}

func (f *TMainForm) ClickInvert(sender vcl.IObject) {
	inv := vcl.AsButton(sender)
	rowIx := f.GetRowIndex(inv)
	for i := 0; i < bitWidth; i++ {
		f.Clicked(f.BitLocs[rowIx][i])
	}
}

func (f *TMainForm) ClickShift(sender vcl.IObject) {
	shift := vcl.AsButton(sender)
	cname := shift.Name()
	rowIx := f.GetRowIndex(shift)
	var col int
	if cname[0] == 'l' {
		col = 0
	} else {
		col = 1
	}
	var str string
	f.BitLocs[rowIx][bitWidth+2].GetTextBuf(&str, 8)
	shiftNum, _ := strconv.ParseInt(str, 10, 16)
	f.BitLocs[rowIx][bitWidth].GetTextBuf(&str, bitWidth*32)
	num, _ := strconv.ParseInt(str, f.base, bitWidth*2)
	switch col {
	case 0:
		num <<= shiftNum
	case 1:
		num >>= shiftNum
	}
	num &= 0xffffffff
	f.UpdateBitNum(num, rowIx)
	f.UpdateBit(rowIx, num)
}

func (f *TMainForm) ClickReverse(sender vcl.IObject) {
	rev := vcl.AsButton(sender)
	rowIx := f.GetRowIndex(rev)
	bins := make([]string, bitWidth)
	for c := 0; c < bitWidth; c++ {
		bitMap := make(map[string]int, Row)
		for i := 0; i < Row; i++ {
			var str string
			if int64(i) == rowIx {
				f.BitLocs[i][c].GetTextBuf(&str, 2)
				bins[bitWidth-c-1] = str
			} else {
				f.BitLocs[i][bitWidth-c-1].GetTextBuf(&str, 2)
			}
			bitMap[str] = 0
		}
		f.UpdateHeader(bitMap, bitWidth-c-1)
	}
	binStr := strings.Join(bins, "")
	bin, _ := strconv.ParseInt(binStr, 2, bitWidth*2)
	f.UpdateBitNum(bin, rowIx)
	for i, bin := range bins {
		f.BitLocs[rowIx][i].SetTextBuf(bin)
		f.BitLocs[rowIx][i].SetColor(color[bin])
	}
}

func (f *TMainForm) AddR(sender vcl.IObject) {
	Row++
	f.RmRow.SetEnabled(true)
	if Row == MaxRow {
		f.AddRow.SetEnabled(false)
	}
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
	f.SetHeight(winY)
	f.Repaint()
	for _, obj := range(f.BitLocs[Row-1]) {
		obj.Show()
	}
	f.UpdateHeaders()
}

func (f *TMainForm) RemoveR(sender vcl.IObject) {
	Row--
	f.AddRow.SetEnabled(true)
	if Row == 1 {
		f.RmRow.SetEnabled(false)
	}
	bits := f.BitLocs[Row]
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
	f.SetHeight(winY)
	f.ClickClear(bits[bitWidth+6])
	f.UpdateHeaders()
	for _, obj := range bits {
		obj.Hide()
	}
}

func (f *TMainForm) GetRowIndex(sender vcl.IWinControl) int64 {
	cname := sender.Name()
	name := string(cname[len(cname)-1])
	rowIx, _ := strconv.ParseInt(name, 10, 0)
	return rowIx
}

func (f *TMainForm) GetColIndex(sender vcl.IWinControl) int64 {
	cname := sender.Name()
	reg := regexp.MustCompile(`^m\d+`)
	name := reg.FindAllString(cname, -1)[0][1:]
	colIx, _ := strconv.ParseInt(name, 10, 0)
	return colIx
}

func (f *TMainForm) UpdateBitNum(bin, r int64) {
	switch f.base {
	case 16:
		f.BitLocs[r][bitWidth].SetTextBuf(fmt.Sprintf("%x", bin))
	case 10:
		f.BitLocs[r][bitWidth].SetTextBuf(fmt.Sprint(bin))
	case 8:
		f.BitLocs[r][bitWidth].SetTextBuf(fmt.Sprintf("%o", bin))
	}
}

func (f *TMainForm) UpdateBit(row, resNum int64) {
	binStr := strconv.FormatInt(resNum, 2)
	n := len(binStr)
	sum := 0
	for c := bitWidth - 1; c >= 0; c-- {
		bitMap := make(map[string]int, Row)
		for r := 0; r < Row; r++ {
			var binString string
			f.BitLocs[r][c].GetTextBuf(&binString, 2)
			if int64(r) == row {
				if sum < n {
					s := string(binStr[n-sum-1])
					f.BitLocs[r][c].SetTextBuf(s)
					f.BitLocs[r][c].SetColor(color[s])
					binString = s
				} else {
					f.BitLocs[r][c].SetTextBuf("0")
					f.BitLocs[r][c].SetColor(color["0"])
					binString = "0"
				}
			}
			bitMap[binString] = 0
		}
		f.UpdateHeader(bitMap, c)
		sum++
	}
}

func (f *TMainForm) UpdateHeader(bitMap map[string]int, c int) {
	if len(bitMap) == 1 {
		f.BitHeader[c].SetColor(color["same"])
	} else {
		f.BitHeader[c].SetColor(color["diff"])
	}
}

func (f *TMainForm) UpdateHeaders() {
	for c := 0; c < bitWidth; c++ {
		bitMap := make(map[string]int, Row)
		for r := 0; r < Row; r++ {
			var str string
			f.BitLocs[r][c].GetTextBuf(&str, 2)
			bitMap[str] = 0
		}
		f.UpdateHeader(bitMap, c)
	}
}

func (f *TMainForm) GetBitString(r int) []string {
	bitList := make([]string, bitWidth)
	for i := 0; i < bitWidth; i++ {
		var bitString string
		f.BitLocs[r][i].GetTextBuf(&bitString, 2)
		bitList[i] = bitString
	}
	return bitList
}

func (f *TMainForm) ClickOnTop(sender vcl.IObject) {
	cb := vcl.AsCheckBox(sender)
	if cb.Checked() {
		f.SetFormStyle(types.FsSystemStayOnTop)
	} else {
		f.SetFormStyle(types.FsNormal)	
	}
}
