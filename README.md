# golang版本的GUI开发包tk9.0
tcl/tk9.0版本的tk开发GUI,是一个正在完善的go gui开发包

### 特点
- 无cgo跨平台编译,支持win,macos,linux
- 官方地址 <https://gitlab.com/cznic/tk9.0> 这只是一个镜像
- 官方文档 <https://pkg.go.dev/modernc.org/tk9.0>
- 这是一个移植的2024-9-28发行的tcl/tk9.0版本,后续对托盘菜单,系统通知等肯定支持,目前golang gui在更新的库在这两方面表现都很差.
- 这个包编译的程序并不大也不含有任何第三方库,也无cgo,
- 经过测试目前支持菜单,文本框,文本域,打开,保存,调用对话框都能正常使用,托盘菜单,系统通知目前还没移植完成.可以编写一些GUI程序,漂亮效果目前只有按钮,等待官方继续开发
- 一个简单demo,编译后6M,使用`go build -ldflags "-s -w -H windowsgui"` upx压缩后5m,能够接受.
~~~
package main

import (
	. "modernc.org/tk9.0"
	"modernc.org/tk9.0/b5"
)

func main() {
	background := White
	primary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#0d6efd", b5.ButtonFocus: "#98c1fe"}
	secondary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#6c757d", b5.ButtonFocus: "#c0c4c8"}
	success := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#198754", b5.ButtonFocus: "#9dccb6"}

	StyleThemeUse("default")
	opts := Opts{}
	Grid(TButton(Txt("Primary"), Style(b5.ButtonStyle("primary.TButton", primary, background, false))),
		TButton(Txt("Secondary"), Style(b5.ButtonStyle("secondary.TButton", secondary, background, false))),
		TButton(Txt("Success"), Style(b5.ButtonStyle("success.TButton", success, background, false))),
		opts)
	Grid(TExit(Txt("退出")), Columnspan(3), opts)
	App.Configure(Background(background)).Wait()
}
~~~

