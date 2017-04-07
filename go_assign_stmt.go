package gocoder

import (
	"go/ast"
	"go/token"
)

type GoAssignStmt struct {
	rootExpr           *GoExpr
	assignStmtRhsExprs []*GoExpr

	astAssignStmt *ast.AssignStmt
}

func newGoAssignStmt(rootExpr *GoExpr, astAssignStmt *ast.AssignStmt) *GoAssignStmt {
	g := &GoAssignStmt{
		rootExpr:      rootExpr,
		astAssignStmt: astAssignStmt,
	}

	g.load()

	return g
}

func (p *GoAssignStmt) load() {
	for i := 0; i < len(p.astAssignStmt.Rhs); i++ {
		goExpr := newGoExpr(p.rootExpr, p.astAssignStmt.Rhs[i])
		p.assignStmtRhsExprs = append(p.assignStmtRhsExprs, goExpr)
	}
}

func (p *GoAssignStmt) Inspect(f func(GoNode) bool) {
	if len(p.assignStmtRhsExprs) == 0 {
		return
	}

	for i := 0; i < len(p.assignStmtRhsExprs); i++ {
		p.assignStmtRhsExprs[i].Inspect(f)
	}
}

func (p *GoAssignStmt) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astAssignStmt)
}

func (p *GoAssignStmt) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astAssignStmt.Pos())
}

func (p *GoAssignStmt) goNode() {}
