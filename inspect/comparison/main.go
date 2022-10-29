package main

import (
	"flag"
	"fmt"
	"strconv"
	"unicode"
)

func arg() string {
	input := flag.String("input", "", "")
	flag.Parse()
	return *input
}

// tokenize
type (
	TokenKind int

	Token struct {
		kind  TokenKind
		next  *Token
		token byte
		val   int
	}
)

const (
	Add TokenKind = iota
	Sub
	Mul
	Reserved
	Number
	EOF
)

func consume(token *Token, op byte) (*Token, bool) {
	if token.kind != Reserved || token.token != op {
		return token, false
	}
	return token.next, true
}

func expect(token *Token, op byte) *Token {
	if token.kind != Reserved || token.token != op {
		panic("Error: expected " + string(op))
	}
	return token.next
}

func expectNumber(token *Token) (*Token, int) {
	if token.kind != Number {
		panic("Error: expected a number")
	}
	val := token.val
	return token.next, val
}

func newToken(kind TokenKind, current *Token, token byte) *Token {
	next := &Token{
		kind:  kind,
		token: token,
	}
	current.next = next
	return next
}

func eof(token *Token) bool {
	return token.kind == EOF
}

func tokenize(code string) *Token {
	head := &Token{}
	current := head
	extractNumber := func(i int) int {
		for ; i < len(code); i++ {
			if !unicode.IsDigit(rune(code[i])) {
				break
			}
		}
		return i
	}

	for i := 0; i < len(code); i++ {
		if unicode.IsSpace(rune(code[i])) {
			continue
		}

		if code[i] == '+' ||
			code[i] == '-' ||
			code[i] == '*' ||
			code[i] == '/' ||
			code[i] == '(' ||
			code[i] == ')' {
			current = newToken(Reserved, current, code[i])
			continue
		}

		if unicode.IsDigit(rune(code[i])) {
			j := extractNumber(i)
			val, _ := strconv.Atoi(code[i:j])
			current = newToken(Number, current, code[i])
			current.val = val
			i = j - 1
			continue
		}
		panic("Error: cannot tokenize")
	}
	current = newToken(EOF, current, code[len(code)-1])
	return head.next
}

// parse
type (
	NodeKind int

	Node struct {
		kind   NodeKind
		lhs    *Node
		rhs    *Node
		number int
	}

	parser struct {
		token *Token
	}
)

const (
	NAdd NodeKind = iota
	NSub
	NMul
	NDiv
	NNumber
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

func newNumber(number int) *Node {
	return &Node{
		kind:   NNumber,
		number: number,
	}
}

// expr = mul() (+ mul() | - mul())*
func (o *parser) expr() *Node {
	node := o.mul()

	ok := false

	for {
		if o.token, ok = consume(o.token, '+'); ok {
			node = newBinary(NAdd, node, o.mul())
		}
		if o.token, ok = consume(o.token, '-'); ok {
			node = newBinary(NSub, node, o.mul())
		}
		return node
	}
}

// mul = primary() (* primary() | / primary())*
func (o *parser) mul() *Node {
	node := o.primary()

	ok := false

	for {
		if o.token, ok = consume(o.token, '*'); ok {
			node = newBinary(NMul, node, o.primary())
		}
		if o.token, ok = consume(o.token, '/'); ok {
			node = newBinary(NDiv, node, o.primary())
		}
		return node
	}
}

// primary = number | "(" expr ")"
func (o *parser) primary() *Node {
	ok := false
	if o.token, ok = consume(o.token, '('); ok {
		node := o.expr()
		expect(o.token, ')')
		return node
	}

	val := 0
	o.token, val = expectNumber(o.token)
	return newNumber(val)
}

func gen(node *Node) {
	if node.kind == NNumber {
		fmt.Printf("  push %d\n", node.number)
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
	code := arg()
	p := &parser{
		tokenize(code),
	}
	node := p.expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global _main\n")
	fmt.Printf("_main:\n")

	gen(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
