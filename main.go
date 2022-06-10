package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

const MAIN = "main"

var (
	stringLiterals []string
)

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

func print() {
	fmt.Printf("# print\n")
	fmt.Printf(".text\n")
	fmt.Printf("runtime.print:\n")
	fmt.Printf("  movq $2, %%rdi\n")        // set 2 to rdi (stderr)
	fmt.Printf("  movq 16(%%rsp), %%rsi\n") // set 16(rsp) to rsi (string)
	fmt.Printf("  movq 8(%%rsp), %%rdx\n")  // set 8(rsp) to rdx (length)
	fmt.Printf("  movq $1, %%rax\n")        // set 1 to rax (syscall number)
	fmt.Printf("  syscall\n")
	fmt.Printf("  ret\n")
}

func generate(file *ast.File) {
	emitSL()
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			emitDeclFunc(MAIN, decl.(*ast.FuncDecl))
		}
	}
}

// semanticAnalyze analyzes the syntax tree and returns an error if there is any problem.
// now semanticAnalyze extract string literals from the syntax tree
func semanticAnalyze(file *ast.File) {
	declsWalk(file.Decls)
}

func declsWalk(decls []ast.Decl) {
	for _, decl := range decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			funcDecl := decl.(*ast.FuncDecl)
			bodyWalk(funcDecl.Body.List)
		default:
			must(fmt.Errorf("unexpected declaration: %T", decl))
		}
	}
}

func bodyWalk(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.ExprStmt:
			expr := stmt.(*ast.ExprStmt).X
			walkExpr(&expr)
		default:
			must(fmt.Errorf("Unexpected stmt type"))
		}
	}
}

func walkExpr(expr *ast.Expr) {
	switch e := (*expr).(type) {
	case *ast.CallExpr:
		for _, arg := range e.Args {
			walkExpr(&arg)
		}
	case *ast.ParenExpr: // "(" or ")" expr
		walkExpr(&e.X)
	case *ast.BasicLit:
		parseBasicLit(e)
	case *ast.BinaryExpr:
		walkExpr(&e.X)
		walkExpr(&e.Y)
	default:
		must(fmt.Errorf("unexpected expr type %T", *expr))
	}
}

func parseBasicLit(expr *ast.BasicLit) {
	switch expr.Kind.String() {
	// TODO INT
	case "STRING":
		stringLiterals = append(stringLiterals, expr.Value)
	default:
		must(fmt.Errorf("unexpected basic literal type %T", expr))
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
func emitDeclFunc(pkg string, funcDecl *ast.FuncDecl) {
	fmt.Printf("# %T\n", funcDecl)
	fmt.Printf(".text\n")
	fmt.Printf("%s.%s:\n", pkg, funcDecl.Name)

	// emit assembly code for function body. parse {...}
	emitFuncBody(funcDecl.Body)
}

func emitFuncBody(body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch stmt.(type) {
		case *ast.ExprStmt:
			expr := stmt.(*ast.ExprStmt).X
			emitExpr(expr)
		default:
			must(fmt.Errorf("unexpected stmt type %T", stmt))
		}
	}
}

// emitSL assmbly string literals in .data section
func emitSL() {
	fmt.Printf(".data\n")
	for i, sl := range stringLiterals {
		fmt.Printf("S%d:\n", i)
		fmt.Printf("  %s\n", sl)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// define file set
	fset := token.NewFileSet()
	// parse source from source/main.go
	f, err := parser.ParseFile(fset, "./source/main.go", nil, parser.ParseComments)
	must(err)

	// generate assembly code
	generate(f)

	// define runtime and os.Exit
	runtime()
	osExit()
	print()
}
