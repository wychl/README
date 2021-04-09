package main

import lua "github.com/yuin/gopher-lua"

func main() {
	runLuaScriptFile()
}

func runLuaScript() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
}

func runLuaScriptFile() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile("hello.lua"); err != nil {
		panic(err)
	}
}
