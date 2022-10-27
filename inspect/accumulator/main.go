package main

import (
	"flag"
	"fmt"
	"strconv"
)

const (
	add = '+'
	sub = '-'
)

func main() {
	equivalent := flag.String("equivalent", "", "equivalent")
	flag.Parse()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl _main\n")
	fmt.Printf("_main:\n")
	i := 0
	extractValueStr := func() string {
		strVal := ""
		for ; i < len(*equivalent); i++ {
			if _, err := strconv.Atoi(string((*equivalent)[i])); err != nil {
				return strVal
			}
			strVal += string((*equivalent)[i])
		}
		return strVal
	}

	fmt.Printf("  mov rax, %s\n", extractValueStr())

	for ; i < len(*equivalent); i++ {
		if (*equivalent)[i] == add {
			i++
			fmt.Printf("  add rax, %s\n", extractValueStr())
			i--
			continue
		}

		if (*equivalent)[i] == sub {
			i++
			fmt.Printf("  sub rax, %s\n", extractValueStr())
			i--
			continue
		}

		if _, err := strconv.Atoi(string((*equivalent)[i])); err != nil {
			panic(err)
		}
	}
	fmt.Printf("  ret\n")
}
