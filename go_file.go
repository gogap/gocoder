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
		filename:          filename,
		GoExpr:            newRootGoExpr(f, fset),
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
	if err = p.loadFuncDecls(); err != nil {
		return
	}

	return
}

func (p *GoFile) FindType(typeName string) (goType *GoType, exist bool) {
	for i := 0; i < len(p.GoExpr.astFile.Decls); i++ {
		ast.Inspect(p.GoExpr.astFile.Decls[i], func(n ast.Node) bool {
			if exist {
				return false
			}

			switch node := n.(type) {
			case *ast.TypeSpec:
				{
					if node.Name.Name == typeName {
						goType = newGoType(p.rootExpr, node)
						exist = true
						return false
					}

					return true
				}
			}
			return true
		})
	}

	return
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

func (p *GoFile) loadFuncDecls() error {
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
