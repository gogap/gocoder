package gocoder

import (
	"context"
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
			case *ast.ValueSpec:
				{
					if expr.Type != nil {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Type))
					}
				}
			case *ast.Field:
				{
					if expr.Type != nil {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Type))
					}
				}
			}

		}
	}
}

func (p *GoIdent) Inspect(f InspectFunc, ctx context.Context) {
	for i := 0; i < len(p.goChildren); i++ {
		p.goChildren[i].Inspect(f, ctx)
	}
}

func (p *GoIdent) Name() string {
	return p.astExpr.Name
}

func (p *GoIdent) goNode() {}
