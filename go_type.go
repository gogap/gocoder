package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoType struct {
	rootExpr *GoExpr

	astExpr ast.Expr
	strType string

	node   GoNode
	parent ast.Node

	funcs []*GoFunc
}

func newGoType(rootExpr *GoExpr, parent ast.Node, astType ast.Expr) *GoType {
	g := &GoType{
		rootExpr: rootExpr,
		astExpr:  astType,
		parent:   parent,
		strType:  astTypeToStringType(astType),
	}

	g.load()

	return g
}

func (p *GoType) Node() GoNode {
	return p.node
}

func (p *GoType) load() {
	node := p.astExpr.(ast.Node)

	switch expr := node.(type) {
	case *ast.StructType:
		{
			spec, _ := p.parent.(*ast.TypeSpec)
			p.node = newGoStruct(p.rootExpr, spec, expr)
		}
	case *ast.Ident:
		{
			p.node = newGoIdent(p.rootExpr, expr)
		}
	case *ast.ArrayType:
		{
			p.node = newGoArray(p.rootExpr, expr)
		}
	case *ast.MapType:
		{
			p.node = newGoMap(p.rootExpr, expr)
		}
	case *ast.InterfaceType:
		{
			p.node = newGoInterface(p.rootExpr, expr)
		}
	case *ast.StarExpr:
		{
			p.node = newGoStar(p.rootExpr, expr)
		}
	case *ast.SelectorExpr:
		{
			p.node = newGoSelector(p.rootExpr, expr)
		}
	}

	typSpec, ok := p.parent.(*ast.TypeSpec)
	if ok {
		for i := 0; i < p.rootExpr.options.GoPackage.NumFuncs(); i++ {
			if p.rootExpr.options.GoPackage.Func(i).Receiver() == typSpec.Name.Name {
				p.funcs = append(p.funcs, p.rootExpr.options.GoPackage.Func(i))
			}
		}
	}
}

func (p *GoType) MethodByName(name string) *GoFunc {
	for i := 0; i < len(p.funcs); i++ {
		if p.funcs[i].Name() == name {
			return p.funcs[i]
		}
	}

	return nil
}

func (p *GoType) Method(i int) *GoFunc {
	return p.funcs[i]
}

func (p *GoType) NumMethods() int {
	return len(p.funcs)
}

func (p *GoType) Inspect(f InspectFunc, ctx context.Context) {
	inspectable, ok := p.node.(GoNodeInspectable)
	if ok {
		inspectable.Inspect(f, ctx)
	}
}

func (p *GoType) IsArray() bool {
	_, ok := p.astExpr.(*ast.ArrayType)
	return ok
}

func (p *GoType) IsInterface() bool {
	_, ok := p.astExpr.(*ast.InterfaceType)
	return ok
}

func (p *GoType) IsMap() bool {
	_, ok := p.astExpr.(*ast.MapType)
	return ok
}

func (p *GoType) IsStruct() bool {
	_, ok := p.astExpr.(*ast.StructType)
	return ok
}

func (p *GoType) IsSelector() bool {
	_, ok := p.astExpr.(*ast.SelectorExpr)
	return ok
}

func (p *GoType) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoType) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoType) String() string {
	return p.strType
}

func (p *GoType) goNode() {}
