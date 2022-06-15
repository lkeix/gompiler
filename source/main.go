package main

import "os"

var hoge int = 30
var fuga int = 12
var errorMessage string = "to stderr\n"

func main() {
	var a int
	var str string
	a = 1
	str = "aaa"
	print(a)
	print(str)
	print("hello world\n")
	print(errorMessage)
	os.Exit(hoge + fuga)
}
