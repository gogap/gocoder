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

type importSpec struct {
	Name string
	Path string
}

type GoFile struct {
	filename      string
	shortFilename string

	goFuncs []*GoFunc
	goTypes []*GoExpr
	goVars  []*GoExpr

	importPackages    []importSpec //path
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

	if strings.HasPrefix(filename, gf.options.GoPath) {
		gf.shortFilename = strings.TrimPrefix(filename, gf.options.GoPath+"/src/")
	} else if strings.HasPrefix(filename, gf.options.GoRoot) {
		gf.shortFilename = strings.TrimPrefix(filename, gf.options.GoRoot+"/src/")
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

func (p *GoFile) NumFuncs() int {
	return len(p.goFuncs)
}

func (p *GoFile) Func(i int) *GoFunc {
	return p.goFuncs[i]
}

func (p *GoFile) NumTypes() int {
	return len(p.goTypes)
}

func (p *GoFile) Type(i int) *GoExpr {
	return p.goTypes[i]
}

func (p *GoFile) NumVars() int {
	return len(p.goVars)
}

func (p *GoFile) Var(i int) *GoExpr {
	return p.goVars[i]
}

func (p *GoFile) Package() *GoPackage {
	return p.options.GoPackage
}

func (p *GoFile) Filename() string {
	return p.filename
}

func (p *GoFile) ShortFilename() string {
	return p.shortFilename
}

func (p *GoFile) String() string {
	return p.shortFilename
}

func (p *GoFile) Imports() []string {

	var imps []string

	for i := 0; i < len(p.importPackages); i++ {
		imps = append(imps, p.importPackages[i].Path)
	}

	return imps
}

func (p *GoFile) InGoRoot() bool {
	return p.Package().inGoRoot
}

func (p *GoFile) FindImportByName(name string) (goPkg *GoPackage, exist bool) {
	for i := 0; i < len(p.importPackages); i++ {
		if name == p.importPackages[i].Name {
			return p.FindImportByPath(p.importPackages[i].Path)
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

func (p *GoFile) FindType(typeName string) (goType *GoExpr, exist bool) {
	for i := 0; i < len(p.GoExpr.astFile.Decls); i++ {
		ast.Inspect(p.GoExpr.astFile.Decls[i], func(n ast.Node) bool {
			if exist {
				return false
			}

			switch node := n.(type) {
			case *ast.TypeSpec:
				{
					if node.Name.Name == typeName {
						goType = newGoExpr(p.rootExpr, node)
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

func (p *GoFile) load() (err error) {

	if err = p.loadImportPackages(); err != nil {
		return
	}

	if err = p.loadFuncDecls(); err != nil {
		return
	}

	if err = p.loadTypeDecls(); err != nil {
		return
	}

	if err = p.loadVarDecls(); err != nil {
		return
	}

	return nil
}

func (p *GoFile) loadImportPackages() (err error) {
	for _, impt := range p.astFile.Imports {

		imptPath := strings.Trim(impt.Path.Value, "\"")

		pathInGopath := filepath.Join(p.options.GoPath, "src", imptPath)
		_, e1 := os.Stat(pathInGopath)

		if e1 == nil {

			spec := importSpec{
				Name: "",
				Path: imptPath,
			}

			if impt.Name != nil {
				spec.Name = impt.Name.Name
			} else {
				spec.Name = filepath.Base(imptPath)
			}

			p.importPackages = append(p.importPackages, spec)
			p.mapImportPackages[imptPath] = nil

			continue
		}

		pathInGoRoot := filepath.Join(p.options.GoRoot, "src", imptPath)
		_, e2 := os.Stat(pathInGoRoot)

		if e2 == nil {

			spec := importSpec{
				Name: "",
				Path: imptPath,
			}

			if impt.Name != nil {
				spec.Name = impt.Name.Name
			} else {
				spec.Name = filepath.Base(imptPath)
			}

			p.importPackages = append(p.importPackages, spec)
			p.mapImportPackages[imptPath] = nil
			continue
		}
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

func (p *GoFile) loadTypeDecls() error {
	for _, decl := range p.astFile.Decls {
		ast.Inspect(decl, func(n ast.Node) bool {

			switch d := n.(type) {
			case *ast.TypeSpec:
				{
					p.goTypes = append(p.goTypes, newGoExpr(p.rootExpr, d))
				}
			}

			return true
		})
	}

	return nil
}

func (p *GoFile) loadVarDecls() error {
	for _, decl := range p.astFile.Decls {
		ast.Inspect(decl, func(n ast.Node) bool {

			switch d := n.(type) {
			case *ast.FuncDecl:
				{
					return false
				}
			case *ast.TypeSpec:
				{
					return false
				}
			case *ast.ValueSpec:
				{
					p.goVars = append(p.goVars, newGoExpr(p.rootExpr, d))
				}
			}

			return true
		})
	}
	return nil
}
