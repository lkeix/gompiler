package main

import (
	"flag"
	"fmt"
	"strconv"
	"unicode"
)

// tokenizer
type (
	TokenKind int

	Token struct {
		kind  TokenKind
		token string
		next  *Token
		val   int
	}
)

const (
	Reserved TokenKind = iota
	Number
	EOF
)

func newToken(kind TokenKind, current *Token, token string) *Token {
	next := &Token{
		kind:  kind,
		token: token,
	}
	current.next = next
	return next
}

func tokenize(code string) *Token {
	extractNumber := func(i int) int {
		for ; i < len(code); i++ {
			if !unicode.IsDigit(rune(code[i])) {
				break
			}
		}
		return i
	}

	head := &Token{}
	current := head

	for i := 0; i < len(code); i++ {
		// operator
		if code[i] == '+' ||
			code[i] == '-' ||
			code[i] == '*' ||
			code[i] == '/' ||
			code[i] == '(' ||
			code[i] == ')' {
			current = newToken(Reserved, current, string(code[i:i+1]))
			continue
		}

		// number
		if unicode.IsDigit(rune(code[i])) {
			j := extractNumber(i)
			val, _ := strconv.Atoi(string(code[i:j]))
			current = newToken(Number, current, string(code[i:j]))
			current.val = val
			i = j - 1
		}
		// relation
		if i >= len(code)-1 {
			continue
		}

		if code[i:i+2] == "==" ||
			code[i:i+2] == "!=" ||
			code[i:i+2] == "<=" ||
			code[i:i+2] == ">=" {
			current = newToken(Reserved, current, code[i:i+2])
			continue
		}

		if code[i] == '<' ||
			code[i] == '>' {
			current = newToken(Reserved, current, string(code[i]))
		}
	}

	current = newToken(EOF, current, "")
	return head.next
}

// parser
type (
	NodeKind int

	Node struct {
		kind   NodeKind
		lhs    *Node
		rhs    *Node
		number int
	}

	Parser struct {
		token *Token
	}
)

const (
	Add NodeKind = iota
	Sub
	Mul
	Div
	Eq
	Ne
	Ge
	Gt
	Le
	Lt
	Equality
	NNumber
)

func newNumberNode(number int) *Node {
	return &Node{
		kind:   NNumber,
		number: number,
	}
}

func newBinary(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func (p *Parser) consume(op string) bool {
	if p.token.kind != Reserved ||
		p.token.token != op {
		return false
	}

	p.token = p.token.next
	return true
}

func (p *Parser) expect(op string) {
	if p.token.kind != Reserved ||
		p.token.token != op {
		panic("Error: expected " + op + ". But " + p.token.token)
	}
	p.token = p.token.next
}

func (p *Parser) expectNumber() int {
	if p.token.kind != Number {
		panic("Error: expected number")
	}
	val := p.token.val
	p.token = p.token.next
	return val
}

// expr = Equality()
func (p *Parser) expr() *Node {
	return p.Equality()
}

// relation = add() (< add() | > add() | < add() | >= add())*
func (p *Parser) Relation() *Node {
	n := p.add()

	for {
		if p.consume("<") {
			n = newBinary(Lt, n, p.add())
		}

		if p.consume("<=") {
			n = newBinary(Le, n, p.add())
		}

		if p.consume(">") {
			n = newBinary(Lt, p.add(), n)
		}

		if p.consume(">=") {
			n = newBinary(Le, p.add(), n)
		}
		return n
	}
}

// Equality = relation() (== Relation() | != Relation())*
func (p *Parser) Equality() *Node {
	n := p.Relation()

	for {
		if p.consume("==") {
			n = newBinary(Eq, n, p.Relation())
		}

		if p.consume("!=") {
			n = newBinary(Ne, n, p.Relation())
		}

		return n
	}
}

// add = mul() (+ mul() | - mul())*
func (p *Parser) add() *Node {
	n := p.mul()

	for {
		if p.consume("+") {
			n = newBinary(Add, n, p.mul())
		}
		if p.consume("-") {
			n = newBinary(Add, n, p.mul())
		}
		return n
	}
}

// mul = unary() (* unary() | / unary())*
func (p *Parser) mul() *Node {
	n := p.unary()

	for {
		if p.consume("*") {
			n = newBinary(Mul, n, p.unary())
		}
		if p.consume("/") {
			n = newBinary(Div, n, p.unary())
		}
		return n
	}
}

// unary = (+ | -) ? primary()
func (p *Parser) unary() *Node {
	if p.consume("+") {
		return p.primary()
	}

	if p.consume("-") {
		return newBinary(Sub, newNumberNode(0), p.primary())
	}

	return p.primary()
}

// primary = number | ( expr )
func (p *Parser) primary() *Node {
	if p.consume("(") {
		n := p.expr()
		p.expect(")")
		return n
	}
	return newNumberNode(p.expectNumber())
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
	case Lt:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzx rax, al\n")
	case Le:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzx rax, al\n")
	case Ne:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzx rax, al\n")
	case Eq:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzx rax, al\n")
	case Add:
		fmt.Printf("  add rax, rdi\n")
	case Sub:
		fmt.Printf("  sub rax, rdi\n")
	case Mul:
		fmt.Printf("  imul rax, rdi\n")
	case Div:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rax, rdi\n")
	}

	fmt.Printf("  push rax\n")
}

// main
func main() {
	input := flag.String("input", "", "")
	flag.Parse()

	token := tokenize(*input)

	p := &Parser{
		token: token,
	}

	node := p.expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global _main\n")
	fmt.Printf("_main:\n")

	gen(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
