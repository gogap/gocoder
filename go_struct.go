package gocoder

import (
	"go/ast"
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
	*GoExpr

	rootExpr *GoExpr

	astExpr *ast.StructType
}

func newGoStruct(rootExpr *GoExpr, expr *ast.StructType, options ...Option) *GoStruct {

	return &GoStruct{
		GoExpr:  newGoExpr(rootExpr, expr, options...),
		astExpr: expr,
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
			Tag:        StructTag(strings.Trim(field.Tag.Value, "\"")),
			IsCombined: len(field.Names) == 0,
			IsExported: true,
			IsPointer:  isPointer,
			Type:       fieldTypeToStringType(field),
			GoExpr:     newGoExpr(p.rootExpr, field.Type),
		}

		if len(field.Names) > 0 {
			structField.Name = field.Names[0].Name
			structField.IsExported = field.Names[0].IsExported()
		}

		goStructFields = append(goStructFields, structField)
	}

	return goStructFields
}

// func (p *GoStruct) parseType(typ ast.Expr, preTyp ast.Expr, structField *GoStructField) {
// 	switch item := typ.(type) {
// 	case *ast.Ident:
// 		{
// 			structField.Types = append(structField.Types, item.Name)
// 		}
// 	case *ast.ArrayType:
// 		{
// 			structField.IsArray = true
// 			p.parseType(item.Elt, item, structField)
// 		}
// 	case *ast.SelectorExpr:
// 		{
// 			ident, ok := item.X.(*ast.Ident)
// 			if ok {
// 				structField.Types = append(structField.Types, ident.Name+"."+item.Sel.Name)
// 			}

// 			pkg, find := p.options.GoFile.FindImportByName(ident.Name)
// 			if find {
// 				structField.UsingPackage = pkg
// 			}
// 		}
// 	case *ast.StarExpr:
// 		{
// 			p.parseType(item.X, item, structField)
// 		}
// 	case *ast.MapType:
// 		{
// 			structField.IsMap = true
// 			p.parseType(item.Key, item, structField)
// 			p.parseType(item.Value, item, structField)

// 		}
// 	case *ast.InterfaceType:
// 		{
// 			structField.Types = append(structField.Types, "interface{}")
// 		}
// 	default:
// 		structField.Types = append(structField.Types, "<parse error>: "+reflect.TypeOf(item).String())
// 	}
// 	return
// }
