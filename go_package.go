package gocoder

import (
	"fmt"
	"os"
	"path/filepath"
)

type GoPackageOption func(*GoPackage) error

type GoPackage struct {
	options *Options

	pkgPath string
	pkgDir  string

	gofiles []*GoFile
}

func NewGoPackage(pkgPath string, options ...Option) (goPackage *GoPackage, err error) {
	pkg := &GoPackage{
		pkgPath: pkgPath,
		options: &Options{},
	}

	if err = pkg.options.init(options...); err != nil {
		return
	}

	if len(pkg.options.GoPath) == 0 {
		pkg.options.GoPath = os.Getenv("GOPATH")
	}

	pkg.pkgDir = filepath.Join(pkg.options.GoPath, "/src/", pkg.pkgPath)

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

func (p *GoPackage) GoFiles() []*GoFile {
	return p.gofiles
}

func (p *GoPackage) load() error {

	walkFn := func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(filename) != ".go" {
			return nil
		}

		gofile, err := NewGoFile(filename, p.options.Copy()...)

		if err != nil {
			return err
		}

		p.gofiles = append(p.gofiles, gofile)

		return nil
	}

	return filepath.Walk(p.pkgDir, walkFn)
}
