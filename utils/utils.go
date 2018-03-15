package utils

import (
	"runtime"
	"os"
)

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func IsExistDir(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func CreateNotExistDir(path string) error {
	if !IsExistDir(path) {
		return CreateDir(path)
	}

	return nil
}
