# -*- coding: utf-8 -*-
from tkinter import Text, Tk, Button
from tkinter.ttk import Combobox

bit_bg_map = {
    '0': '#ffffff',
    '1': '#55ffff'
}

default_param = {
    'bg': bit_bg_map['0'],
    'exportselection': False,
    'width': 2,
    'height': 1,
    'relief': 'solid',
    'takefocus': False,
}

base_map = {
    '16进制': (hex, 16),
    '10进制': (str, 10),
    '8进制': (oct, 8),
}

base_func = {
    16: hex,
    10: str,
    8: oct,
}


# def event_handler(func, **kwargs):
#     return lambda event, fun=func: fun(event, **kwargs)


def func_handler(func, **kwargs):
    return lambda fun=func: fun(**kwargs)


class FormText(Text):
    def __init__(self, master, cnf=None, **kwargs):
        if not cnf:
            cnf = {}
        self.cnf = cnf
        super().__init__(master, **kwargs)
        self.grid(
            row=self.cnf.get('row', 0),
            column=self.cnf.get('column', 0),
            padx=self.cnf.get('padx', 2),
            pady=self.cnf.get('pady', 2)
        )
        self.tag_configure("center", justify='center')
        self.insert('1.0', self.cnf.get('txt', '0'), 'center')
        self.configure(state=kwargs.get('state', 'normal'))


class BitLoc(FormText):
    def change_bit(self, _):
        self.configure(state='normal')
        var = '0' if self.get('1.0') == '1' else '1'
        self.delete('1.0')
        self.insert('1.0', var, 'center')
        self.configure(bg=bit_bg_map[var], state='disabled')


class BitHeader(FormText):
    def change_header(self, bits):
        bit_set = set(bit.get('1.0') for bit in bits)
        if len(bit_set) > 1:
            self.configure(bg='#ffaaff')
        else:
            self.configure(bg='#ffffff')


class NumRes(FormText): 
    def change_num(self, bits, base=16):
        bin_str = ''.join([bit.get('1.0') for bit in bits])
        dec_num = int(bin_str, 2)
        self.delete('1.0', 'end')
        self.insert('1.0', f'{base_func[base](dec_num)}', 'center')

    def change_bit_loc(self, bits, func, base=16, data_width=32):
        bin_str = bin(int(self.get('1.0', 'end'), base))[2:]
        bin_str = '0' * (data_width - len(bin_str)) + bin_str
        for i in range(data_width):
            if bin_str[i] != bits[i].get('1.0'):
                bits[i].change_bit(None)
        func(None)


class NumShift:
    def __init__(self, master=None, row=1, base=16, **kwargs):
        self.l_shift = Button(
            master, text='左移', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.shift, **kwargs, mode='l')
        )
        self.r_shift = Button(
            master, text='右移', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.shift, **kwargs, mode='r')
        )
        self.shift_num = Text(master, **default_param, bd=1)
        self.shift_num.tag_configure("center", justify='center')
        self.shift_num.insert('1.0', '1', 'center')
        self.l_shift.grid(row=row, column=35)
        self.shift_num.grid(row=row, column=36)
        self.r_shift.grid(row=row, column=37)
        self.base = base

    def shift(self, mode, num_text, bits, function):
        shift_num = int(self.shift_num.get('1.0', 'end'))
        num = int(num_text.get('1.0', 'end'), self.base)
        num_text.delete('1.0', 'end')
        res = num << shift_num if mode == 'l' else num >> shift_num
        num_text.insert('1.0', f'{base_func[self.base](res)}', 'center')
        num_text.change_bit_loc(bits, function, base=self.base)


