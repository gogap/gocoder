package gocoder

import (
	"context"
	"go/token"
)

type GoNode interface {
	goNode()
	Position() (begin token.Position, end token.Position)
	Print() error
}

type GoNodeInspectable interface {
	Inspect(f InspectFunc, ctx context.Context)
}
