package hmitags

type Tag interface {
	Name() string
	Read(key string) (any, error)
	Write(key string, val any) error
}

type BasicTag struct {
	name  string
	value any
}
