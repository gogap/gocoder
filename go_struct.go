package gocoder

import (
	"go/ast"
	"go/token"
	"strings"
)

type GoStructField struct {
	Name       string
	Type       string
	Tag        StructTag
	IsPointer  bool
	IsArray    bool
	IsMap      bool
	IsCombined bool
	IsExported bool

	UsingPackage *GoPackage `json:"-"`
	*GoExpr
}

type GoStruct struct {
	rootExpr *GoExpr

	astExpr *ast.StructType
}

func newGoStruct(rootExpr *GoExpr, expr *ast.StructType, options ...Option) *GoStruct {

	return &GoStruct{
		rootExpr: rootExpr,
		astExpr:  expr,
	}
}

func (p *GoStruct) NumFields() int {
	return p.astExpr.Fields.NumFields()
}

func (p *GoStruct) Fields() []GoStructField {

	var goStructFields []GoStructField

	for i := 0; i < p.astExpr.Fields.NumFields(); i++ {

		field := p.astExpr.Fields.List[i]

		_, isPointer := field.Type.(*ast.StarExpr)

		structField := GoStructField{
			IsCombined: len(field.Names) == 0,
			IsExported: true,
			IsPointer:  isPointer,
			Type:       fieldTypeToStringType(field),
		}

		if field.Tag != nil {
			structField.Tag = StructTag(strings.Trim(field.Tag.Value, "\""))
		}

		if len(field.Names) > 0 {
			structField.Name = field.Names[0].Name
			structField.IsExported = field.Names[0].IsExported()
		}

		goStructFields = append(goStructFields, structField)
	}

	return goStructFields
}

func (p *GoStruct) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoStruct) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoStruct) goNode() {}
