package gocoder

import (
	"go/ast"
)

type GoResults struct {
	goFunc     *GoFunc
	astResults *ast.FieldList
}

func (p *GoResults) GoFunc() *GoFunc {
	return p.goFunc
}

func (p *GoResults) NumFields() int {
	return p.astResults.NumFields()
}

func (p *GoResults) TypesAre(resultsType ...string) bool {

	if p.astResults.NumFields() != len(resultsType) {
		return false
	}

	for i := 0; i < p.astResults.NumFields(); i++ {
		if !isFieldTypeOf(p.astResults.List[i], resultsType[i]) {
			return false
		}
	}

	return true
}
