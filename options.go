package gocoder

import (
	"fmt"
	"os"
)

type Option func(*Options) error

type Options struct {
	GoPath    string
	GoPackage *GoPackage
	ImportBy  *GoPackage
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

func OptionImportBy(pkg *GoPackage) Option {
	return func(g *Options) (err error) {
		g.ImportBy = pkg
		return
	}
}

func OptionGoPackage(goPkg *GoPackage) Option {
	return func(g *Options) (err error) {
		g.GoPackage = goPkg
		return nil
	}
}
