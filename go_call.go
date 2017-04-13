package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoCall struct {
	// *GoExpr

	rootExpr *GoExpr
	args     []*GoExpr
	goFun    *GoExpr

	astCall *ast.CallExpr
}

func newGoCall(rootExpr *GoExpr, astCall *ast.CallExpr) *GoCall {
	g := &GoCall{
		rootExpr: rootExpr,
		astCall:  astCall,
		// GoExpr:   newGoExpr(rootExpr, astCall),
	}

	g.load()

	return g
}

func (p *GoCall) load() {

	p.goFun = newGoExpr(p.rootExpr, p.astCall.Fun)

	for _, arg := range p.astCall.Args {

		newArg := newGoExpr(
			p.rootExpr,
			arg)

		p.args = append(p.args, newArg)
	}
}

func (p *GoCall) Name() string {
	ident, ok := p.astCall.Fun.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func (p *GoCall) Args() []*GoExpr {
	return p.args
}

func (p *GoCall) Inspect(f InspectFunc, ctx context.Context) {
	p.goFun.Inspect(f, ctx)
	for i := 0; i < len(p.args); i++ {
		p.args[i].Inspect(f, ctx)
	}
}

func (p *GoCall) goNode() {
}

func (p *GoCall) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astCall.Pos())
}

func (p *GoCall) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astCall)
}
