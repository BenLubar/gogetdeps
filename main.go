// gogetdeps: `go get` with dependency version pinning
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup
var GOPATH []string
var mainPackage string

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

	for i, p := range GOPATH {
		if !filepath.IsAbs(p) {
			log.Fatalf("GOPATH %q is not absolute.", p)
		}
		GOPATH[i], err = filepath.EvalSymlinks(p)
		if err != nil {
			log.Fatalf("Error cleaning GOPATH %q: %v", p, err)
		}
		GOPATH[i] = filepath.Join(GOPATH[i], "src")
		if filepath.HasPrefix(wd, GOPATH[i]) {
			mainPackage, err = filepath.Rel(GOPATH[i], wd)
			if err != nil {
				log.Panicf("relpath failed: %v", err)
			}
		}
	}

	if mainPackage == "" {
		flag.Usage()
	}

	Find(mainPackage)
	wg.Wait()
}
