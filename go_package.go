package gocoder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GoPackageOption func(*GoPackage) error

type GoPackage struct {
	options *Options

	pkgPath string
	pkgDir  string

	goFiles map[string]*GoFile
	files   []string

	loadLocker sync.Mutex

	inGoRoot bool
}

func NewGoPackage(pkgPath string, options ...Option) (goPackage *GoPackage, err error) {
	pkg := &GoPackage{
		pkgPath: pkgPath,
		options: &Options{},
		goFiles: make(map[string]*GoFile),
	}

	if err = pkg.options.init(options...); err != nil {
		return
	}

	if len(pkg.options.GoPath) == 0 {
		pkg.options.Fallback(OptionGoPath(os.Getenv("GOPATH")))
	}

	if len(pkg.options.GoRoot) == 0 {
		goroot := ""
		goroot, err = execCommand("go", "env", "GOROOT")
		if err != nil {
			return
		}

		pkg.options.Fallback(OptionGoRoot(goroot))
	}

	fiInGoPath, _ := os.Stat(filepath.Join(pkg.options.GoPath, "/src/", pkg.pkgPath))
	fiInGoRoot, _ := os.Stat(filepath.Join(pkg.options.GoRoot, "/src/", pkg.pkgPath))

	if fiInGoPath != nil {
		pkg.pkgDir = filepath.Join(pkg.options.GoPath, "/src/", pkg.pkgPath)
	} else if fiInGoRoot != nil {
		pkg.pkgDir = filepath.Join(pkg.options.GoRoot, "/src/", pkg.pkgPath)
		pkg.inGoRoot = true
	} else {
		return nil, fmt.Errorf("package %s not exist in GOPATH and GOROOT", pkgPath)
	}

	fi, err := os.Stat(pkg.pkgDir)
	if err != nil {
		return
	}

	if !fi.IsDir() {
		err = fmt.Errorf("package path of %s is not a dir", pkg.pkgDir)
		return
	}

	if err = pkg.load(); err != nil {
		return
	}

	if err = pkg.checkCircularImport(); err != nil {
		return
	}

	goPackage = pkg
	return
}

func (p *GoPackage) Name() string {
	return filepath.Base(p.pkgPath)
}

func (p *GoPackage) InGoRoot() bool {
	return p.inGoRoot
}

func (p *GoPackage) Options() Options {
	return *p.options
}

func (p *GoPackage) Path() string {
	return p.pkgPath
}

func (p *GoPackage) PackageDir() string {
	return p.pkgDir
}

func (p *GoPackage) checkCircularImport() (err error) {
	return
}

func (p *GoPackage) NumFile() int {
	return len(p.files)
}

func (p *GoPackage) File(i int) *GoFile {
	filename := p.files[i]

	p.loadLocker.Lock()
	defer p.loadLocker.Unlock()

	gf, exist := p.goFiles[filename]
	if exist {
		return gf
	}

	opts := p.options.Copy()
	opts = append(opts, OptionGoPackage(p))

	gofile, err := NewGoFile(filename, opts...)
	if err != nil {
		panic(err)
	}

	p.goFiles[filename] = gofile

	return gofile
}

func (p *GoPackage) NumFuncs() int {
	num := 0
	for i := 0; i < p.NumFile(); i++ {
		num += p.File(i).NumFuncs()
	}
	return num
}

func (p *GoPackage) Func(funcIndex int) *GoFunc {

	for i := 0; i < p.NumFile(); i++ {
		max := p.File(i).NumFuncs()
		if funcIndex >= max {
			funcIndex -= max
			continue
		}

		return p.File(i).Func(funcIndex)
	}
	return nil
}

func (p *GoPackage) FindType(typeName string) (goType *GoExpr, exist bool) {
	for i := 0; i < p.NumFile(); i++ {
		goType, exist = p.File(i).FindType(typeName)
		if exist {
			return
		}
	}
	return
}

func (p *GoPackage) load() error {

	walkFn := func(filename string, info os.FileInfo, err error) error {

		if strings.Contains(filename, "/.") {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(filename) != ".go" {
			return nil
		}

		p.files = append(p.files, filename)

		return nil
	}

	return filepath.Walk(p.pkgDir, walkFn)
}
