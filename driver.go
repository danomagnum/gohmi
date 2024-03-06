package main

import "io"

type Driver interface {
	Read(key string) (any, error)
	Write(key string, value string) error
	Start() error
	Stop() error
	Status() string
	Name() string

	Load(r io.Reader) error
	Save(w io.Writer) error
}
