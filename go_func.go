package gocoder

import (
	"go/ast"
)

type GoFunc struct {
	goFile *GoFile
	decl   *ast.FuncDecl

	callExprs []*ast.CallExpr
}

func newGoFunc(gofile *GoFile, decl *ast.FuncDecl) (gofunc *GoFunc) {
	g := &GoFunc{
		goFile: gofile,
		decl:   decl,
	}

	g.load()

	gofunc = g

	return
}

func (p *GoFunc) String() string {
	return p.decl.Name.String()
}

func (p *GoFunc) GoFile() *GoFile {
	return p.goFile
}

func (p *GoFunc) Params() *GoParams {
	return &GoParams{
		goFunc:    p,
		astParams: p.decl.Type.Params,
	}
}

func (p *GoFunc) Results() *GoResults {
	return &GoResults{
		goFunc:     p,
		astResults: p.decl.Type.Results,
	}
}

func (p *GoFunc) FindCall(funcName string) (call *GoCall, exist bool) {

	for i := 0; i < len(p.callExprs); i++ {
		if isCallingFuncOf(p.callExprs[i], funcName) {
			goCall := newGoCall(p, p.callExprs[i])
			return goCall, true
		}
	}

	return nil, false
}

func (p *GoFunc) load() {
	ast.Inspect(p.decl, func(n ast.Node) bool {
		switch exprType := n.(type) {
		case *ast.CallExpr:
			{
				p.callExprs = append(p.callExprs, exprType)
			}
		}
		return true
	})
}
