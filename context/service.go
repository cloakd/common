package context

type Service interface {
	Id() string
	Configure(ctx *Context) error
	Start() error
}
