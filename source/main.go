package main

import "os"

var globalstring string = "I'm global string\n"
var globalint1 int = 30
var globalint2 int = 2

func f1(x int) int {
	return x + 1
}

func sum(a int, b int) int {
	return a + b
}

func join(a string, b string) string {
	return a + b
}

func returnString() string {
	return "aaaaa\n"
}

func main() {
	globalint2 = f1(1)

	var localstring1 string

	print("Start!\n")
	localstring1 = "I'm local string1\n"
	print(globalstring)
	print(localstring1)
	localstring1 = "local string1 changed\n"
	globalstring = "globalstring changed\n"
	print(localstring1)
	print(globalstring)

	var localint1 int
	localint1 = 10
	var tmp string
	tmp = returnString()
	print(tmp)

	/*
		var joined string
		joined = join(globalstring, localstring1)
		globalstring = "concat string for localint1 "
		print(joined)
	*/

	print("end!\n")

	os.Exit(globalint1 + sum(localint1, globalint2))
}
