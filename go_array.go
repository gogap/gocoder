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
	}

	return g
}

func (p *GoArray) Name() string {
	return astTypeToStringType(p.astExpr)
}

func (p *GoArray) Inspect(f InspectFunc, ctx context.Context) {
	p.astEle.Inspect(f, ctx)
}

func (p *GoArray) Ele() *GoExpr {
	return p.astEle
}

func (p *GoArray) goNode() {}

func (p *GoArray) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoArray) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}
