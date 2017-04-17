package gocoder

import (
	"context"
	"fmt"
	"strings"
)

type AnyOf struct {
	Type []string `json:"type"`
}

type JsonSchema struct {
	Schema     string                 `json:"$schema,omitempty"`
	Id         string                 `json:"id,omitempty"`
	Type       string                 `json:"type,omitempty"`
	AnyOf      *AnyOf                 `json:"anyOf,omitempty"`
	Items      *JsonSchema            `json:"items,omitempty"`
	Properties map[string]*JsonSchema `json:"properties,omitempty"`
}

func NodeToJsonSchema(goNode GoNode) (schema JsonSchema, err error) {
	s := JsonSchema{
		Schema:     "http://json-schema.org/draft-04/schema#",
		Properties: make(map[string]*JsonSchema),
	}

	if err = toJsonSchema(&s, goNode); err != nil {
		return
	}

	return s, nil
}

func toJsonSchema(root *JsonSchema, goNode GoNode) (err error) {
	switch node := goNode.(type) {
	case *GoIdent:
		{
			if !node.HasObject() {

				if isBasicType(node.Name()) {
					root.Type = node.Name()
				}

				typ, exist := node.rootExpr.options.GoPackage.FindType(node.Name())
				if exist {
					if err = toJsonSchema(root, typ.Node()); err != nil {
						return
					}
				}

				break
			}

			prop := &JsonSchema{
				Id:         root.Id,
				Properties: make(map[string]*JsonSchema),
			}

			node.Inspect(func(n GoNode, ctx context.Context) bool {

				if err = toJsonSchema(prop, n); err != nil {
					return false
				}

				for k, v := range prop.Properties {
					root.Properties[k] = v
				}

				return true
			}, nil)
		}
	case *GoStruct:
		{

			if len(root.Type) == 0 {
				root.Type = "object"
			}

			for i := 0; i < node.NumFields(); i++ {
				if err = toJsonSchema(root, node.Field(i)); err != nil {
					return
				}
			}
		}
	case *GoSelector:
		{

			if !node.IsInOtherPackage() {
				break
			}

			pkg := node.UsingPackage()

			typ, exist := pkg.FindType(node.GetSelName())

			if !exist {
				err = fmt.Errorf("could not found type %s", node.GetSelName())
				return
			}

			switch goType := typ.Node().(type) {
			case *GoType:
				{
					if goType.MethodByName("String") != nil {
						root.Type = "string"
						break
					}

					root.Type = "object"

					prop := &JsonSchema{
						Id:         fmt.Sprintf("%s/items", root.Id),
						Properties: make(map[string]*JsonSchema),
					}

					if err = toJsonSchema(prop, goType.Node()); err != nil {
						return
					}

					root.Items = prop

					break
				}
			}

			if len(root.Type) == 0 {
				root.Type = goTypeToJsonType(typ.String())
			}
		}
	case *GoArray:
		{
			prop := &JsonSchema{
				Id:         fmt.Sprintf("%s/items", root.Id),
				Properties: make(map[string]*JsonSchema),
			}

			root.Type = "array"

			node.Inspect(func(n GoNode, ctx context.Context) bool {
				if err = toJsonSchema(prop, n); err != nil {
					return false
				}

				return false
			}, nil)

			root.Items = prop
		}
	case *GoMap:
		{
			root.Type = "object"

			prop := &JsonSchema{
				Id:         fmt.Sprintf("%s/items", root.Id),
				Properties: make(map[string]*JsonSchema),
			}

			if err = toJsonSchema(prop, node.Value().Node()); err != nil {
				return
			}

			root.Items = prop
		}
	case *GoInterface:
		{
			root.Type = "object"
		}
	case *GoField:
		{
			goFieldToJsonSchema(root, node)
		}
	case *GoExpr:
		{
			node.Inspect(func(n GoNode, ctx context.Context) bool {
				toJsonSchema(root, n)
				return false
			}, nil)
		}
	case *GoStar:
		{
			if err = toJsonSchema(root, node.X()); err != nil {
				return
			}
		}
	case *GoType:
		{
			if err = toJsonSchema(root, node.Node()); err != nil {
				return
			}
		}
	case *GoBasicLit:
		{
			name := strings.ToLower(node.Kind())
			prop := &JsonSchema{
				Id:   fmt.Sprintf("%s/properties/%s", root.Id, name),
				Type: goTypeToJsonType(name),
			}

			root.Properties[name] = prop
		}
	case *GoCompositeLit:
		{
			node.Inspect(func(n GoNode, ctx context.Context) bool {
				toJsonSchema(root, n)
				return false
			}, nil)
		}
	}

	return
}

func goFieldToJsonSchema(root *JsonSchema, field *GoField) (err error) {

	if !field.IsExported() {
		return
	}

	typ := field.Type()

	name := ""
	isCombine := field.NumName() == 0

	if !isCombine {
		name = field.Name(0).Name()
	}

	strType := typ.String()
	jsonTag, ok := field.Tag().Lookup("json")

	if ok {
		tags := strings.Split(jsonTag, ",")

		if len(tags) > 0 {

			if tags[0] == "-" { // ignore field
				return
			}

			name = tags[0]

			if len(tags) > 1 {
				if tags[1] == "string" { // json tag of string
					strType = "string"
				}
			}
		}
	}

	if isBasicType(strType) {

		if name == "" {
			name = strType // if no-name, it may combined
		}

		prop := &JsonSchema{
			Id:   fmt.Sprintf("%s/properties/%s", root.Id, name),
			Type: goTypeToJsonType(strType),
		}

		root.Properties[name] = prop

		return
	}

	nextGoType := typ.Node()

	var propName string

	if len(name) == 0 {
		propName = root.Id
	} else {
		propName = fmt.Sprintf("%s/properties/%s", root.Id, name)
	}

	prop := &JsonSchema{
		Id:         propName,
		Properties: make(map[string]*JsonSchema),
	}

	toJsonSchema(prop, nextGoType)

	if len(name) == 0 {
		for k, v := range prop.Properties {
			root.Properties[k] = v
		}
	} else {
		root.Properties[name] = prop
	}

	return
}

func goTypeToJsonType(typ string) (ret string) {
	switch typ {
	case "string":
		ret = "string"
	case "byte", "int", "int32", "int64", "int8", "rune":
		ret = "integer"
	case "float32", "float64":
		ret = "number"
	case "bool":
		ret = "boolean"
	case "interface{}":
		ret = "object"
	}

	return ret
}

func isBasicType(typ string) bool {
	if typ == "string" ||
		typ == "byte" ||
		typ == "int" ||
		typ == "int8" ||
		typ == "int32" ||
		typ == "int64" ||
		typ == "float32" ||
		typ == "float64" ||
		typ == "bool" ||
		typ == "rune" {
		return true
	}

	return false
}
