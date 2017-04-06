package gocoder

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GoFile struct {
	filename string
	gopath   string

	goFuncs []*GoFunc

	importPackages    []string //path
	mapImportPackages map[string]*GoPackage
	mapStructs        map[string]*GoStruct

	syncNewImportLocker sync.Mutex

	*GoExpr
}

func NewGoFile(filename string, options ...Option) (goFile *GoFile, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return
	}

	gf := &GoFile{
		filename: filename,
		GoExpr: &GoExpr{
			astFile:    f,
			astFileSet: fset,
		},
		mapStructs:        make(map[string]*GoStruct),
		mapImportPackages: make(map[string]*GoPackage),
	}

	options = append(options, OptionExprInGoFile(gf))

	if err = gf.GoExpr.options.init(options...); err != nil {
		return
	}

	if err = gf.load(); err != nil {
		return
	}

	if err = gf.loadImportPackages(); err != nil {
		return
	}

	gf.loadStructs()

	goFile = gf

	return
}

func (p *GoFile) Print() error {
	return ast.Print(p.astFileSet, p.astFile)
}

func (p *GoFile) Funcs() []*GoFunc {
	return p.goFuncs
}

func (p *GoFile) Package() *GoPackage {
	return p.options.GoPackage
}

func (p *GoFile) Filename() string {
	return p.filename
}

func (p *GoFile) ShortFilename() string {
	return strings.TrimPrefix(p.filename, p.options.GoPath+"/src/")
}

func (p *GoFile) Imports() []string {
	return p.importPackages
}

func (p *GoFile) FindImportByName(name string) (goPkg *GoPackage, exist bool) {
	for i := 0; i < len(p.importPackages); i++ {
		if name == filepath.Base(p.importPackages[i]) {
			return p.FindImportByPath(p.importPackages[i])
		}
	}
	return nil, false
}

func (p *GoFile) FindImportByPath(importPath string) (*GoPackage, bool) {

	pkg, exist := p.mapImportPackages[importPath]
	if exist {
		var err error
		if pkg == nil {
			p.syncNewImportLocker.Lock()
			defer p.syncNewImportLocker.Unlock()

			pkg, err = NewGoPackage(importPath,
				OptionImportByPackage(p.Package()),
				OptionImportByFile(p),
			)

			if err != nil {
				return nil, false
			}
		}
		p.mapImportPackages[importPath] = pkg
		return pkg, true
	}

	return nil, false
}

func (p *GoFile) load() (err error) {
	if err = p.loadDecls(); err != nil {
		return
	}

	return
}

func (p *GoFile) FindStruct(name string) (goStruct *GoStruct, exist bool) {
	goStruct, exist = p.mapStructs[name]
	return
}

func (p *GoFile) loadStructs() {
	for i := 0; i < len(p.GoExpr.astFile.Decls); i++ {
		genDecl, ok := p.GoExpr.astFile.Decls[i].(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok != token.TYPE {
			continue
		}

		if len(genDecl.Specs) != 1 {
			continue
		}

		typeSpec := genDecl.Specs[0].(*ast.TypeSpec)

		// identType, ok := typeSpec.Type.(*ast.Ident)
		// if ok {
		// 	p.structType = identType.Name
		// }

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		p.mapStructs[typeSpec.Name.Name] = newGoStruct(p.rootExpr, structType, OptionExprInGoFile(p))
	}
}

func (p *GoFile) loadImportPackages() (err error) {
	for _, impt := range p.astFile.Imports {
		imptPath := strings.Trim(impt.Path.Value, "\"")
		pathInGopath := filepath.Join(p.options.GoPath, "src", imptPath)

		_, e := os.Stat(pathInGopath)
		if e != nil {
			continue
		}

		p.importPackages = append(p.importPackages, imptPath)
		p.mapImportPackages[imptPath] = nil
	}

	return nil
}

func (p *GoFile) loadDecls() error {
	for _, decl := range p.astFile.Decls {
		ast.Inspect(decl, func(n ast.Node) bool {

			switch d := n.(type) {
			case *ast.FuncDecl:
				{
					p.goFuncs = append(p.goFuncs, newGoFunc(p.GoExpr, d))
				}
			}

			return true
		})
	}
	return nil
}
