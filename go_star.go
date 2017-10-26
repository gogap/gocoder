package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoStar struct {
	rootExpr *GoExpr
	astExpr  *ast.StarExpr
	goChild  *GoExpr
}

func newGoStar(rootExpr *GoExpr, astStar *ast.StarExpr) *GoStar {
	g := &GoStar{
		rootExpr: rootExpr,
		astExpr:  astStar,
		goChild:  newGoExpr(rootExpr, astStar.X),
	}

	return g
}

func (p *GoStar) X() *GoExpr {
	return p.goChild
}

func (p *GoStar) Inspect(f InspectFunc, ctx context.Context) {
	if p.goChild != nil {
		p.goChild.Inspect(f, ctx)
	}
}

func (p *GoStar) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoStar) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoStar) goNode() {}
