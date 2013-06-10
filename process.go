package main

import (
	"log"
	"os"
	"path/filepath"
)

func Process(pkg string) {
	defer wg.Done()

	if dir, err := os.Open(filepath.Join(GOROOT, pkg)); err == nil {
		// Don't do anything to GOROOT packages.
		dir.Close()
		return
	}

	_ = Find(pkg)

	log.Print(pkg)
}
