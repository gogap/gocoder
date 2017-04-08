package gocoder

import (
	"context"
	"go/ast"
)

type GoUnary struct {
	*GoExpr
	rootExpr *GoExpr
	goChild  *GoExpr

	astExpr *ast.UnaryExpr
}

func newGoUnary(rootExpr *GoExpr, unary *ast.UnaryExpr) *GoUnary {
	g := &GoUnary{
		rootExpr: rootExpr,
		astExpr:  unary,
		GoExpr:   newGoExpr(rootExpr, unary),
	}

	g.load()

	return g
}

func (p *GoUnary) load() {
	if p.astExpr.X != nil {
		p.goChild = newGoExpr(p.rootExpr, p.astExpr.X)
	}
}

func (p *GoUnary) Inspect(f InspectFunc, ctx context.Context) {
	if p.goChild != nil {
		p.goChild.Inspect(f, ctx)
	}
}

func (p *GoUnary) goNode() {}
