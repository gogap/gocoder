package gocoder

import (
	"go/ast"
)

type GoCall struct {
	goFunc *GoFunc

	astCall *ast.CallExpr
	args    *GoCallArgs
}

func newGoCall(goFunc *GoFunc, astCall *ast.CallExpr) *GoCall {
	g := &GoCall{
		goFunc:  goFunc,
		astCall: astCall,
	}

	g.load()

	return g
}

func (p *GoCall) GoFunc() *GoFunc {
	return p.goFunc
}

func (p *GoCall) load() {

}

func (p *GoCall) NumArgs() int {
	return len(p.astCall.Args)
}

func (p *GoCall) Args() *GoCallArgs {
	return p.args
}
