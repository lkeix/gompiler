package main

import "os"

var hoge int = 30
var fuga int = 12
var errorMessage string = "to stderr\n"

func main() {
	print("hello world\n")
	os.Exit(hoge + fuga)
}
