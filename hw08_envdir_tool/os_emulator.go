package main

import (
	"os"
	"os/exec"
)

type iOpSystem interface {
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	Environ() []string
	Run(cmd *exec.Cmd) error
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

func (os OpSystem) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
