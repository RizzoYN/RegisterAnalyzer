package main

import (
	"fmt"
	"strings"
	// "strconv"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

const (
	padx     = 4
	pady     = 4
	bitBgX   = 20
	bitBgY   = 20
	bitWidth = 32
	winX     = bitWidth * (bitBgX + padx)
	winY     = bitBgY * 3 + pady * 2
)

var (
	color = map[string]types.TColor{
		"0":    types.TColor(0xffffff),
		"1":    types.TColor(0xffff55),
		"diff": types.TColor(0xff00ff),
		"same": types.TColor(0xf0f0f0),
	}
	Row = 2
)


type TMainForm struct {
	*vcl.TForm
	BitLocs        [][]*vcl.TMemo
	BitHeader      []*vcl.TMemo
	BitNum         []*vcl.TEdit
	BaseChoise     *vcl.TCheckGroup
	LeftShfits     []*vcl.TButton
	ShiftNums      []*vcl.TEdit
	RightShfits    []*vcl.TButton
	ReverseButtons []*vcl.TButton
	InvertButtons  []*vcl.TButton
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

func (f *TMainForm)newMemo(owner vcl.IComponent, parent vcl.IWinControl, x, y, w, h, ix int32, row int, color types.TColor, text string) *vcl.TMemo {
	memo := vcl.NewMemo(owner)
	var maxLength int32
	if row == 0 {
		maxLength = 3
	} else {
		maxLength = 1
	}
	memo.SetParent(parent)
	memo.SetTextBuf(text)
	memo.SetColor(color)
	memo.SetMaxLength(maxLength)
	memo.SetAlignment(types.TaCenter)
	memo.SetReadOnly(true)
	memo.SetBounds(x + ix % bitWidth * w, y + int32(row) * h, w, h)
	memo.SetHideSelection(true)
	memo.SetComponentIndex(ix)
	if ix < bitWidth {
		memo.SetBorderStyle(types.BsNone)
		memo.SetControlState(types.CsNoStdEvents)
	} else {
		memo.SetOnClick(f.Clicked)
	}
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
	for r := 0; r <= rows; r++ {
		if r == 0 {
			for col := 0; col < cols; col++ {
				ix := int32(r * bitWidth + col)
				headers[col] = f.newMemo(owner, parent, padx, pady, bitBgX, bitBgY, ix, r, color["same"], fmt.Sprint(bitWidth - col - 1))
			}
			
		} else {
			bitRow := make([]*vcl.TMemo, cols)
			for col := 0; col < cols; col++ {
				ix := int32(r * bitWidth + col)
				bitRow[col] = f.newMemo(owner, parent, padx, pady, bitBgX, bitBgY, ix, r, color["0"], "0")
			}
			bits[r - 1] = bitRow
			nums[r - 1] = nil
		}				
	}
	f.BitHeader = headers
	f.BitLocs = bits
	f.BitNum = nums
	f.LeftShfits = lshifts
	f.ShiftNums = shiftnums
	f.RightShfits = rshifts
	f.ReverseButtons = reverse
	f.InvertButtons = invert
}

func (f *TMainForm) Typed(sender vcl.IObject, key *types.Char, shift types.TShiftState) {
	var str string
	num := vcl.AsEdit(sender)
	num.GetTextBuf(&str, 32)
	res := str + string(rune(*key))
	fmt.Println(res)
}

func (f *TMainForm)Clicked(sender vcl.IObject) {
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
	rowIx := int(bit.ComponentIndex() / 32) - 1
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
	fmt.Println(binStr)
}
