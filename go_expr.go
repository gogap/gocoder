package gocoder

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"sync"
)

type InspectFunc func(GoNode, context.Context) bool

type GoExpr struct {
	rootExpr *GoExpr

	astFileSet *token.FileSet
	astFile    *ast.File
	expr       ast.Node

	walkCache []GoNode
	walkOnce  sync.Once

	options Options

	isRoot bool

	exprNode GoNode
}

func newRootGoExpr(astFile *ast.File, astFileSet *token.FileSet) *GoExpr {
	expr := &GoExpr{
		astFile:    astFile,
		astFileSet: astFileSet,
		isRoot:     true,
	}

	expr.rootExpr = expr

	return expr
}

func newGoExpr(rootExpr *GoExpr, expr ast.Node, options ...Option) *GoExpr {

	goExpr := &GoExpr{
		rootExpr: rootExpr,
		expr:     expr,
	}

	var opts []Option

	var rootOptions []Option
	if rootExpr != nil {
		rootOptions = rootExpr.options.Copy()
	}

	if len(rootOptions) > 0 {
		opts = append(opts, rootOptions...)
	}

	if len(options) > 0 {
		opts = append(opts, options...)
	}

	goExpr.options.init(opts...)

	return goExpr
}

func (p *GoExpr) Print() error {
	return ast.Print(p.astFileSet, p.expr)
}

func (p *GoExpr) Root() *GoExpr {
	return p.rootExpr
}

func (p *GoExpr) Node() GoNode {

	if p.exprNode != nil {
		return p.exprNode
	}

	p.exprNode = p.exprToGoNode(p.expr)

	return p.exprNode
}

func (p *GoExpr) goNode() {}

func (p *GoExpr) Options() Options {
	return p.options
}

func (p *GoExpr) exprToGoNode(n ast.Node) GoNode {
	var goNode GoNode

	switch nodeType := n.(type) {
	case *ast.Ident:
		{
			goNode = newGoIdent(p.rootExpr, nodeType)
		}
	case *ast.CallExpr:
		{
			goNode = newGoCall(p.rootExpr, nodeType)
		}
	case *ast.FuncDecl:
		{
			goNode = newGoFunc(p.rootExpr, nodeType)
		}
	case *ast.AssignStmt:
		{
			goNode = newGoAssignStmt(p.rootExpr, nodeType)
		}
	case *ast.FieldList:
		{
			goNode = newFieldList(p.rootExpr, nodeType)
		}
	case *ast.Field:
		{
			goNode = newGoField(p.rootExpr, nodeType)
		}
	case *ast.UnaryExpr:
		{
			goNode = newGoUnary(p.rootExpr, nodeType)
		}
	case *ast.BasicLit:
		{
			goNode = newGoBasicLit(p.rootExpr, nodeType)
		}
	case *ast.CompositeLit:
		{
			goNode = newGoCompositeLit(p.rootExpr, nodeType)
		}
	case *ast.SelectorExpr:
		{
			goNode = newGoSelector(p.rootExpr, nodeType)
		}
	case *ast.TypeSpec:
		{
			goNode = newGoType(p.rootExpr, nodeType, nodeType.Type)
		}
	case *ast.StarExpr:
		{
			goNode = newGoStar(p.rootExpr, nodeType)
		}
	case *ast.MapType:
		{
			goNode = newGoMap(p.rootExpr, nodeType)
		}
	case *ast.InterfaceType:
		{
			goNode = newGoInterface(p.rootExpr, nodeType)
		}
	case *ast.ArrayType:
		{
			goNode = newGoArray(p.rootExpr, nodeType)
		}
	}

	return goNode
}

func (p *GoExpr) walk() {

	ast.Inspect(p.expr, func(n ast.Node) bool {

		goNode := p.exprToGoNode(n)

		if goNode == nil {
			return false
		}

		p.walkCache = append(p.walkCache, goNode)

		return false
	})
}

func (p *GoExpr) Inspect(f InspectFunc, ctx context.Context) {
	p.walkOnce.Do(func() {
		p.walk()
	})

	for i := 0; i < len(p.walkCache); i++ {
		if !f(p.walkCache[i], ctx) {
			return
		}
	}
}

func (p *GoExpr) Position() (begin token.Position, end token.Position) {
	if p.astFileSet == nil {
		return p.rootExpr.astFileSet.Position(p.expr.Pos()), p.rootExpr.astFileSet.Position(p.expr.End())
	}
	return p.astFileSet.Position(p.expr.Pos()), p.astFileSet.Position(p.expr.End())
}

func (p *GoExpr) IsIdent() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.Ident)

	return ok
}

func (p *GoExpr) IsArray() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.ArrayType)
	return ok
}

func (p *GoExpr) IsInterface() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.InterfaceType)
	return ok
}

func (p *GoExpr) IsMap() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.MapType)
	return ok
}

func (p *GoExpr) IsStruct() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.StructType)
	return ok
}

func (p *GoExpr) IsSelector() bool {

	expr := p.expr

	if n, ok := p.expr.(*ast.TypeSpec); ok {
		expr = n
	}

	_, ok := expr.(*ast.SelectorExpr)
	return ok
}

func (p *GoExpr) Name() string {

	switch n := p.expr.(type) {
	case *ast.FuncDecl:
		{
			return n.Name.Name
		}
	case *ast.TypeSpec:
		{
			return n.Name.Name
		}
	case *ast.ImportSpec:
		{
			return n.Name.Name
		}
	case *ast.File:
		{
			return n.Name.Name
		}
	case *ast.ValueSpec:
		{
			if len(n.Names) > 0 {
				return n.Names[0].Name
			}
		}
	}

	return ""
}

// TODO
func (p *GoExpr) Snippet() (code string, err error) {
	return
}

func (p *GoExpr) Type() string {
	switch exp := p.expr.(type) {
	case *ast.TypeSpec:
		{
			switch exp.Type.(type) {
			case *ast.StructType:
				{
					return fmt.Sprintf("struct.%s", exp.Name)
				}
			case *ast.InterfaceType:
				{
					return fmt.Sprintf("interface{}.%s", exp.Name)
				}
			}
		}
	}

	str, ok := p.Node().(fmt.Stringer)
	if ok {
		return str.String()
	}

	return astTypeToStringType(p.expr)
}

func (p *GoExpr) String() string {
	return fmt.Sprintf("%s: %s", p.Name(), p.Type())
}
