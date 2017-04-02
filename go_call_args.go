package gocoder

import (
	"go/ast"
)

type GoCallArgs struct {
	goCall *GoCall

	astArgs []ast.Expr
}

func newGoCallArgs(goCall *GoCall, args []ast.Expr) *GoCallArgs {
	return &GoCallArgs{
		goCall:  goCall,
		astArgs: args,
	}
}

func (p *GoCallArgs) GoCall() *GoCall {
	return p.goCall
}
