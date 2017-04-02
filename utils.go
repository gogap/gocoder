package gocoder

import (
	"fmt"
	"go/ast"
)

func fieldTypeToStringType(field *ast.Field) string {
	exprStr := ""

	selector := 0

	ast.Inspect(field.Type, func(n ast.Node) bool {
		switch ex := n.(type) {
		case *ast.SelectorExpr:
			{
				selector += 1
			}
		case *ast.Ident:
			exprStr += ex.Name
			if selector > 0 {
				exprStr += "."
				selector--
			}
		case *ast.InterfaceType:
			exprStr += "interface{}"
		case *ast.MapType:
			exprStr += fmt.Sprintf("map[%s]%s", ex.Key, ex.Value)
		case *ast.StarExpr:
			exprStr += "*"
		case *ast.ArrayType:
			exprStr += "[]"
		case *ast.Ellipsis:
			exprStr += "..."
		}

		return true
	})

	return exprStr
}

func isFieldTypeOf(field *ast.Field, strType string) bool {
	return fieldTypeToStringType(field) == strType
}

func isCallingFuncOf(expr interface{}, name string) bool {
	switch ex := (expr).(type) {
	case *ast.CallExpr:
		{
			switch fn := ex.Fun.(type) {
			case *ast.SelectorExpr:
				{
					return fn.Sel.Name == name
				}
			case *ast.Ident:
				{
					return fn.Name == name
				}
			}
		}
	}
	return false
}
