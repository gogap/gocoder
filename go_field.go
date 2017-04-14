package gocoder

import (
	"go/ast"
	"go/token"
	"strings"
)

type GoField struct {
	rootExpr *GoExpr

	astExpr *ast.Field

	// goExpr  *GoExpr
	goNames []*GoIdent

	fieldType *GoType
}

func newGoField(rootExpr *GoExpr, field *ast.Field) *GoField {
	g := &GoField{
		rootExpr: rootExpr,
		astExpr:  field,
	}

	g.load()

	return g
}

func (p *GoField) NumName() int {
	return len(p.astExpr.Names)
}

func (p *GoField) IsExported() bool {
	if len(p.astExpr.Names) == 0 {
		return true
	}

	return p.astExpr.Names[0].IsExported()
}

func (p *GoField) Name(i int) *GoIdent {
	return p.goNames[i]
}

func (p *GoField) Type() *GoType {
	return p.fieldType
}

func (p *GoField) Tag() StructTag {
	if p.astExpr.Tag == nil {
		return StructTag("")
	}

	tag := strings.Trim(p.astExpr.Tag.Value, "\"")
	tag = strings.Trim(tag, "`")

	return StructTag(tag)
}

func (p *GoField) load() {

	for i := 0; i < len(p.astExpr.Names); i++ {
		p.goNames = append(p.goNames, newGoIdent(p.rootExpr, p.astExpr.Names[i]))
	}

	p.fieldType = newGoType(p.rootExpr, p.astExpr, p.astExpr.Type)
}

func (p *GoField) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoField) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoField) goNode() {}
