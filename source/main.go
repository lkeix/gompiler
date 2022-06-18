package main

import "os"

var (
	globalstring string = "I'm global string\n"
	globalint1   int    = 30
	globalint2   int    = 0
)

func f1(x int) int {
	return x + 1
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

	print("end!\n")

	os.Exit(globalint1 + globalint2 + localint1)
}
