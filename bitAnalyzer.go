package main

import (
	"strconv"
	"fmt"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

const (
	padx = 4
	pady = 4
	bitBgX = 20
	bitBgY = 20
	bitWidth = 32
	winX = bitWidth * (bitBgX + padx)
	winY = bitBgY * 3 + pady * 2
)

var (
	color = map[string]types.TColor{
		"0": types.TColor(0xffffff),
		"1": types.TColor(0xffff55),
		"diff": types.TColor(0xff00ff),
		"same": types.TColor(0xffffff),
	}
	row = 2
)


type Bit struct {
    Memos []*vcl.TMemo
	Header *vcl.TMemo
	color map[string]types.TColor
}

func NewBit(owner vcl.IComponent, parent vcl.IWinControl, x, y, w, h int32, rows, ix int, color map[string]types.TColor) *Bit {
	bit := &Bit{}
	memos := make([]*vcl.TMemo, rows)
	for i := 0; i < rows; i++ {
		memo := vcl.NewMemo(owner)
		memo.SetParent(parent)
		memo.SetTextBuf("0")
		memo.SetColor(color["0"])
		memo.SetMaxLength(1)
		memo.SetOnClick(bit.click)
		memo.SetAlignment(types.TaCenter)
		memo.SetReadOnly(true)
		memo.SetBounds(x + int32(ix) * w, y + int32(i + 1) * h, w, h)
		memo.SetHideSelection(true)
		memo.SetComponentIndex(int32((ix + i) * bitWidth))
		memos[i] = memo
	}
	header := vcl.NewMemo(owner)
	header.SetParent(parent)
	header.SetBounds(x + int32(ix) * w, y, w, h)
	header.SetTextBuf(strconv.Itoa(bitWidth - ix - 1))
	header.SetColor(color["same"])
	header.SetMaxLength(2)
	header.SetAlignment(types.TaCenter)
	header.SetReadOnly(true)
	header.SetHideSelection(true)
	header.SetBorderStyle(types.BsNone)
	header.SetComponentIndex(int32(ix))
	bit.color = color
	bit.Header = header
	bit.Memos = memos
	return bit
}

func (b *Bit) click(sender vcl.IObject) {
	var str string
    bit := vcl.AsMemo(sender)
	bit.GetTextBuf(&str, 2)
	bit.SetMaxLength(2)
	if str == "1" {
		bit.SetTextBuf("0")
		bit.SetColor(b.color["0"])
	} else {
		bit.SetTextBuf("1")
		bit.SetColor(b.color["1"])
	}
	bit.SetAlignment(types.TaCenter)
	ix := bit.ComponentIndex()
	fmt.Print(ix)
	fmt.Print(" ")
	idx := int((ix - ix % 2 - bitWidth) / 2)
	fmt.Println(idx)
	res := make(map[string]int, 0)
	for i := 0; i < row; i++ {
		var bitString string
		mainForm.bits[idx].Memos[i].GetTextBuf(&bitString, 2)
		res[bitString] = 0
	}
	if len(res) == 1 {
		mainForm.bits[idx].Header.SetColor(color["same"])
	} else {
		mainForm.bits[idx].Header.SetColor(color["diff"])
	}
}

type TMainForm struct {
    *vcl.TForm
    bits []*Bit
	nums []*vcl.TMemo
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
	f.initComponents()
}

func (f *TMainForm) initComponents() {
	f.bits = make([]*Bit, bitWidth)
	for i := 0; i < bitWidth; i++ {
		f.bits[i] = NewBit(f, f, padx, pady, bitBgX, bitBgY, 2, i, color)
	}
	// for i := 0; i < row; i++ {
	// 	num := vcl.NewMemo(f)
	// 	num.SetParent(f)
	// 	num.SetBounds(padx + 32 * bitBgX, pady + int32(i + 2) * bitBgY, 60, bitBgY)
	// 	f.nums[i] = num
	// }
}
