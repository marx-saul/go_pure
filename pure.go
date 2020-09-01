package pure

import (
	"go/ast"

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
	FactTypes:	[]analysis.Fact{new(isWrapper)},
}

type isWrapper struct {}
func (f *isWrapper) AFact() {}
func (f *isWrapper) String() string {
	return "FACT"
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	
	// inspect all function declarations
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			if pureAttributed(node) {
				checkFuncPurity(pass, node)
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
func checkFuncPurity(pass *analysis.Pass, fd *ast.FuncDecl) bool {
	if fd == nil {
		return true
	}
	result := true
	
	// see all identifiers
	ast.Inspect(fd, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			result = result && checkIdent(pass, node, fd)
		}
		return true
	})
	
	return result
}

// fd : the func-decl we are in
func checkIdent(pass *analysis.Pass, ident *ast.Ident, fd *ast.FuncDecl) bool {
	// look for the declaration of the identifier
	if ident.Obj == nil {
		if pass.TypesInfo.Types[ident].IsType() {
			return true
		} else {
			pass.Reportf(ident.NamePos, "\x1b[1m%s\x1b[0m was not found in this module. Cross-module-purity-check has not been implemented yet.\n", ident.Name)
			return false
		}
	}
	
	switch decl := ident.Obj.Decl.(type) {
	
	case *ast.FuncDecl:
		// do not call checkFuncPurity; otherwise it can loop forever
		result := pureAttributed(decl)
		// call of impure function inside a pure function
		if !result {
			pass.Reportf(ident.NamePos, "Pure function \x1b[1m%s\x1b[0m cannot call impure function \x1b[1m%s\x1b[0m\n", fd.Name.Name, ident.Name)
		}
		return result
	
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

