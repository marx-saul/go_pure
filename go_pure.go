package go_pure

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	
	"fmt"
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
	pureFuncDict := make(map[types.Object]bool)
	
	// collect all functions and their purity
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.FuncDecl:
				if pureAttributed(node) {
					pureFuncDict[pass.TypesInfo.Defs[node.Name]] = true
				}
			}
		}
	}
	
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		fmt.Println("Unexpected Error")
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	
	// inspect all function declarations
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			if pureFuncDict[pass.TypesInfo.Defs[node.Name]] {
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
		if strings.Contains(comment.Text, "@pure") {
			return true
		}
	}
	return false
}

// look inside function body
func checkFuncPurity(pass *analysis.Pass, fd *ast.FuncDecl, dict map[types.Object]bool) bool {
	if fd == nil {
		return true
	}
	result := true
	
	//fmt.Println(fd.Name.Name)
	// see all identifiers
	ast.Inspect(fd.Body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			result = checkIdent(pass, node, fd, dict) && result
		}
		return true
	})
	
	return result
}

// fd : the func-decl we are in
func checkIdent(pass *analysis.Pass, ident *ast.Ident, fd *ast.FuncDecl, dict map[types.Object]bool) bool {
	/*fmt.Println(ident, ":", ident.Name)
	fmt.Println("\t", ident.Obj)
	fmt.Println("\t", pass.TypesInfo.Defs[ident])
	fmt.Println("\t", pass.TypesInfo.Uses[ident])
	fmt.Println()
	*/
	// look for the declaration of the identifier
	if ident.Obj == nil {
		if pass.TypesInfo.Types[ident].IsType() {
			return true
		}
		
		use := pass.TypesInfo.Uses[ident]
		if use == nil {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m refers to a unresolvable symbol \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
			return false
		}
		
		// it is a function
		if _, ok := use.Type().(*types.Signature); ok {
			// it is not pure attributed
			if !dict[use]  {
				pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot call impure function \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
				return false
			}
		}
		
		return true
	}
	
	switch decl := ident.Obj.Decl.(type) {
	
	case *ast.FuncDecl:
		// do not call checkFuncPurity; otherwise it can loop forever
		// call of impure function inside a pure function
		if dict[pass.TypesInfo.Defs[decl.Name]] {
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
