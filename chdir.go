package main

import (
	"os"
	"path"
	"runtime"
)

func chdir() {
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Join(path.Dir(filename))

	os.Chdir(dir)
}
