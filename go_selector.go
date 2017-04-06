package gocoder

import (
	"go/ast"
)

type GoSelector struct {
	*GoExpr
	rootExpr *GoExpr

	goSelIdent *GoIdent
	goXExpr    *GoExpr

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
	if p.astExpr.X != nil {
		p.goXExpr = newGoExpr(p.rootExpr, p.astExpr.X)
	}

	if p.astExpr.Sel != nil {
		p.goSelIdent = newGoIdent(p.rootExpr, p.astExpr.Sel)
	}
}

func (p *GoSelector) IsPackage() bool {
	xIdent, ok := p.astExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	if len(xIdent.Name) == 0 {
		return false
	}

	gofile := p.Root().Options().GoFile

	if gofile == nil {
		return false
	}

	_, b := gofile.FindImportByName(xIdent.Name)

	return b
}

func (p *GoSelector) ToPackage() *GoPackage {
	xIdent, ok := p.astExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}

	if len(xIdent.Name) == 0 {
		return nil
	}

	gofile := p.Root().Options().GoFile

	if gofile == nil {
		return nil
	}

	pkg, _ := gofile.FindImportByName(xIdent.Name)

	return pkg
}

func (p *GoSelector) Inspect(f func(GoNode) bool) {
	p.goXExpr.Inspect(f)
	p.goSelIdent.Inspect(f)
}

func (p *GoSelector) goNode() {}
