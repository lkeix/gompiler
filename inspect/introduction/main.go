package main

import (
	"flag"
	"fmt"
)

func main() {
	value := flag.Int64("value", 0, "value")
	flag.Parse()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl _main\n")
	fmt.Printf("_main:\n")
	fmt.Printf("  mov rax, %d\n", *value)
	fmt.Printf("  ret\n")
}
