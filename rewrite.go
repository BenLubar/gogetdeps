package main

import (
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RewritePath(in string) string {
	if *undo {
		return strings.TrimPrefix(in, mainPackage + "/external/")
	}

	if dir, err := os.Open(filepath.Join(GOROOT, in)); err == nil {
		// Don't do anything to GOROOT packages.
		dir.Close()
		return in
	}

	if strings.HasPrefix(in, mainPackage) || strings.HasPrefix(mainPackage, in) {
		return in
	}

	return mainPackage + "/external/" + in
}

func Rewrite(pkg string, files []string) {
	outpkg := RewritePath(pkg)
	outpath := filepath.Join(TheGOPATH, outpkg)
	if outpkg != pkg {
		if *undo {
			return
		}

		err := os.MkdirAll(outpath, 0755)
		if err != nil {
			log.Fatalf("Error making directory for package cache for %q: err", pkg, err)
		}
	}

	fset := token.NewFileSet()
	for _, fn := range files {
		ast, err := parser.ParseFile(fset, fn, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("Error parsing file %q: %v", fn, err)
		}

		for _, i := range ast.Imports {
			i.Path.Value = `"` + RewritePath(i.Path.Value[1:len(i.Path.Value)-1]) + `"`
		}

		f, err := os.Create(filepath.Join(outpath, filepath.Base(fn)))
		if err != nil {
			log.Fatalf("Error creating file %q: %v", filepath.Join(outpath, filepath.Base(fn)), err)
		}
		err = printer.Fprint(f, fset, ast)
		f.Close()
		if err != nil {
			log.Fatalf("Error writing file %q: %v", filepath.Join(outpath, filepath.Base(fn)), err)
		}
	}
}
