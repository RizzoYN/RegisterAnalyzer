package main

import (
	"fmt"
	"strconv"
	"strings"


	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/keys"
)

const (
	padx      = 4
	pady      = 4
	bitBgX    = 20
	bitBgY    = 25
	bitWidth  = 32
	bitNumEdX = int32(2.5 * bitWidth) + 10
	ButtonS   = 30
	winX      = (bitWidth+1)*bitBgX+2*padx + bitNumEdX + 5*ButtonS
	winY      = bitBgY*3 + pady*2 + 50 // bitBgY*(Row+1) + pady*Row + 50
)

var (
	color = map[string]types.TColor{
		"0":    types.TColor(0xffffff),
		"1":    types.TColor(0xffff55),
		"diff": types.TColor(0xff00ff),
		"same": types.TColor(0xf0f0f0),
	}
	Row = 2
	FirstIdx = Row * bitWidth + 64
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
}

var mainForm *TMainForm

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&mainForm)
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
		memo.SetHeight(18)
		memo.SetControlState(types.CsNoStdEvents)
	} else {
		memo.SetOnClick(f.Clicked)
	}
	memo.SetHideSelection(true)
	return memo
}

func (f *TMainForm) initComponents(owner vcl.IComponent, parent vcl.IWinControl, cols, rows int, color map[string]types.TColor) {
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
			numEdit.SetOnKeyPress(f.Typed)
			lshift := vcl.NewButton(owner)
			lshift.SetParent(parent)
			lshift.SetBounds(padx+bitWidth*bitBgX+bitNumEdX, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			lshift.SetTextBuf("<<")
			lshift.SetOnClick(f.ClickShift)
			shiftnum := vcl.NewEdit(owner)
			shiftnum.SetParent(parent)
			shiftnum.SetBounds(padx+bitWidth*bitBgX+bitNumEdX+ButtonS, padx+int32(r)*bitBgY, bitBgX, bitBgY)
			shiftnum.SetTextBuf("1")
			shiftnum.SetAlignment(types.TaCenter)
			rshift := vcl.NewButton(owner)
			rshift.SetParent(parent)
			rshift.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			rshift.SetTextBuf(">>")
			rshift.SetOnClick(f.ClickShift)
			rev := vcl.NewButton(owner)
			rev.SetParent(parent)
			rev.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*2, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			rev.SetTextBuf("倒序")
			rev.SetOnClick(f.ClickReverse)
			invt := vcl.NewButton(owner)
			invt.SetParent(parent)
			invt.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*3, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			invt.SetTextBuf("转换")
			invt.SetOnClick(f.ClickInvert)
			cler := vcl.NewButton(owner)
			cler.SetParent(parent)
			cler.SetBounds(padx+(bitWidth+1)*bitBgX+bitNumEdX+ButtonS*4, padx+int32(r)*bitBgY, ButtonS, bitBgY)
			cler.SetTextBuf("清空")
			cler.SetOnClick(f.ClickClear)
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
	f.base = 16
}

func (f *TMainForm) Typed(sender vcl.IObject, key *types.Char) {
	var str string
	num := vcl.AsEdit(sender)
	num.GetTextBuf(&str, bitWidth)
	keyNum := rune(*key)
	res := str + string(keyNum)
	if (keyNum == keys.VkBack) && (len(str) > 0) {
		res = str[:len(str) - 1]
	}
	fmt.Println(res)
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
	if len(bitMap) == 1 {
		f.BitHeader[ix].SetColor(color["same"])
	} else {
		f.BitHeader[ix].SetColor(color["diff"])
	}
	bitList := make([]string, bitWidth)
	for i := 0; i < bitWidth; i++ {
		var bitString string
		f.BitLocs[rowIx][i].GetTextBuf(&bitString, 2)
		bitList[i] = bitString
	}
	binStr := strings.Join(bitList, "")
	bin, _ := strconv.ParseInt(binStr, 2, bitWidth)
	switch f.base {
	case 16:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%x", bin))
	case 10:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprint(bin))
	case 8:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%o", bin))
	}
}

