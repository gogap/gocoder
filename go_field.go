package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoField struct {
	rootExpr *GoExpr

	astExpr *ast.Field

	goType  *GoExpr
	goNames []*GoIdent
}

func newGoField(rootExpr *GoExpr, field *ast.Field) *GoField {
	g := &GoField{
		rootExpr: rootExpr,
		astExpr:  field,
	}

	g.load()

	return g
}

func (p *GoField) Names() []*GoIdent {
	return p.goNames
}

func (p *GoField) Type() *GoExpr {
	return p.goType
}

func (p *GoField) load() {

	if p.astExpr.Type != nil {
		p.goType = newGoExpr(p.rootExpr, p.astExpr.Type)
	}

	for i := 0; i < len(p.astExpr.Names); i++ {
		p.goNames = append(p.goNames, newGoIdent(p.rootExpr, p.astExpr.Names[i]))
	}
}

func (p *GoField) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoField) Inspect(f InspectFunc, ctx context.Context) {
	p.goType.Inspect(f, ctx)
}

func (p *GoField) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoField) goNode() {}
