package gocoder

import (
	"go/ast"
)

type GoBasicLit struct {
	*GoExpr

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
	p.GoExpr = newGoExpr(p.rootExpr, p.astExpr)
}

func (p *GoBasicLit) Inspect(f func(GoNode) bool) {
	return
}

func (p *GoBasicLit) Value() string {
	return p.astExpr.Value
}

func (p *GoBasicLit) Kind() string {
	return p.astExpr.Kind.String()
}

func (p *GoBasicLit) goNode() {}
