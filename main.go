// gogetdeps: `go get` with dependency version pinning
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	wg          sync.WaitGroup
	TheGOPATH   string
	GOPATH      []string
	GOROOT      string
	mainPackage string

	undo = flag.Bool("undo", false, "remove all traces of gogetdeps from this project")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Call %s from your project's main package.\n\nPossible flags are:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if len(flag.Args()) != 0 {
		flag.Usage()
	}

	GOPATH = filepath.SplitList(os.Getenv("GOPATH"))
	if len(GOPATH) == 0 {
		log.Fatal("GOPATH must be set.")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Panicf("pwd failed: %v", err)
	}
	wd, err = filepath.EvalSymlinks(wd)
	if err != nil {
		log.Panicf("pwd failed: %v", err)
	}

	GOROOT = runtime.GOROOT()
	if !filepath.IsAbs(GOROOT) {
		log.Fatalf("GOROOT %q i not absolute.", GOROOT)
	}
	GOROOT, err = filepath.EvalSymlinks(filepath.Join(GOROOT, "src", "pkg"))
	if err != nil {
		log.Fatalf("Error cleaning GOROOT: %v", err)
	}

	for i, p := range GOPATH {
		if !filepath.IsAbs(p) {
			log.Fatalf("GOPATH %q is not absolute.", p)
		}
		GOPATH[i], err = filepath.EvalSymlinks(p)
		if err != nil {
			log.Fatalf("Error cleaning GOPATH %q: %v", p, err)
		}
		GOPATH[i] = filepath.Join(GOPATH[i], "src")
		if strings.HasPrefix(wd, GOPATH[i]) {
			TheGOPATH = GOPATH[i]
			mainPackage, err = filepath.Rel(TheGOPATH, wd)
			if err != nil {
				log.Panicf("relpath failed: %v", err)
			}
		}
	}

	if mainPackage == "" {
		flag.Usage()
	}

	err = os.RemoveAll(filepath.Join(TheGOPATH, mainPackage, "external"))
	if err != nil {
		log.Fatalf("Error removing old package cache: %v", err)
	}

	wg.Add(1)
	Process(mainPackage)
	wg.Wait()
}
