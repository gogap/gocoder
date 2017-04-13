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

	funcs []*GoFunc
}

func newGoStruct(rootExpr *GoExpr, spec *ast.TypeSpec, expr *ast.StructType, options ...Option) *GoStruct {

	g := &GoStruct{
		rootExpr: rootExpr,
		astExpr:  expr,
		spec:     spec,
	}

	// fmt.Println(rootExpr.options)

	// g.rootExpr.options.Fallback(options...)

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

func (p *GoStruct) Method(name string) *GoFunc {
	for i := 0; i < len(p.funcs); i++ {
		if p.funcs[i].Name() == name {
			return p.funcs[i]
		}
	}

	return nil
}

func (p *GoStruct) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoStruct) load() {

	var goFileds []*GoField

	for i := 0; i < p.astExpr.Fields.NumFields(); i++ {
		field := p.astExpr.Fields.List[i]
		goFileds = append(goFileds, newGoField(p.rootExpr, field))
	}

	p.goFileds = goFileds

	for i := 0; i < p.rootExpr.options.GoPackage.NumFuncs(); i++ {
		if p.rootExpr.options.GoPackage.Func(i).Receiver() == p.Name() {
			p.funcs = append(p.funcs, p.rootExpr.options.GoPackage.Func(i))
		}
	}
}

func (p *GoStruct) goNode() {}
