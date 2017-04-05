package gocoder

import (
	"go/ast"
)

type GoFunc struct {
	rootExpr *GoExpr

	decl      *ast.FuncDecl
	callExprs []*ast.CallExpr
}

func newGoFunc(rootExpr *GoExpr, decl *ast.FuncDecl) (gofunc *GoFunc) {
	g := &GoFunc{
		rootExpr: rootExpr,
		decl:     decl,
	}

	g.load()

	gofunc = g

	return
}

func (p *GoFunc) String() string {
	return p.decl.Name.String()
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
