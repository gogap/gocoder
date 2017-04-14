package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoBasicLit struct {
	rootExpr *GoExpr

	astExpr *ast.BasicLit
}

func newGoBasicLit(rootExpr *GoExpr, astBasicLit *ast.BasicLit) *GoBasicLit {
	g := &GoBasicLit{
		rootExpr: rootExpr,
		astExpr:  astBasicLit,
	}

	g.load()

	return g
}

func (p *GoBasicLit) load() {
}

func (p *GoBasicLit) Inspect(f InspectFunc, ctx context.Context) {
	return
}

func (p *GoBasicLit) Value() string {
	return p.astExpr.Value
}

func (p *GoBasicLit) Kind() string {
	return p.astExpr.Kind.String()
}

func (p *GoBasicLit) goNode() {}

func (p *GoBasicLit) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoBasicLit) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}
