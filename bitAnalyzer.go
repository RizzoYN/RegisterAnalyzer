package main

import (
	"fmt"
	"strconv"
	"strings"

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
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
)

type TMainForm struct {
	*vcl.TForm
	BitLocs        [][]*vcl.TMemo
	BitHeader      []*vcl.TMemo
	BitNum         []*vcl.TEdit
	BaseChoise     *vcl.TRadioGroup
	LeftShfits     []*vcl.TButton
	ShiftNums      []*vcl.TEdit
	RightShfits    []*vcl.TButton
	ReverseButtons []*vcl.TButton
	InvertButtons  []*vcl.TButton
	ClearButtons   []*vcl.TButton
	base           int
	AddRow         *vcl.TButton
	RmRow          *vcl.TButton
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
	f.initComponents(f, f, bitWidth, Row, color)
}

func (f *TMainForm) newMemo(owner vcl.IComponent, parent vcl.IWinControl, x, y, w, h, ix int32, row int, color types.TColor, text string) *vcl.TMemo {
	memo := vcl.NewMemo(owner)
	menu := vcl.NewPopupMenu(owner)
	var maxLength int32
	if row == 0 {
		maxLength = 3
	} else {
		maxLength = 1
	}
	memo.SetParent(parent)
	memo.SetPopupMenu(menu)
	memo.SetTextBuf(text)
	memo.SetColor(color)
	memo.SetMaxLength(maxLength)
	memo.SetAlignment(types.TaCenter)
	memo.SetReadOnly(true)
	memo.SetBounds(x+ix%bitWidth*w, y+int32(row)*h, w, h)
	memo.SetComponentIndex(ix)
	if ix < bitWidth {
		memo.SetBorderStyle(types.BsNone)
		memo.SetHeight(17)
		memo.SetControlState(types.CsNoStdEvents)
	} else {
		memo.SetOnClick(f.Clicked)
	}
	memo.SetHideSelection(true)
	return memo
}

