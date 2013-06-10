package main

import (
	"os"
	"path/filepath"
	"strings"
)

func Process(pkg string) {
	defer wg.Done()

	if dir, err := os.Open(filepath.Join(GOROOT, pkg)); err == nil {
		// Don't do anything to GOROOT packages.
		dir.Close()
		return
	}

	if strings.HasPrefix(pkg, mainPackage+"/external/") {
		if *undo {
			return
		}
		wg.Add(1)
		Process(pkg[len(mainPackage+"/external/"):])
		return
	}

	files := Find(pkg)
	if files == nil {
		// already processed
		return
	}

	Rewrite(pkg, files)
}
