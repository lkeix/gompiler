package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

const MAIN = "main"

type (
	stringLiteral struct {
		tag   string
		value string
	}
	globalVariable struct {
		tag   string
		value string
		typ   *ast.Ident
	}
)

var (
	stringLiterals  []stringLiteral
	globalVariables []globalVariable

	globalString = &ast.Object{
		Kind: ast.Typ,
		Name: "string",
		Decl: nil,
		Data: nil,
		Type: nil,
	}

	globalInt = &ast.Object{
		Kind: ast.Typ,
		Name: "int",
		Decl: nil,
		Data: nil,
		Type: nil,
	}
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
	fmt.Printf("  ret\n\n")
}

func declsWalk(decls []ast.Decl) {
	for _, decl := range decls {
		switch decl.(type) {
		case *ast.GenDecl:
			// extract global variables before analyze declaration functions
			parseGlobalVariables(decl.(*ast.GenDecl))
		case *ast.FuncDecl:
			funcDecl := decl.(*ast.FuncDecl)
			bodyWalk(funcDecl.Body.List)
		default:
			must(fmt.Errorf("unexpected declaration: %T", decl))
		}
	}
}

func declWalk(decl *ast.Decl) {
	switch (*decl).(type) {
	case *ast.GenDecl:
		// extract global variables before analyze declaration functions
		parseGlobalVariables((*decl).(*ast.GenDecl))
	case *ast.FuncDecl:
		funcDecl := (*decl).(*ast.FuncDecl)
		bodyWalk(funcDecl.Body.List)
	default:
		must(fmt.Errorf("unexpected declaration: %T", decl))
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
		parseStringLiteral(e)
	case *ast.BinaryExpr:
		walkExpr(&e.X)
		walkExpr(&e.Y)
	default:
		must(fmt.Errorf("unexpected expr type %T", *expr))
	}
}

func parseStringLiteral(expr *ast.BasicLit) {
	switch expr.Kind.String() {
	// TODO INT
	case "INT":
		break
	case "STRING":
		stringLiterals = append(stringLiterals, stringLiteral{tag: "", value: expr.Value})
	default:
		must(fmt.Errorf("unexpected basic literal type %T", expr))
	}
}

func parseGlobalVariables(decl *ast.GenDecl) {
	switch decl.Tok {
	case token.VAR:
		valSpec, ok := decl.Specs[0].(*ast.ValueSpec)
		if !ok {
			must(fmt.Errorf("unexpected value spec type %T", decl.Specs[0]))
		}
		fmt.Printf("# spec.Name=%v, spec.Value=%v\n", valSpec.Names[0], valSpec.Values[0])
		parseGrobalVariable(valSpec)
		typeIdent, ok := valSpec.Type.(*ast.Ident)
		if !ok {
			must(fmt.Errorf("unexpected type %T", valSpec.Type))
		}
		globalVariables = append(globalVariables, globalVariable{tag: valSpec.Names[0].Name, value: valSpec.Values[0].(*ast.BasicLit).Value, typ: typeIdent})
	}
}