class UIForm:
    def __init__(self, master, data_width=32):
        self.data_width = data_width
        self.bits = []
        self.headers = []
        self.base_sel = Combobox(master, width=6)
        self.base_sel['value'] = ('16进制', '10进制', '8进制')
        self.base_sel.current(0)
        self.base_sel.bind('<<ComboboxSelected>>', self.change_base)
        self.base_sel.grid(row=0, column=33)
        base = self.base_sel.get()
        self.base = base_map[base][1]
        for i in range(self.data_width):
            header = BitHeader(
                master, {'row': 0, 'column': i, 'txt': self.data_width - i - 1},
                **default_param, bd=0
            )
            header.configure(state='disabled')
            self.headers.append(header)
            bit_loc = [
                BitLoc(
                    master, {'row': m + 1, 'column': i}, **default_param
                ) for m in range(2)
            ]
            for bit in bit_loc:
                bit.bind('<Button-1>', bit.change_bit)
            self.bits.append(bit_loc)
        self.num_res = [
            NumRes(
                master, {
                    'row': n + 1, 'column': self.data_width + 1, 'txt': '0x0'
                }, height=1, width=15, relief='solid'
            ) for n in range(2)
        ]
        self.clear_button_f = Button(
            master, text='清除', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.clear, row=0)
        )
        self.clear_button_f.grid(row=1, column=34)
        self.clear_button_s = Button(
            master, text='清除', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.clear, row=1)
        )
        self.clear_button_s.grid(row=2, column=34)
        self.shift_f = NumShift(
            master, num_text=self.num_res[0],
            bits=[bits[0] for bits in self.bits],
            function=self.update_ui_click, base=self.base
        )
        self.shift_s = NumShift(
            master, row=2, num_text=self.num_res[1],
            bits=[bits[1] for bits in self.bits],
            function=self.update_ui_click, base=self.base
        )
        self.in_f = Button(
            master, text='取反', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.invert, row=0)
        )
        self.in_f.grid(row=1, column=38)
        self.in_s = Button(
            master, text='取反', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.invert, row=1)
        )
        self.in_s.grid(row=2, column=38)
        self.re_f = Button(
            master, text='反向', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.reverse, row=0)
        )
        self.re_f.grid(row=1, column=39)
        self.re_s = Button(
            master, text='反向', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.reverse, row=1)
        )
        self.re_s.grid(row=2, column=39)
        self.re_s = Button(
            master, text='反向', bg=bit_bg_map['0'],
            height=1, width=4, relief='solid', bd=0,
            command=func_handler(self.reverse, row=1)
        )
        self.re_s.grid(row=2, column=39)

    def update_ui_click(self, _):
        for i in range(self.data_width):
            self.headers[i].change_header(self.bits[i])
            self.num_res[0].change_num([bit[0] for bit in self.bits], base=self.base)
            self.num_res[1].change_num([bit[1] for bit in self.bits], base=self.base)

    def update_ui_key(self, event):
        if event.char == '\r':
            self.update_ui_click(event)
        num1 = self.num_res[0].get('1.0', 'end').strip()
        if 'x' in num1:
            num1 = num1[2:].strip()
        num1 = bin(int(num1 if num1 else '0x0', self.base))[2:]
        num1 = '0' * (self.data_width - len(num1)) + num1
        num2 = self.num_res[1].get('1.0', 'end').strip()
        if 'x' in num2:
            num2 = num2[2:].strip()
        num2 = bin(int(num2 if num2 else '0x0', self.base))[2:]
        num2 = '0' * (self.data_width - len(num2)) + num2
        for i in range(self.data_width):
            bits = self.bits[i]
            bit1 = bits[0].get('1.0')
            bit2 = bits[1].get('1.0')
            if num1[i] != bit1:
                self.bits[i][0].change_bit(event)
            if num2[i] != bit2:
                self.bits[i][1].change_bit(event)
            self.headers[i].change_header(self.bits[i])

    def clear(self, row):
        self.num_res[row].delete('1.0', 'end')
        self.num_res[row].insert('1.0', '0x0', 'center')
        for i in range(32):
            self.bits[i][row].configure(state='normal')
            self.bits[i][row].delete('1.0')
            self.bits[i][row].insert('1.0', '0', 'center')
            self.bits[i][row].configure(bg=bit_bg_map['0'], state='disabled')
        self.update_ui_click(None)

    def invert(self, row=0):
        for i in range(self.data_width):
            self.bits[i][row].change_bit(None)
        self.update_ui_click(None)

    def reverse(self, row=0):
        bins = [self.bits[i][row].get('1.0') for i in range(self.data_width - 1, -1, -1)]
        for i in range(self.data_width):
            self.bits[i][row].configure(state='normal')
            self.bits[i][row].delete('1.0')
            self.bits[i][row].insert('1.0', bins[i], 'center')
            self.bits[i][row].configure(bg=bit_bg_map[bins[i]], state='disabled')
        self.update_ui_click(None)

    def change_base(self, _):
        num1 = int(self.num_res[0].get('1.0', 'end'), self.base)
        num2 = int(self.num_res[1].get('1.0', 'end'), self.base)
        base = self.base_sel.get()
        self.num_res[0].delete('1.0', 'end')
        self.num_res[1].delete('1.0', 'end')
        self.num_res[0].insert('1.0', base_map[base][0](num1), 'center')
        self.num_res[1].insert('1.0', base_map[base][0](num2), 'center')
        self.base = base_map[base][1]
        self.shift_f.base = base_map[base][1]
        self.shift_s.base = base_map[base][1]


if __name__ == '__main__':
    tk = Tk()
    tk.title('寄存器对比')
    tk.configure(bg='#ffffff')
    form = UIForm(tk)
    tk.bind('<Button-1>', form.update_ui_click)
    tk.bind('<Key>', form.update_ui_key, add='+')
    tk.resizable(False, False)
    tk.mainloop()
