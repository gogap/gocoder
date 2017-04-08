package gocoder

import (
	"context"
	"go/ast"
	"go/token"
	"sync"
)

type InspectFunc func(GoNode, context.Context) bool

type GoExpr struct {
	rootExpr *GoExpr

	astFileSet *token.FileSet
	astFile    *ast.File
	expr       ast.Expr

	walkCache []GoNode
	walkOnce  sync.Once

	options Options

	isRoot bool
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

func newGoExpr(rootExpr *GoExpr, expr ast.Expr, options ...Option) *GoExpr {

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

func (p *GoExpr) Options() Options {
	return p.options
}

func (p *GoExpr) walk() {

	ast.Inspect(p.expr, func(n ast.Node) bool {

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
		case *ast.StructType:
			{
				goNode = newGoStruct(p.rootExpr, nodeType)
			}
		case *ast.TypeSpec:
			{
				goNode = newGoType(p.rootExpr, nodeType)
			}
		}

		if goNode == nil {
			return true
		}

		p.walkCache = append(p.walkCache, goNode)

		return true
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

func (p *GoExpr) Position() token.Position {
	if p.astFileSet == nil {
		return p.rootExpr.astFileSet.Position(p.expr.Pos())
	}
	return p.astFileSet.Position(p.expr.Pos())
}
