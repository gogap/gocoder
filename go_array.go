package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoArray struct {
	rootExpr *GoExpr

	astExpr *ast.ArrayType
	astEle  *GoExpr
}

func newGoArray(rootExpr *GoExpr, astExpr *ast.ArrayType) *GoArray {
	g := &GoArray{
		rootExpr: rootExpr,
		astExpr:  astExpr,
		astEle:   newGoExpr(rootExpr, astExpr.Elt),
		// GoExpr:   newGoExpr(rootExpr, astExpr),
	}

	return g
}

func (p *GoArray) Inspect(f InspectFunc, ctx context.Context) {
	p.astEle.Inspect(f, ctx)
}

func (p *GoArray) Ele() *GoExpr {
	return p.astEle
}

func (p *GoArray) goNode() {}

func (p *GoArray) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoArray) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}
