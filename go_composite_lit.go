package gocoder

import (
	"go/ast"
)

type GoCompositeLit struct {
	*GoExpr
	rootExpr *GoExpr
	goChild  *GoExpr

	astExpr *ast.CompositeLit
}

func newGoCompositeLit(rootExpr *GoExpr, astCompositeLit *ast.CompositeLit) *GoCompositeLit {
	g := &GoCompositeLit{
		rootExpr: rootExpr,
		astExpr:  astCompositeLit,
		GoExpr:   newGoExpr(rootExpr, astCompositeLit),
	}

	g.load()

	return g
}

func (p *GoCompositeLit) load() {
	p.goChild = newGoExpr(p.rootExpr, p.astExpr.Type)
}

func (p *GoCompositeLit) Inspect(f func(GoNode) bool) {
	p.goChild.Inspect(f)
}

func (p *GoCompositeLit) goNode() {}
