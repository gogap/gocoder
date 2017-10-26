package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoSelector struct {
	rootExpr *GoExpr

	goSelIdent *GoIdent
	goXExpr    *GoExpr

	astExpr *ast.SelectorExpr
}

func newGoSelector(rootExpr *GoExpr, astSelector *ast.SelectorExpr) *GoSelector {
	g := &GoSelector{
		rootExpr: rootExpr,
		astExpr:  astSelector,
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

func (p *GoSelector) X() *GoExpr {
	return p.goXExpr
}

func (p *GoSelector) IsInOtherPackage() bool {
	xIdent, ok := p.astExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	// if xIdent.Obj != nil {
	// 	return false
	// }

	if len(xIdent.Name) == 0 {
		return false
	}

	gofile := p.rootExpr.Options().GoFile

	if gofile == nil {
		return false
	}

	_, b := gofile.FindImportByName(xIdent.Name)

	return b
}

func (p *GoSelector) UsingPackage() *GoPackage {

	gofile := p.rootExpr.Options().GoFile

	xIdent, ok := p.astExpr.X.(*ast.Ident)
	if !ok {
		return gofile.Package()
	}

	// if xIdent.Obj != nil {
	// 	return gofile.Package()
	// }

	if len(xIdent.Name) == 0 {
		return gofile.Package()
	}

	if gofile == nil {
		return nil
	}

	pkg, exist := gofile.FindImportByName(xIdent.Name)
	if exist {
		return pkg
	}

	return gofile.Package()
}

func (p *GoSelector) GetSelName() string {
	return p.astExpr.Sel.Name
}

func (p *GoSelector) Inspect(f InspectFunc, ctx context.Context) {
	p.goXExpr.Inspect(f, ctx)
	p.goSelIdent.Inspect(f, ctx)
}

func (p *GoSelector) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoSelector) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoSelector) goNode() {}
