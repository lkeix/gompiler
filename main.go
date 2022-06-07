package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strconv"
)

// AT&T syntax
func emitOSExit() {
	fmt.Printf(".text\n")                  // start of text section
	fmt.Printf(".global _start\n")         // .global label: _start can be called from other files
	fmt.Printf("_start:\n")                // _start label: entry point
	fmt.Printf("  callq main.main\n\n")    // call main.main
	fmt.Printf("os.Exit:\n")               // os.Exit label: exit
	fmt.Printf("  movq 8(%%rsp), %%rdi\n") // rsp(stack pointer register) + 8 address  value(42(decimal) = 2a(hex)) to rdi(destination register)
	fmt.Printf("  movq $60, %%rax\n")      // rax(accumulator register) = 60
	fmt.Printf("  syscall\n\n")            // emit syscall
}

func emitSwap() {

}

func emitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		val := e.Value
		ival, err := strconv.Atoi(val)
		must(err)
		fmt.Printf("# %T\n", expr)
		fmt.Printf("  movq $%d, %%rax\n", ival)
		fmt.Printf("  pushq %%rax\n")
	case *ast.BinaryExpr:
		fmt.Printf("# start %T\n", expr)
		emitExpr(e.X)
		emitExpr(e.Y)
		fmt.Printf("  popq %%rbx\n")
		fmt.Printf("  popq %%rdi\n")
		switch e.Op.String() {
		case "+":
			fmt.Printf("  addq %%rdi, %%rax\n")
			fmt.Printf("  pushq %%rax\n")
		case "-":
			fmt.Printf("  subq %%rbx, %%rdi\n")
			fmt.Printf("  movq %%rdi, %%rax\n")
			fmt.Printf("  pushq %%rax\n")
		default:
			panic(fmt.Errorf("unexpected binary operator: %s", e.Op.String()))
		}
	default:
		must(fmt.Errorf("unexpected expr type %T", expr))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	source := "40 - 2"
	expr, err := parser.ParseExpr(source)
	must(err)

	emitOSExit()

	fmt.Printf(".text\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main.main:\n")
	emitExpr(expr)
	fmt.Printf("  popq %%rax\n")    // pop rax value and increament rs
	fmt.Printf("  pushq %%rax\n")   // decrement rsp and push rax value to rsp
	fmt.Printf("  callq os.Exit\n") // decrement (rsp address) - 8 and write next rip address(os.Exit) to rsp address
	fmt.Printf("  ret\n")
}
