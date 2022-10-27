package main

import (
	"flag"
	"fmt"
	"strconv"
)

type (
	TokenKind int
	Token     struct {
		Kind  TokenKind
		Next  *Token
		Value int
		Str   string
	}
)

const (
	Reversed TokenKind = iota
	Number
	EOF
)

var token *Token

func consume(op byte) bool {
	if token.Kind != Reversed || token.Str[len(token.Str)-1] != op {
		return false
	}
	token = token.Next
	return true
}

func expect(op byte) {
	if token.Kind != Reversed || token.Str[len(token.Str)-1] != op {
		panic("Error: " + string(op))
	}
	token = token.Next
}

func expectNumber() int {
	if token.Kind != Number {
		panic("Error: is not number")
	}
	val := token.Value
	token = token.Next
	return val
}

func atEOF() bool {
	return token.Kind == EOF
}

func newToken(kind TokenKind, current *Token, str string) *Token {
	token := &Token{}
	token.Kind = kind
	current.Next = token
	token.Str = str
	return token
}

func Tokenize(code string) *Token {
	head := &Token{}

	current := head

	extractNumber := func(start int, str string) int {
		for ; start < len(str); start++ {
			if _, err := strconv.Atoi(string(str[start])); err != nil {
				break
			}
		}
		return start
	}
	for i := 0; i < len(code); i++ {
		if code[i] == ' ' {
			continue
		}

		if code[i] == '+' || code[i] == '-' {
			current = newToken(Reversed, current, string(code[:i+1]))
			continue
		}

		if _, err := strconv.Atoi(string(code[i])); err == nil {
			current = newToken(Number, current, string(code[:i]))
			end := extractNumber(i, code)
			val, _ := strconv.Atoi(code[i:end])
			current.Value = val
			i = end
			continue
		}

		panic("Error: cannot tokenize")

	}

	newToken(EOF, current, code)
	return head.Next
}

func strToInt(code string) (int, int) {
	conv := ""
	i := 0
	for ; i < len(code); i++ {
		if _, err := strconv.Atoi(string(code[i])); err != nil {
			break
		}
		conv += string(code[i])
	}

	val, _ := strconv.Atoi(conv)
	return i, val
}

func main() {
	input := flag.String("input", "input", "")

	flag.Parse()

	token = Tokenize(*input)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global _main\n")
	fmt.Printf("_main:\n")
	fmt.Printf("  mov rax, %d\n", expectNumber())

	for !atEOF() {
		if consume('+') {
			fmt.Printf("  add rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Printf("  sub rax, %d\n", expectNumber())
	}

	fmt.Printf("  ret\n")
}
