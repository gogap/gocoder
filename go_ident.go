package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoIdent struct {
	rootExpr   *GoExpr
	goChildren []GoNode

	astExpr *ast.Ident

	objKind string
}

func newGoIdent(rootExpr *GoExpr, ident *ast.Ident) *GoIdent {
	g := &GoIdent{
		rootExpr: rootExpr,
		astExpr:  ident,
	}

	g.load()

	return g
}

func (p *GoIdent) HasObject() bool {
	return p.astExpr.Obj != nil
}

func (p *GoIdent) ObjectKind() string {
	return p.objKind
}

func (p *GoIdent) load() {

	if p.astExpr.Obj == nil {
		return
	}

	switch p.astExpr.Obj.Kind {
	case ast.Var:
		{
			p.objKind = "var"
			switch expr := p.astExpr.Obj.Decl.(type) {
			case *ast.AssignStmt:
				{
					for i := 0; i < len(expr.Rhs); i++ {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Rhs[i]))
					}
				}
			case *ast.ValueSpec:
				{
					if expr.Type != nil {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Type))
					}
				}
			case *ast.Field:
				{
					if expr.Type != nil {
						p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr.Type))
					}
				}
			}
		}
	case ast.Typ:
		{
			p.objKind = "type"
			switch expr := p.astExpr.Obj.Decl.(type) {
			case *ast.TypeSpec:
				{
					p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr))
				}
			}
		}
	case ast.Fun:
		{
			p.objKind = "func"
			switch expr := p.astExpr.Obj.Decl.(type) {
			case *ast.FuncDecl:
				{
					p.goChildren = append(p.goChildren, newGoExpr(p.rootExpr, expr))
				}
			}
		}
	}
}

func (p *GoIdent) Inspect(f InspectFunc, ctx context.Context) {
	for i := 0; i < len(p.goChildren); i++ {
		child, inspectable := p.goChildren[i].(GoNodeInspectable)
		if inspectable {
			child.Inspect(f, ctx)
		}
	}
}

func (p *GoIdent) Name() string {
	return p.astExpr.Name
}

func (p *GoIdent) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoIdent) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoIdent) goNode() {}
