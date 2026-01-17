package main

import (
	"os"
)

type iOpSystem interface {
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	Environ() []string
}

type OpSystem struct {
}

func (os OpSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (os OpSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (os OpSystem) Environ() []string {
	return os.Environ()
}
