package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const version = "2.0"

var extensions = []string{"jar", "war", "ear"}

func errLog(msg string) {
	fmt.Fprintf(os.Stderr, msg)
}

func processZipReader(reader *zip.Reader, searchPatterns []string, class string, path string) (error, found bool) {
	for _, f := range reader.File {
		for _, searchPattern := range searchPatterns {
			if searchPattern == f.Name {
				fmt.Printf("%s:\t %s\n", path, class)
				found = true
			}
		}
		for _, extension := range extensions {
			if strings.HasSuffix(f.Name, extension) {
				fReader, err := f.Open()
				if err != nil {
					errLog("Unable to open " + f.Name + " in " + path)
				} else {
					defer fReader.Close()
					b, errB := ioutil.ReadAll(fReader)
					if errB != nil {
						errLog("Error reading " + f.Name + " in " + path)
					}
					nestedReader, nestedErr := zip.NewReader(bytes.NewReader(b), int64(len(b)))
					nestedPath := path + ">" + f.Name
					if nestedErr != nil {
						errLog(nestedPath + ": Invalid file\n")
					} else {
						nestedError, nestedFound := processZipReader(nestedReader, searchPatterns, class, nestedPath)
						error = error || nestedError
						found = found || nestedFound
					}
				}
			}
		}
	}
	return
}

func processFiles(class string, fs []string) (error, found bool) {

	canonicalClassName := strings.Replace(class, ".", "/", -1) + ".class"
	searchPatterns := []string{canonicalClassName, "WEB-INF/classes/" + canonicalClassName}

	for _, zipFileName := range fs {
		reader, err := zip.OpenReader(zipFileName)
		if err != nil {
			errLog(zipFileName + ": Invalid file\n")
			error = true
		} else {
			defer reader.Close()
			nestedError, nestedFound := processZipReader(&reader.Reader, searchPatterns, class, zipFileName)
			error = error || nestedError
			found = found || nestedFound
		}
	}
	return
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-version") {
		fmt.Printf("grepj %s \n", version)
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

	error, found := processFiles(class, fs)
	exit := 0
	if error {
		exit = 2
	} else if !found {
		// Not found exit code is 1
		exit = 1
	}
	os.Exit(exit)
}
