package gocoder

import (
	"go/ast"
	"go/token"
)

type GoFunc struct {
	rootExpr *GoExpr

	decl      *ast.FuncDecl
	callExprs []*ast.CallExpr

	assignments map[string]*GoAssignStmt
}

func newGoFunc(rootExpr *GoExpr, decl *ast.FuncDecl) (gofunc *GoFunc) {
	g := &GoFunc{
		rootExpr:    rootExpr,
		decl:        decl,
		assignments: make(map[string]*GoAssignStmt),
	}

	g.load()

	gofunc = g

	return
}

func (p *GoFunc) Name() string {
	return p.decl.Name.String()
}

func (p *GoFunc) Receiver() string {

	if p.decl.Recv == nil {
		return ""
	}

	if len(p.decl.Recv.List) == 0 {
		return ""
	}

	nextExpr := p.decl.Recv.List[0].Type

nextCase:
	switch vType := nextExpr.(type) {
	case *ast.Ident:
		{
			return vType.Name
		}
	case *ast.StarExpr:
		{
			nextExpr = vType.X
			goto nextCase
		}
	}

	return ""
}

func (p *GoFunc) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.decl)
}

func (p *GoFunc) Root() *GoExpr {
	return p.rootExpr
}

func (p *GoFunc) Params() *GoFieldList {
	return newFieldList(p.rootExpr, p.decl.Type.Params)
}

func (p *GoFunc) Results() *GoFieldList {
	return newFieldList(p.rootExpr, p.decl.Type.Results)
}

func (p *GoFunc) FindCall(funcName string) (call *GoCall, exist bool) {

	for i := 0; i < len(p.callExprs); i++ {
		if isCallingFuncOf(p.callExprs[i], funcName) {
			goCall := newGoCall(p.rootExpr, p.callExprs[i])
			return goCall, true
		}
	}

	return nil, false
}

func (p *GoFunc) FindAssigment(identName string) (assignStmt *GoAssignStmt, exist bool) {

	assignStmt, exist = p.assignments[identName]

	if exist {
		return
	}

	ast.Inspect(p.decl, func(n ast.Node) bool {
		if exist {
			return false
		}

		switch node := n.(type) {
		case *ast.AssignStmt:
			{
				if len(node.Lhs) == 0 {
					return true
				}

				ident, ok := node.Lhs[0].(*ast.Ident)

				if !ok {
					return true
				}

				if identName != ident.Name {
					return true
				}

				assignStmt = newGoAssignStmt(p.rootExpr, node)
				exist = true

				p.assignments[ident.Name] = assignStmt

				return false
			}
		}

		return true
	})
	return
}

func (p *GoFunc) GetReturnStmt() *GoReturnStmt {

	var stmt *ast.ReturnStmt

	ast.Inspect(p.decl, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ReturnStmt:
			{
				stmt = node
			}
		}
		return true
	})

	if stmt != nil {
		return newGoReturnStmt(p.rootExpr, stmt)
	}

	return nil
}

func (p *GoFunc) load() {
	ast.Inspect(p.decl, func(n ast.Node) bool {
		switch exprType := n.(type) {
		case *ast.CallExpr:
			{
				p.callExprs = append(p.callExprs, exprType)
			}
		}
		return true
	})
}

func (p *GoFunc) goNode() {
}

func (p *GoFunc) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.decl.Pos())
}
