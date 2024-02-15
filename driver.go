package main

type Driver interface {
	Read(key string) (any, error)
	Write(key string, value any) error
	Start() error
	Stop() error
	Status() string
	Name() string
}
