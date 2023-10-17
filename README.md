# RegisterAnalyzer
32位寄存器bit位对比工具
使用tkinter,无第三方库,推荐python>=3.9

go为GoVcl库实现,需要动态库,https://github.com/ying32/govcl/releases

fltk为go-fltk实现,windows下需要mingw64 编译包含dll文件: go build -ldflags="-H windowsgui -s -w -linkmode external -extldflags -static" 