func (f *TMainForm) initComponents(owner vcl.IComponent, parent vcl.IWinControl, cols, rows int, color map[string]types.TColor) {
	f.base = 16
	headers := make([]*vcl.TMemo, cols)
	bits := make([][]*vcl.TMemo, rows)
	nums := make([]*vcl.TEdit, rows)
	lshifts := make([]*vcl.TButton, rows)
	shiftnums := make([]*vcl.TEdit, rows)
	rshifts := make([]*vcl.TButton, rows)
	reverse := make([]*vcl.TButton, rows)
	invert := make([]*vcl.TButton, rows)
	clear := make([]*vcl.TButton, rows)
	for r := 0; r <= rows; r++ {
		if r == 0 {
			for col := 0; col < cols; col++ {
				ix := int32(r*bitWidth + col)
				headers[col] = f.newMemo(owner, parent, padx, pady, bitBgX, bitBgY, ix, r, color["same"], fmt.Sprint(bitWidth-col-1))
			}
		} else {
			bitRow := make([]*vcl.TMemo, cols)
			for col := 0; col < cols; col++ {
				ix := int32(r*bitWidth + col)
				bitRow[col] = f.newMemo(owner, parent, padx, pady, bitBgX, bitBgY, ix, r, color["0"], "0")
			}
			bits[r-1] = bitRow
			numEdit := vcl.NewEdit(owner)
			numEdit.SetParent(parent)
			numEdit.SetBounds(padx+bitWidth*bitBgX, padx+int32(r)*bitBgY, bitNumEdX, bitBgY)
			numEdit.SetOnKeyUp(f.Typed)
			numEdit.SetName(fmt.Sprintf("numEdit%d", r-1))
			numEdit.SetTextBuf("0")
			lshift := vcl.NewButton(owner)
			lshift.SetParent(parent)
			lshift.SetBounds(padx+bitWidth*bitBgX+bitNumEdX, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			lshift.SetTextBuf("<<")
			lshift.SetOnClick(f.ClickShift)
			lshift.SetName(fmt.Sprintf("lshift%d", r-1))
			shiftnum := vcl.NewEdit(owner)
			shiftnum.SetParent(parent)
			shiftnum.SetBounds(padx+bitWidth*bitBgX+bitNumEdX+ButtonS, padx+int32(r)*bitBgY, bitBgX, bitBgY)
			shiftnum.SetTextBuf("1")
			shiftnum.SetMaxLength(2)
			shiftnum.SetAlignment(types.TaCenter)
			rshift := vcl.NewButton(owner)
			rshift.SetParent(parent)
			rshift.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			rshift.SetTextBuf(">>")
			rshift.SetOnClick(f.ClickShift)
			rshift.SetName(fmt.Sprintf("rshift%d", r-1))
			rev := vcl.NewButton(owner)
			rev.SetParent(parent)
			rev.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*2, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			rev.SetTextBuf("倒序")
			rev.SetOnClick(f.ClickReverse)
			rev.SetName(fmt.Sprintf("rev%d", r-1))
			invt := vcl.NewButton(owner)
			invt.SetParent(parent)
			invt.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*3, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			invt.SetTextBuf("转换")
			invt.SetOnClick(f.ClickInvert)
			invt.SetName(fmt.Sprintf("invt%d", r-1))
			cler := vcl.NewButton(owner)
			cler.SetParent(parent)
			cler.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*4, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			cler.SetTextBuf("清空")
			cler.SetOnClick(f.ClickClear)
			cler.SetName(fmt.Sprintf("cler%d", r-1))
			nums[r-1] = numEdit
			lshifts[r-1] = lshift
			shiftnums[r-1] = shiftnum
			rshifts[r-1] = rshift
			reverse[r-1] = rev
			invert[r-1] = invt
			clear[r-1] = cler
		}
	}
	checkgroup := vcl.NewRadioGroup(owner)
	checkgroup.SetParent(parent)
	checkgroup.SetCaption("进制")
	checkgroup.SetBounds(int32(winX/2-60), winY-50, 120, 46)
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
	addrow := vcl.NewButton(owner)
	addrow.SetParent(parent)
	addrow.SetBounds(padx, winY-25-bitBgY/2, ButtonS*2, bitBgY)
	addrow.SetTextBuf("增加一行")
	addrow.SetOnClick(f.AddR)
	if Row == 9 {
		addrow.SetEnabled(false)
	}
	rmrow := vcl.NewButton(owner)
	rmrow.SetParent(parent)
	rmrow.SetBounds(padx+ButtonS*2, winY-25-bitBgY/2, ButtonS*2, bitBgY)
	if Row == 1 {
		rmrow.SetEnabled(false)
	}
	rmrow.SetTextBuf("删除一行")
	rmrow.SetOnClick(f.RemoveR)
	f.BaseChoise = checkgroup
	f.BitHeader = headers
	f.BitLocs = bits
	f.BitNum = nums
	f.LeftShfits = lshifts
	f.ShiftNums = shiftnums
	f.RightShfits = rshifts
	f.ReverseButtons = reverse
	f.InvertButtons = invert
	f.ClearButtons = clear
	f.AddRow = addrow
	f.RmRow = rmrow
}

func (f *TMainForm) Typed(sender vcl.IObject, key *types.Char, shift types.TShiftState) {
	var str string
	num := vcl.AsEdit(sender)
	num.GetTextBuf(&str, bitWidth)
	rowIx := f.GetRowIndex(num)
	resNum, err := strconv.ParseInt(str, f.base, bitWidth*2)
	if err != nil && str != "" {
		bitList := f.GetBitString(int(rowIx))
		binStr := strings.Join(bitList, "")
		bin, _ := strconv.ParseInt(binStr, 2, bitWidth*2)
		f.UpdateBitNum(bin, rowIx)
		var bitString string
		f.BitNum[rowIx].GetTextBuf(&bitString, bitWidth*2)
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
	ix := bit.ComponentIndex() % bitWidth
	rowIx := int(bit.ComponentIndex()/32) - 1
	for i := 0; i < Row; i++ {
		var bitString string
		f.BitLocs[i][ix].GetTextBuf(&bitString, 2)
		bitMap[bitString] = 0
	}
	f.UpdateHeader(bitMap, int(ix))
	bitList := f.GetBitString(rowIx)
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
		f.BitNum[i].GetTextBuf(&bitString, bitWidth)
		num, _ := strconv.ParseInt(bitString, oldbase, bitWidth*2)
		f.UpdateBitNum(num, int64(i))
	}
}

func (f *TMainForm) ClickClear(sender vcl.IObject) {
	cler := vcl.AsButton(sender)
	rowIx := f.GetRowIndex(cler)
	f.BitNum[rowIx].SetTextBuf("0")
	for c := 0; c < bitWidth; c++ {
		bitMap := make(map[string]int, Row)
		for i := 0; i < Row; i++ {
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
	f.ShiftNums[rowIx].GetTextBuf(&str, 8)
	shiftNum, _ := strconv.ParseInt(str, 10, 16)
	f.BitNum[rowIx].GetTextBuf(&str, bitWidth*32)
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
	if Row == 1 {
		f.RmRow.SetEnabled(true)
	}
	Row++
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
	f.Free()
	mainForm.Free()
	vcl.Application.Terminate()
	main()
}

func (f *TMainForm) RemoveR(sender vcl.IObject) {
	if Row > 1 {
		Row--
	}
	winY = int32(bitBgY*(Row+1)+pady*2) + 50
	f.Free()
	mainForm.Free()
	vcl.Application.Terminate()
	main()
}

func (f *TMainForm) GetRowIndex(sender vcl.IWinControl) int64 {
	cname := sender.Name()
	name := string(cname[len(cname)-1])
	rowIx, _ := strconv.ParseInt(name, 10, 0)
	return rowIx
}

func (f *TMainForm) UpdateBitNum(bin, r int64) {
	f.BitNum[r].Clear()
	switch f.base {
	case 16:
		f.BitNum[r].SetTextBuf(fmt.Sprintf("%x", bin))
	case 10:
		f.BitNum[r].SetTextBuf(fmt.Sprint(bin))
	case 8:
		f.BitNum[r].SetTextBuf(fmt.Sprintf("%o", bin))
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

func (f *TMainForm) GetBitString(r int) []string {
	bitList := make([]string, bitWidth)
	for i := 0; i < bitWidth; i++ {
		var bitString string
		f.BitLocs[r][i].GetTextBuf(&bitString, 2)
		bitList[i] = bitString
	}
	return bitList
}
