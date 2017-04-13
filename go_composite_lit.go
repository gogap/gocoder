package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoCompositeLit struct {
	// *GoExpr
	rootExpr *GoExpr
	goChild  *GoExpr

	astExpr *ast.CompositeLit
}

func newGoCompositeLit(rootExpr *GoExpr, astCompositeLit *ast.CompositeLit) *GoCompositeLit {
	g := &GoCompositeLit{
		rootExpr: rootExpr,
		astExpr:  astCompositeLit,
		// GoExpr:   newGoExpr(rootExpr, astCompositeLit),
	}

	g.load()

	return g
}

func (p *GoCompositeLit) load() {
	p.goChild = newGoExpr(p.rootExpr, p.astExpr.Type)
}

func (p *GoCompositeLit) Inspect(f InspectFunc, ctx context.Context) {
	p.goChild.Inspect(f, ctx)
}

func (p *GoCompositeLit) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoCompositeLit) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoCompositeLit) goNode() {}
