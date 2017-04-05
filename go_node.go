package gocoder

type GoNode interface {
	goNode()
	Print() error
}

type GoNodeInspectable interface {
	Inspect(f func(GoNode) bool)
}
