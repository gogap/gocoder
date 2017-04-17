package gocoder

import (
	"go/ast"
	"go/token"
)

type GoStruct struct {
	rootExpr *GoExpr

	astExpr *ast.StructType

	goFileds []*GoField

	spec *ast.TypeSpec
}

func newGoStruct(rootExpr *GoExpr, spec *ast.TypeSpec, expr *ast.StructType, options ...Option) *GoStruct {

	g := &GoStruct{
		rootExpr: rootExpr,
		astExpr:  expr,
		spec:     spec,
	}

	g.load()

	return g
}

func (p *GoStruct) Name() string {
	return p.spec.Name.String()
}

func (p *GoStruct) NumFields() int {
	return p.astExpr.Fields.NumFields()
}

func (p *GoStruct) Field(i int) *GoField {
	return p.goFileds[i]
}

func (p *GoStruct) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoStruct) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoStruct) load() {

	var goFileds []*GoField

	for i := 0; i < len(p.astExpr.Fields.List); i++ {
		field := p.astExpr.Fields.List[i]
		goFileds = append(goFileds, newGoField(p.rootExpr, field))
	}

	p.goFileds = goFileds
}

func (p *GoStruct) goNode() {}
