package gocoder

import (
	"go/ast"
)

type GoCall struct {
	*GoExpr

	rootExpr *GoExpr
	args     []*GoExpr
	goFun    *GoExpr

	astCall *ast.CallExpr
}

func newGoCall(rootExpr *GoExpr, astCall *ast.CallExpr) *GoCall {
	g := &GoCall{
		rootExpr: rootExpr,
		astCall:  astCall,
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

	p.GoExpr = newGoExpr(p.rootExpr, p.astCall)
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

func (p *GoCall) Inspect(f func(GoNode) bool) {
	p.goFun.Inspect(f)
	for i := 0; i < len(p.args); i++ {
		p.args[i].Inspect(f)
	}
}

func (p *GoCall) goNode() {
}
