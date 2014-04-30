package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "1.0"

func errLog(msg string) {
	fmt.Fprintf(os.Stderr, msg)
}

func process(class string, fs []string) int {

	foundOnce := false
	errorOnce := false

	canonicalClassName := strings.Replace(class, ".", "/", -1) + ".class"

	for _, zipFile := range fs {
		zip, err := zip.OpenReader(zipFile)
		if err != nil {
			errLog(zipFile + ": Invalid jar file\n")
			errorOnce = true
		} else {
			defer zip.Close()
			for _, f := range zip.File {
				if canonicalClassName == f.Name {
					fmt.Printf("%s:\t %s\n", zipFile, class)
					foundOnce = true
				}
			}
		}
	}
	if errorOnce {
		return 2
	}
	if !foundOnce {
		// Not found exit code is 1
		return 1
	}
	return 0
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-version") {
		fmt.Printf("grepj %s \n", VERSION)
		fmt.Printf("Written by Rahul Revo, see <https://github.com/rrevo/grepj>.\n")
		os.Exit(0)
	}
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <class-name> <file 1> ... <file n> \n", filepath.Base(os.Args[0]))
		fmt.Printf("          class-name is searched in the files provided\n")
		os.Exit(2)
	}

	class := os.Args[1]
	fs := os.Args[2:]

	os.Exit(process(class, fs))
}
