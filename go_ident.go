package gocoder

import (
	"go/ast"
)

type GoIdent struct {
	*GoExpr

	rootExpr   *GoExpr
	goChildren []*GoExpr

	astExpr *ast.Ident
}

func newGoIdent(rootExpr *GoExpr, ident *ast.Ident) *GoIdent {
	g := &GoIdent{
		rootExpr: rootExpr,
		astExpr:  ident,
		GoExpr:   newGoExpr(rootExpr, ident),
	}

	g.load()

	return g
}

func (p *GoIdent) load() {

	if p.astExpr.Obj == nil {
		return
	}

	switch p.astExpr.Obj.Kind {
	case ast.Var:
		{

			switch expr := p.astExpr.Obj.Decl.(type) {
			case *ast.AssignStmt:
				{
					for i := 0; i < len(expr.Rhs); i++ {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Rhs[i]))
					}
				}
			}

		}
	}
}

func (p *GoIdent) Inspect(f func(GoNode) bool) {
	if len(p.goChildren) > 0 {
		for i := 0; i < len(p.goChildren); i++ {
			p.goChildren[i].Inspect(f)
		}
	}
}

func (p *GoIdent) Name() string {
	return p.astExpr.Name
}

func (p *GoIdent) goNode() {}
