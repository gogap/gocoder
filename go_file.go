package gocoder

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type GoFile struct {
	filename string
	gopath   string

	astFileSet *token.FileSet
	astFile    *ast.File

	goFuncs []*GoFunc
}

func NewGoFile(filename, gopath string) (goFile *GoFile, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return
	}

	gf := &GoFile{
		filename:   filename,
		gopath:     gopath,
		astFileSet: fset,
		astFile:    f,
	}

	if err = gf.load(); err != nil {
		return
	}

	goFile = gf

	return
}

func (p *GoFile) Print() error {
	return ast.Print(p.astFileSet, p.astFile)
}

func (p *GoFile) Funcs() []*GoFunc {
	return p.goFuncs
}

func (p *GoFile) load() (err error) {
	if err = p.loadDecls(); err != nil {
		return
	}

	return
}

func (p *GoFile) loadDecls() error {
	for _, decl := range p.astFile.Decls {
		ast.Inspect(decl, func(n ast.Node) bool {

			switch d := n.(type) {
			case *ast.FuncDecl:
				{
					p.parseDeclFunc(d)
				}
			}

			return true
		})
	}
	return nil
}

func (p *GoFile) parseDeclFunc(decl *ast.FuncDecl) {
	p.goFuncs = append(p.goFuncs, newGoFunc(p, decl))
}