func parseGrobalVariable(valSpec *ast.ValueSpec) {
	typeIdent, ok := valSpec.Type.(*ast.Ident)
	if !ok {
		must(fmt.Errorf("unexpected type ident %v", typeIdent))
	}
	switch typeIdent.Obj {
	case globalInt:
		_, ok := valSpec.Values[0].(*ast.BasicLit)
		if !ok {
			must(fmt.Errorf("unexpected type ident %v", typeIdent))
		}
	case globalString:
		lit, ok := valSpec.Values[0].(*ast.BasicLit)
		if !ok {
			must(fmt.Errorf("unexpected type ident %v", typeIdent))
		}
		parseStringLiteral(lit)
	default:
		must(fmt.Errorf("Unexpected global ident"))
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
	if expr.Kind.String() == "INT" {
		val := expr.Value
		ival, err := strconv.Atoi(val)
		must(err)
		fmt.Printf("# %T\n", expr)
		fmt.Printf("  movq $%d, %%rax\n", ival)
		fmt.Printf("  pushq %%rax\n")
	} else if expr.Kind.String() == "STRING" {
		fmt.Printf("  leaq %s, %%rax\n", stringLiterals[0].tag)
		fmt.Printf("  pushq %%rax\n")
		fmt.Printf("  pushq $%d\n", len(stringLiterals[0].value)-1-2)
		stringLiterals = stringLiterals[1:]
	} else {
		must(fmt.Errorf("unexpected basic literal type %T", expr))
	}
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
	fun := expr.Fun
	fmt.Printf("# fun = %T\n", fun)
	switch fn := fun.(type) {
	case *ast.Ident:
		if fn.Name == "print" {
			// build in print
			emitExpr(expr.Args[0]) // push string pointer, push string len
			fmt.Printf("  call runtime.print\n")
			fmt.Printf("  addq $8, %%rsp\n")
		}
	case *ast.SelectorExpr:
		emitExpr(expr.Args[0])
		fmt.Printf("  popq %%rax\n")
		fmt.Printf("  pushq %%rax\n")
		symbol := fmt.Sprintf("%s.%s", fn.X, fn.Sel)
		fmt.Printf("  callq %s\n\n", symbol)
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

func emitGlobalVariables() {
	for _, valSpec := range globalVariables {
		tag := valSpec.tag
		value := valSpec.value
		ident := valSpec.typ
		if ident.Obj == globalString {
			fmt.Printf("%s:\n", tag)
			// FIXME: searchTag time computational complexity is O(n) where n is the number of string literals.
			fmt.Printf("  .quad %s\n", searchTag(value))
			fmt.Printf("  .quad %d\n", len(value)-1-2)
		} else if ident.Obj == globalInt {
			fmt.Printf("%s:\n", tag)
			fmt.Printf("  .quad %s\n", value)
		} else {
			must(fmt.Errorf("unexpected type ident %v", ident))
		}
	}
	fmt.Printf("\n")
}

// emitSL assmbly string literals in .data section
func emitSL() {
	fmt.Printf(".data\n")
	for i, sl := range stringLiterals {
		fmt.Printf(".S%d:\n", i)
		fmt.Printf("  .string %s\n", sl.value)
		stringLiterals[i].tag = fmt.Sprintf(".S%d", i)
	}
	fmt.Printf("\n")
}

func searchTag(value string) string {
	for _, sl := range stringLiterals {
		if sl.value == value {
			return sl.tag
		}
	}
	must(fmt.Errorf("unexpected string literal value %s", value))
	return ""
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func setup(fset *token.FileSet, file *ast.File) {
	// setup universe block
	// detail on https://motemen.github.io/go-for-go-book/#%E3%82%B9%E3%82%B3%E3%83%BC%E3%83%97
	universe := &ast.Scope{
		Outer:   nil,
		Objects: make(map[string]*ast.Object),
	}

	universe.Insert(globalInt)
	universe.Insert(globalString)
	// insert build-in print function into universe block
	universe.Insert(&ast.Object{
		Kind: ast.Fun,
		Name: "print",
		Decl: nil,
		Data: nil,
		Type: nil,
	})

	universe.Insert(&ast.Object{
		Kind: ast.Pkg,
		Name: "os", // why ???
		Decl: nil,
		Data: nil,
		Type: nil,
	})

	ap, _ := ast.NewPackage(fset, map[string]*ast.File{"": file}, nil, universe)

	var unresolved []*ast.Ident
	for _, ident := range file.Unresolved {
		if obj := universe.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
		} else {
			unresolved = append(unresolved, ident)
		}
	}

	fmt.Printf("# Package:   %s\n", ap.Name)
}

// semanticAnalyze analyzes the syntax tree and returns an error if there is any problem.
// now semanticAnalyze extract string literals from the syntax tree
func semanticAnalyze(file *ast.File) {

	fmt.Printf("# global variables\n")
	for _, decl := range file.Decls {
		declWalk(&decl)
	}
	// declsWalk(file.Decls)
}

func generate(file *ast.File) {
	// emit string literals
	emitSL()

	// emit global variables
	fmt.Printf("# global variables\n")
	emitGlobalVariables()

	// emit declaration functions
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			emitDeclFunc(MAIN, decl.(*ast.FuncDecl))
		}
	}
}

func main() {
	// define file set
	fset := token.NewFileSet()
	// parse source from source/main.go
	f, err := parser.ParseFile(fset, "./source/main.go", nil, parser.ParseComments)
	must(err)

	// setup
	setup(fset, f)

	// semantic Analyze
	semanticAnalyze(f)
	// generate assembly code
	generate(f)

	// define runtime and os.Exit
	runtime()
	osExit()
	print()
}
