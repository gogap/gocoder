package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoInterface struct {
	rootExpr *GoExpr
	astExpr  *ast.InterfaceType
}

func newGoInterface(rootExpr *GoExpr, astInterface *ast.InterfaceType) *GoInterface {
	g := &GoInterface{
		rootExpr: rootExpr,
		astExpr:  astInterface,
	}

	return g
}

func (p *GoInterface) Inspect(f InspectFunc, ctx context.Context) {
	return
}

func (p *GoInterface) goNode() {}

func (p *GoInterface) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoInterface) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}
