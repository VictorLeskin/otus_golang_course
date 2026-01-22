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

type OpSystem struct{}

// For this derived all calls are system calls.
func (OpSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OpSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OpSystem) Environ() []string {
	return os.Environ()
}

func (OpSystem) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
