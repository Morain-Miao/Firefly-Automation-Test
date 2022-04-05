//testplugin.go
package main

import (
	"fmt"
	"plugin"
)

func init() {
	fmt.Println("world")
	//我们还可以做其他更高阶的事情，比如 platform.RegisterPlugin({"func": Hello}) 之类的，向插件平台自动注册该插件的函数
}

func Hello() {
	fmt.Println("hello")
}

func main() {
	p, err := plugin.Open("testplugin.so")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("Hello")
	if err != nil {
		panic(err)
	}

	f.(func())()
}
