package pure

import (
	"go/ast"
	//"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	
	//"fmt"
	"strings"
	//"reflect"
)

const doc = "pure is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name:		"pure",
	Doc:		doc,
	Run:		run,
	Requires:	[]*analysis.Analyzer{
		inspect.Analyzer,
	},
	//FactTypes:	[]analysis.Fact{new(isWrapper)},
}
/*
type isWrapper struct {}
func (f *isWrapper) AFact() {}
func (f *isWrapper) String() string {
	return "FACT"
}
*/
func run(pass *analysis.Pass) (interface{}, error) {
	pureFuncDict := make(map[*ast.FuncDecl]bool)
	
	// collect all functions and their purity
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.FuncDecl:
				if pureAttributed(node) {
					pureFuncDict[node] = true
				}
			}
		}
	}
	
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	
	// inspect all function declarations
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			if pureFuncDict[node] {
				checkFuncPurity(pass, node, pureFuncDict)
			}
		}
	})
	
	return nil, nil
}

// is there a comment `@pure`
func pureAttributed(fd *ast.FuncDecl) bool {
	if fd.Doc == nil {
		return false
	}
	for _, comment := range fd.Doc.List {
		if strings.Index(comment.Text, "@pure") != -1 {
			return true
		}
	}
	return false
}

// look inside function body
func checkFuncPurity(pass *analysis.Pass, fd *ast.FuncDecl, dict map[*ast.FuncDecl]bool) bool {
	if fd == nil {
		return true
	}
	result := true
	
	// see all identifiers
	ast.Inspect(fd, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			result = result && checkIdent(pass, node, fd, dict)
		}
		return true
	})
	
	return result
}

// fd : the func-decl we are in
func checkIdent(pass *analysis.Pass, ident *ast.Ident, fd *ast.FuncDecl, dict map[*ast.FuncDecl]bool) bool {
	/*
	obj := pass.TypesInfo.Defs[ident]
	if obj == nil {
		fmt.Println(ident.String(), " : not found")
		return true
	}
	
	// reference to a function
	if _, ok := obj.Type().(*types.Signature); ok {
		if !dict[obj] {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot call impure function \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
		}
	} else {
		fmt.Println(obj.String(), " : ", obj.Type().String())
	}
	return true
	*/
	
	
	// look for the declaration of the identifier
	if ident.Obj == nil {
		if pass.TypesInfo.Types[ident].IsType() {
			return true
		} else {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m refers to a symbol \x1b[1m%s\x1b[0m which is not declared within the same package.\n", fd.Name.Name, ident.Name)
			return false
		}
	}
	
	switch decl := ident.Obj.Decl.(type) {
	
	case *ast.FuncDecl:
		// do not call checkFuncPurity; otherwise it can loop forever
		// call of impure function inside a pure function
		if dict[decl] {
			return true
		} else {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot call impure function \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
			return false
		}
	
	case *ast.ValueSpec:
		variable := decl.Names[0]
		pos := variable.NamePos
		// reference to the variable declared out of the function scope, and it is mutable
		if !(fd.Pos() < pos && pos < fd.End()) && ident.Obj.Kind != ast.Con {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot refer to a mutable variable \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
			return false
		} else {
			return true
		}
	
	case *ast.AssignStmt:
		pos := decl.Pos()
		// reference to the variable declared out of the function scope
		if !(fd.Pos() < pos && pos < fd.End()) {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot refer to a mutable variable \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
			return false
		} else {
			return true
		}
	
	default:
		return true
	}
	
}

