package fceux

import (
	lua "github.com/yuin/gopher-lua"
)

func left(L *lua.LState) int {
	L.DoString("joypad.write(1,{left=true});")
}
