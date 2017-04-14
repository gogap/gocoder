package gocoder

import (
	"fmt"
	"os"
)

type Option func(*Options) error

type Options struct {
	GoPath          string
	GoRoot          string
	GoPackage       *GoPackage
	GoFile          *GoFile
	ImportByPackage *GoPackage
	ImportByFile    *GoFile

	options []Option
}

func (p *Options) init(options ...Option) (err error) {
	for i := 0; i < len(options); i++ {
		if err = options[i](p); err != nil {
			return
		}
	}

	p.options = options

	if len(p.GoPath) == 0 {
		p.GoPath = os.Getenv("GOPATH")
	}

	return
}

func (p *Options) Fallback(opts ...Option) (err error) {

	for i := 0; i < len(opts); i++ {
		if err = opts[i](p); err != nil {
			return
		}
	}

	p.options = append(p.options, opts...)

	return nil
}

func (p *Options) Copy() []Option {
	var options []Option

	for i := 0; i < len(p.options); i++ {
		options = append(options, p.options[i])
	}

	return options
}

func OptionGoPath(gopath string) Option {
	return func(g *Options) (err error) {
		if len(gopath) == 0 {
			gopath = os.Getenv("GOPATH")
		}

		g.GoPath = gopath

		fi, err := os.Stat(gopath)
		if err != nil {
			return
		}

		if !fi.IsDir() {
			err = fmt.Errorf("gopath of %s is not a dir", gopath)
			return
		}

		return
	}
}

func OptionGoRoot(goroot string) Option {
	return func(g *Options) (err error) {

		if len(goroot) == 0 {
			goroot := ""
			goroot, err = execCommand("go", "env", "GOROOT")
			if err != nil {
				return
			}

			g.GoRoot = goroot
			return nil
		}

		g.GoRoot = goroot

		return nil
	}
}

func OptionExprInGoFile(gofile *GoFile) Option {
	return func(g *Options) (err error) {
		g.GoFile = gofile
		return
	}
}

func OptionImportByPackage(pkg *GoPackage) Option {
	return func(g *Options) (err error) {
		g.ImportByPackage = pkg
		return
	}
}

func OptionImportByFile(file *GoFile) Option {
	return func(g *Options) (err error) {
		g.ImportByFile = file
		return
	}
}

func OptionGoPackage(goPkg *GoPackage) Option {
	return func(g *Options) (err error) {
		g.GoPackage = goPkg
		return nil
	}
}
