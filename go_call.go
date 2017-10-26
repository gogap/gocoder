package gocoder

import (
	"go/ast"
	"go/token"
)

type GoCall struct {
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
}

func (p *GoCall) Name() string {

	switch callFunc := p.astCall.Fun.(type) {
	case *ast.Ident:
		{
			return callFunc.Name
		}
	case *ast.SelectorExpr:
		{
			return callFunc.Sel.Name
		}
	}

	return ""
}

func (p *GoCall) Func() *GoExpr {
	return p.goFun
}

func (p *GoCall) Arg(i int) *GoExpr {
	return p.args[i]
}

func (p *GoCall) NumArgs() int {
	return len(p.args)
}

func (p *GoCall) goNode() {
}

func (p *GoCall) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astCall.Pos()), p.rootExpr.astFileSet.Position(p.astCall.End())
}

func (p *GoCall) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astCall)
}
