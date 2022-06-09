package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

const MAIN = "main"

// AT&T syntax
func osExit() {
	fmt.Printf(".text\n")
	fmt.Printf("os.Exit:\n")               // os.Exit label: exit
	fmt.Printf("  movq 8(%%rsp), %%rdi\n") // rsp(stack pointer register) + 8 address  value(42(decimal) = 2a(hex)) to rdi(destination register)
	fmt.Printf("  movq $60, %%rax\n")      // rax(accumulator register) = 60
	fmt.Printf("  syscall\n\n")            // emit syscall
}

func runtime() {
	fmt.Printf("# runtime\n")
	fmt.Printf(".global _start\n")
	fmt.Printf("_start:\n")
	fmt.Printf("  callq main.main\n\n")
}

func generate(file *ast.File) {
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			emitDeclFunc(MAIN, decl.(*ast.FuncDecl))
		}
	}
}

func emitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		emitFunc(e)
	case *ast.ParenExpr: // "(" or ")" expr
		emitExpr(e.X)
	case *ast.BasicLit:
		emitBasicLit(e)
	case *ast.BinaryExpr:
		emitBinaryExpr(e)
	default:
		must(fmt.Errorf("unexpected expr type %T", expr))
	}
}

func emitBasicLit(expr *ast.BasicLit) {
	val := expr.Value
	ival, err := strconv.Atoi(val)
	must(err)
	fmt.Printf("# %T\n", expr)
	fmt.Printf("  movq $%d, %%rax\n", ival)
	fmt.Printf("  pushq %%rax\n")
}

func emitBinaryExpr(expr *ast.BinaryExpr) {
	fmt.Printf("# start %T\n", expr)
	emitExpr(expr.X) // left
	emitExpr(expr.Y) // right
	fmt.Printf("  popq %%rdi # right\n")
	fmt.Printf("  popq %%rax # left\n")
	switch expr.Op.String() {
	case "+":
		fmt.Printf("  addq %%rdi, %%rax\n")
		fmt.Printf("  pushq %%rax\n")
	case "-":
		fmt.Printf("  subq %%rdi, %%rax\n")
		fmt.Printf("  pushq %%rax\n")
	case "*":
		fmt.Printf("  imulq %%rdi, %%rax\n")
		fmt.Printf("  pushq %%rax\n")
	default:
		panic(fmt.Errorf("unexpected binary operator: %s", expr.Op.String()))
	}
}

func emitFunc(expr *ast.CallExpr) {
	pkg := expr.Args[0]
	fun := expr.Fun
	fmt.Printf("# fun = %T\n", fun)
	switch fn := fun.(type) {
	case *ast.SelectorExpr:
		emitExpr(pkg)
		fmt.Printf("  popq %%rax\n")
		fmt.Printf("  pushq %%rax\n")
		symbol := fmt.Sprintf("%s.%s", fn.X, fn.Sel)
		fmt.Printf("  callq %s\n", symbol)
	}
}

// emitDeclFunc emits assembly code for a declarated function. parse func XXX(...) {...}
func emitDeclFunc(pkg string, DeclFunc *ast.FuncDecl) {
	fmt.Printf("# %T\n", DeclFunc)
	fmt.Printf(".text\n")
	fmt.Printf("%s.%s:\n", pkg, DeclFunc.Name)

	// emit assembly code for function body. parse {...}
	for _, stmt := range DeclFunc.Body.List {
		switch stmt.(type) {
		case *ast.ExprStmt:
			expr := stmt.(*ast.ExprStmt).X
			emitExpr(expr)
		default:
			must(fmt.Errorf("unexpected stmt type %T", stmt))
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// define runtime and os.Exit
	runtime()
	osExit()

	// define file set
	fset := token.NewFileSet()
	// parse source from source/main.go
	f, err := parser.ParseFile(fset, "./source/main.go", nil, parser.ParseComments)
	must(err)

	// generate assembly code
	generate(f)
}
