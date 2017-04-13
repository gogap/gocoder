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
			p.node = newGoStruct(p.rootExpr, p.parent.(*ast.TypeSpec), expr)
		}
	case *ast.SelectorExpr:
		{
			p.node = newGoSelector(p.rootExpr, expr)
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
		// case *ast.TypeSpec:
		// 	{
		// 		fmt.Println("...")
		// 		switch n := expr.Type.(type) {
		// 		case *ast.StructType:
		// 			{
		// 				p.node = newGoStruct(p.rootExpr, expr, n)
		// 			}
		// 		}
		// 	}
	}
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

func (p *GoType) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoType) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoType) String() string {
	return p.strType
}

func (p *GoType) goNode() {}
