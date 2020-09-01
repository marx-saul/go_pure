package pure

import (
	"go/ast"

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
	Name: "pure",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
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
	if ident.Obj == nil {
		fmt.Printf("%s : Obj is nil\n", ident.Name)
		return true
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


/*
func checkStmtPurity(stmt ast.Stmt) bool {
	if stmt == nil {
		return true
	}
	
	switch stmt := stmt.(type) {
	case *ast.AssignStmt:
		result := true
		for _, exp := range stmt.Lhs {
			result = result && checkExprPurity(exp)
		}
		for _, exp := range stmt.Rhs {
			result = result && checkExprPurity(exp)
		}
		return result
		
	case *ast.BadStmt:
		return false
	
	case *ast.BlockStmt:
		result := true
		for _, substmt := range stmt.List {
			result = result && checkStmtPurity(substmt)
		}
		return result
	
	case *ast.ExprStmt:
		return checkExprPurity(stmt.X)
	
	case *ast.ForStmt:
		return checkStmtPurity(stmt.Init) && checkExprPurity(stmt.Post) && checkStmtPurity(stmt.Post) && checkStmtPurity(stmt.Body)
	
	case *ast.GoStmt:
		return checkExprPurity(stmt.Call)
	
	case *ast.IfStmt:
		return checkStmtPurity(stmt.Init) && checkExprPurity(stmt.Cond) && checkStmtPurity(stmt.Body) && checkStmtPurity(stmt.Else)
	
	case *ast.IncDecStmt:
		return checkExprPurity(stmt.X)
	
	case *ast.LabelStmt:
		return check
	
	default:
		fmt.Printf("checkStmtPurity() : AST %s not known\n", reflect.TypeOf(stmt))
		return true
	}
}

func checkExprPurity(exp ast.Expr) bool {
	if exp == nil {
		return true
	}
	
	switch exp := exp.(type) {
	case *ast.BadExpr:
		return false
	
	case *ast.BinaryExpr:
		return checkExprPurity(exp.X) && checkExprPurity(exp.Y)
	
	case *ast.CallExpr:
		result := checkExprPurity(exp.Fun)
		for _, arg := range exp.Args {
			result = result && checkExprPurity(arg)
		}
		return result
	
	case *ast.IndexExpr:
		return checkExprPurity(exp.X) && checkExprPurity(exp.Index)
	
	case *ast.SliceExpr:
		return checkExprPurity(exp.X) && checkExprPurity(exp.Low) && checkExprPurity(exp.High) && checkExprPurity(exp.Max)
	
	case *ast.ParenExpr:
		return checkExprPurity(exp.X)
	
	case *ast.StarExpr:
		return checkExprPurity(exp.X)
	
	case *ast.UnaryExpr:
		return checkExprPurity(exp.X)
	
	case *ast.TypeAssertExpr:
		return checkExprPurity(exp.X)
		
	case *ast.Ident:
		return checkIdentPurity(*exp)
	
	default:
		fmt.Printf("checkExprPurity() : AST %s not known\n", reflect.TypeOf(exp))
		return true
	}
}

func checkIdentPurity(ident ast.Ident) bool {
	fmt.Printf("checkIdentPurity(%s)\n", ident.Name)
	fmt.Printf("\t%s\n", reflect.TypeOf(ident))
	return true
}
*/