func (f *TMainForm) BaseChange(sender vcl.IObject) {
	var str string
	ra := vcl.AsRadioButton(sender)
	ra.GetTextBuf(&str, 4)
	for i := 0; i < Row; i++ {
		var bitString string
		f.BitNum[i].GetTextBuf(&bitString, bitWidth)
		num, _ := strconv.ParseInt(bitString, f.base, bitWidth)
		switch str {
		case "16":
			f.BitNum[i].SetTextBuf(fmt.Sprintf("%x", num))
		case "10":
			f.BitNum[i].SetTextBuf(fmt.Sprint(num))
		case "8":
			f.BitNum[i].SetTextBuf(fmt.Sprintf("%o", num))
		}
	}
	base, _ := strconv.ParseInt(str, 10, 12)
	f.base = int(base)
}

func (f *TMainForm) ClickClear(sender vcl.IObject) {
	cler := vcl.AsButton(sender)
	rowIx := int((cler.ComponentIndex() - int32(FirstIdx)) / 49)
	bitMap := make(map[string]int, Row)
	f.BitNum[rowIx].SetTextBuf("0")
	for i := 0; i < bitWidth; i++ {
		f.BitLocs[rowIx][i].SetTextBuf("0")
		f.BitLocs[rowIx][i].SetColor(color["0"])
	}
	for c := 0; c < bitWidth; c++ {
		for i := 0; i < Row; i++ {
			var bitString string
			f.BitLocs[i][c].GetTextBuf(&bitString, 2)
			bitMap[bitString] = 0
		}
		if len(bitMap) == 1 {
			f.BitHeader[c].SetColor(color["same"])
		} else {
			f.BitHeader[c].SetColor(color["diff"])
		}
	}
}

func (f *TMainForm) ClickInvert(sender vcl.IObject) {
	inv := vcl.AsButton(sender)
	rowIx := int((inv.ComponentIndex() - int32(FirstIdx)) / 49)
	for i := 0; i < bitWidth; i++ {
		f.Clicked(f.BitLocs[rowIx][i])
	}
}

func (f *TMainForm) ClickShift(sender vcl.IObject) {
	shift := vcl.AsButton(sender)
	rowIx := int((shift.ComponentIndex() - int32(FirstIdx)) / 49)
	col := (shift.ComponentIndex() - int32(FirstIdx)) % 39
	var str string
	f.ShiftNums[rowIx].GetTextBuf(&str, 8)
	shiftNum, _ := strconv.ParseInt(str, 10, 16)
	f.BitNum[rowIx].GetTextBuf(&str, bitWidth*32)
	num, _ := strconv.ParseInt(str, f.base, bitWidth)
	switch col {
	case 33:
		num <<= shiftNum
	case 35:
		num >>= shiftNum
	}
	f.BitNum[rowIx].Clear()
	switch f.base {
	case 16:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%x", num))
	case 10:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprint(num))
	case 8:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%o", num))
	}
	
}

func (f *TMainForm) ClickReverse(sender vcl.IObject) {
	rev := vcl.AsButton(sender)
	rowIx := int((rev.ComponentIndex() - int32(FirstIdx)) / 49)
	bins := make([]string, bitWidth)
	for c := 0; c < bitWidth; c++ {
		bitMap := make(map[string]int, Row)
		for i := 0; i < Row; i++ {
			var str string
			if i == rowIx {
				f.BitLocs[i][c].GetTextBuf(&str, 2)
				bins[bitWidth - c - 1] = str
			} else {
				f.BitLocs[i][bitWidth - c - 1].GetTextBuf(&str, 2)
			}
			bitMap[str] = 0
		}
		if len(bitMap) == 1 {
			f.BitHeader[bitWidth - c - 1].SetColor(color["same"])
		} else {
			f.BitHeader[bitWidth - c - 1].SetColor(color["diff"])
		}
	}
	binStr := strings.Join(bins, "")
	bin, _ := strconv.ParseInt(binStr, 2, bitWidth)
	switch f.base {
	case 16:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%x", bin))
	case 10:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprint(bin))
	case 8:
		f.BitNum[rowIx].SetTextBuf(fmt.Sprintf("%o", bin))
	}
	for i, bin := range(bins) {
		f.BitLocs[rowIx][i].SetTextBuf(bin)
		f.BitLocs[rowIx][i].SetColor(color[bin])
	}
}
