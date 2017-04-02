package gocoder

import (
	"go/ast"
)

type GoParams struct {
	goFunc    *GoFunc
	astParams *ast.FieldList
}

func (p *GoParams) Print() error {
	return ast.Print(p.goFunc.goFile.astFileSet, p.astParams)
}

func (p *GoParams) GoFunc() *GoFunc {
	return p.goFunc
}

func (p *GoParams) NumFields() int {
	return p.astParams.NumFields()
}

func (p *GoParams) TypesAre(paramsType ...string) bool {

	if p.astParams.NumFields() != len(paramsType) {
		return false
	}

	for i := 0; i < p.astParams.NumFields(); i++ {
		if !isFieldTypeOf(p.astParams.List[i], paramsType[i]) {
			return false
		}
	}

	return true
}
