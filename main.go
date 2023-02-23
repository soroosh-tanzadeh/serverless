package main

import (
	"fmt"
	v8 "rogchap.com/v8go"
)

func main() {
	iso := v8.NewIsolate()
	context := v8.NewContext(iso)
	context.RunScript("const multiply = (a, b) => a * b", "math.js")
	value, _ := context.RunScript("multiply(2,6)", "main.js")
	fmt.Println(value.Integer())
}
