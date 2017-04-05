package gocoder

import (
	"go/ast"
)

type GoSelector struct {
	*GoExpr
	rootExpr *GoExpr

	astExpr *ast.SelectorExpr
}

func newGoSelector(rootExpr *GoExpr, astSelector *ast.SelectorExpr) *GoSelector {
	g := &GoSelector{
		rootExpr: rootExpr,
		astExpr:  astSelector,
		GoExpr:   newGoExpr(rootExpr, astSelector),
	}

	g.load()

	return g
}

func (p *GoSelector) load() {
}

func (p *GoSelector) goNode() {}
