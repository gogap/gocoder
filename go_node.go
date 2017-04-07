package gocoder

import (
	"go/token"
)

type GoNode interface {
	goNode()
	Position() token.Position
	Print() error
}

type GoNodeInspectable interface {
	Inspect(f func(GoNode) bool)
}
