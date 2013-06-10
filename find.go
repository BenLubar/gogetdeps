package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var Visited = struct {
	Pkg map[string]bool
	sync.Mutex
}{
	Pkg: make(map[string]bool),
}

// If the returned slice is nil, the package has already been processed. Otherwise, the returned slice holds the
// full paths of all source files inside the package, excluding tests.
func Find(base string) (found []string) {
	Visited.Lock()
	if Visited.Pkg[base] {
		Visited.Unlock()
		return nil
	}
	Visited.Pkg[base] = true
	Visited.Unlock()

	found = []string{}

	for _, p := range GOPATH {
		dir, err := os.Open(filepath.Join(p, base))
		if err != nil {
			continue
		}
		defer dir.Close()

		files, err := dir.Readdir(0)
		if err != nil {
			log.Fatalf("Error listing contents of package %q: %v", base, err)
		}

		fset := token.NewFileSet()

		for _, f := range files {
			if f.Mode().IsRegular() {
				if strings.HasSuffix(f.Name(), ".c") || strings.HasSuffix(f.Name(), ".h") || strings.HasSuffix(f.Name(), ".s") {
					found = append(found, filepath.Join(p, base, f.Name()))
				} else if strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
					name := filepath.Join(p, base, f.Name())
					found = append(found, name)
					ast, err := parser.ParseFile(fset, name, nil, parser.ImportsOnly)
					if err != nil {
						log.Fatalf("Error parsing file %q: %v", name, err)
					}
					for _, i := range ast.Imports {
						log.Print(i.Path.Value[1 : len(i.Path.Value)-1])
					}
				}
			}
		}

		return found
	}

	log.Fatalf("Could not find package %q", base)
	panic("unreachable")
}
