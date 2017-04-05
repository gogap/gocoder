package gocoder

import (
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
	}

	g.load()

	return g
}

func (p *GoUnary) load() {
	p.GoExpr = newGoExpr(p.rootExpr, p.astExpr)
	if p.astExpr.X != nil {
		p.goChild = newGoExpr(p.rootExpr, p.astExpr.X)
	}
}

func (p *GoUnary) Inspect(f func(GoNode) bool) {
	if p.goChild != nil {
		p.goChild.Inspect(f)
	}
}

func (p *GoUnary) goNode() {}
