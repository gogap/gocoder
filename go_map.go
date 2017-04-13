package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoMap struct {
	rootExpr *GoExpr
	astExpr  *ast.MapType

	key   *GoExpr
	value *GoExpr
}

func newGoMap(rootExpr *GoExpr, astMap *ast.MapType) *GoMap {
	g := &GoMap{
		rootExpr: rootExpr,
		astExpr:  astMap,
		key:      newGoExpr(rootExpr, astMap.Key),
		value:    newGoExpr(rootExpr, astMap.Value),
	}

	return g
}

func (p *GoMap) Inspect(f InspectFunc, ctx context.Context) {
	return
}

func (p *GoMap) Value() *GoExpr {
	return p.value
}

func (p *GoMap) Key() *GoExpr {
	return p.key
}

func (p *GoMap) goNode() {}

func (p *GoMap) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoMap) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}
