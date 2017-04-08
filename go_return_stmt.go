package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoReturnStmt struct {
	rootExpr *GoExpr
	results  []*GoExpr

	astReturnStmt *ast.ReturnStmt
}

func newGoReturnStmt(rootExpr *GoExpr, astReturnStmt *ast.ReturnStmt) *GoReturnStmt {
	g := &GoReturnStmt{
		rootExpr:      rootExpr,
		astReturnStmt: astReturnStmt,
	}

	g.load()

	return g
}

func (p *GoReturnStmt) load() {
	for i := 0; i < len(p.astReturnStmt.Results); i++ {
		goExpr := newGoExpr(p.rootExpr, p.astReturnStmt.Results[i])
		p.results = append(p.results, goExpr)
	}
}

func (p *GoReturnStmt) Inspect(f InspectFunc, ctx context.Context) {
	for i := 0; i < len(p.results); i++ {
		p.results[i].Inspect(f, ctx)
	}
}

func (p *GoReturnStmt) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astReturnStmt)
}

func (p *GoReturnStmt) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astReturnStmt.Pos())
}

func (p *GoReturnStmt) goNode() {}
