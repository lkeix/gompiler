package main

import (
	"flag"
	"fmt"
	"strconv"
	"unicode"
)

type (
	TokenKind int

	Token struct {
		kind TokenKind
		next *Token
		str  string
		val  int
	}
)

const (
	Add TokenKind = iota
	Sub
	Mul
	Div
	Reversed
	Num
	EOF
)

var token *Token

// Tokenizer
func consume(op byte) bool {
	if token.kind != Reversed || token.str[len(token.str)-1] != op {
		return false
	}
	token = token.next
	return true
}

func expect(op byte) {
	if token.kind != Reversed || token.str[len(token.str)-1] != op {
		panic("Error: expected " + string(op))
	}
	token = token.next
}

func newToken(kind TokenKind, current *Token, str string) *Token {
	next := &Token{
		kind: kind,
		str:  str,
	}
	current.next = next
	return next
}

func atEOF() bool {
	return token.kind == EOF
}

func expectNumber() int {
	if token.kind != Num {
		panic("Error: expected a number")
	}
	val := token.val
	token = token.next
	return val
}

func tokenize(input string) *Token {
	head := &Token{}
	current := head
	extractNumber := func(i int) int {
		for ; i < len(input); i++ {
			if !unicode.IsDigit(rune(input[i])) {
				break
			}
		}
		return i
	}
	for i := 0; i < len(input); i++ {
		if unicode.IsSpace(rune(input[i])) {
			continue
		}

		if input[i] == '+' ||
			input[i] == '-' ||
			input[i] == '*' ||
			input[i] == '/' ||
			input[i] == '(' ||
			input[i] == ')' {
			current = newToken(Reversed, current, input[:i+1])
			continue
		}

		if unicode.IsDigit(rune(input[i])) {
			end := extractNumber(i)
			current = newToken(Num, current, input[:i])
			val, _ := strconv.Atoi(input[i:end])
			current.val = val
			i = end - 1
			continue
		}

		panic("Error: cannot tokenize")
	}

	newToken(EOF, current, input)
	return head.next
}

// Parser
type (
	NodeKind int

	Node struct {
		kind NodeKind
		lhs  *Node
		rhs  *Node
		num  int
	}
)

const (
	NAdd NodeKind = iota
	NSub
	NMul
	NDiv
	NNum
)

func newNode(kind NodeKind) *Node {
	return &Node{
		kind: kind,
	}
}

func newBinary(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNum(num int) *Node {
	return &Node{
		kind: NNum,
		num:  num,
	}
}

func expr() *Node {
	node := mul()

	for {
		if consume('+') {
			node = newBinary(NAdd, node, mul())
			continue
		}
		if consume('-') {
			node = newBinary(NSub, node, mul())
			continue
		}
		return node
	}
}

func mul() *Node {
	node := primary()

	for {
		if consume('*') {
			node = newBinary(NMul, node, primary())
		}
		if consume('/') {
			node = newBinary(NDiv, node, primary())
		}
		return node
	}
}

func primary() *Node {
	if consume('(') {
		node := expr()
		expect(')')
		return node
	}

	return newNum(expectNumber())
}

func gen(node *Node) {
	if node.kind == NNum {
		fmt.Printf("  push %d\n", node.num)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case NAdd:
		fmt.Printf("  add rax, rdi\n")
	case NSub:
		fmt.Printf("  sub rax, rdi\n")
	case NMul:
		fmt.Printf("  imul rax, rdi\n")
	case NDiv:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rax, rdi\n")
	}

	fmt.Printf("  push rax\n")
}

func main() {
	input := flag.String("input", "", "")
	flag.Parse()

	token = tokenize(*input)
	node := expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global _main\n")
	fmt.Printf("_main:\n")
	gen(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
