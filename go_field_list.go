package gocoder

import (
	"go/ast"
	"go/token"
)

type GoFieldList struct {
	rootExpr *GoExpr
	astExpr  *ast.FieldList
}

func newFieldList(rootExpr *GoExpr, astFieldList *ast.FieldList) *GoFieldList {
	return &GoFieldList{
		rootExpr: rootExpr,
		astExpr:  astFieldList,
	}
}

func (p *GoFieldList) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoFieldList) NumFields() int {
	return p.astExpr.NumFields()
}

func (p *GoFieldList) TypesAre(paramsType ...string) bool {

	if p.astExpr.NumFields() != len(paramsType) {
		return false
	}

	for i := 0; i < p.astExpr.NumFields(); i++ {
		if !isFieldTypeOf(p.astExpr.List[i], paramsType[i]) {
			return false
		}
	}

	return true
}

func (p *GoFieldList) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoFieldList) goNode() {}
