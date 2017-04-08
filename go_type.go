package gocoder

import (
	"go/ast"
)

type GoType struct {
	*GoExpr
	rootExpr *GoExpr

	astExpr *ast.TypeSpec
}

func newGoType(rootExpr *GoExpr, astType *ast.TypeSpec) *GoType {
	g := &GoType{
		rootExpr: rootExpr,
		astExpr:  astType,
		GoExpr:   newGoExpr(rootExpr, astType.Type),
	}

	return g
}

func (p *GoType) IsArray() bool {
	_, ok := p.astExpr.Type.(*ast.ArrayType)
	return ok
}

func (p *GoType) IsInterface() bool {
	_, ok := p.astExpr.Type.(*ast.InterfaceType)
	return ok
}

func (p *GoType) IsMap() bool {
	_, ok := p.astExpr.Type.(*ast.MapType)
	return ok
}

func (p *GoType) IsStruct() bool {
	_, ok := p.astExpr.Type.(*ast.StructType)
	return ok
}

func (p *GoType) goNode() {}
