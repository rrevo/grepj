package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const version = "2.1"

var extensions = []string{"jar", "war", "ear"}

func errLog(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func debugLog(msg string) {
	if debugFlag {
		fmt.Fprintln(os.Stdout, "debug: "+msg)
	}
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
				debugLog("Search in " + f.Name)
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
						errLog(nestedPath + ": Invalid file")
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
		debugLog("Search in " + zipFileName)
		reader, err := zip.OpenReader(zipFileName)
		if err != nil {
			errLog(zipFileName + ": Invalid file")
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

var versionFlag bool
var debugFlag bool

func init() {
	flag.BoolVar(&versionFlag, "version", false, "Version information")
	flag.BoolVar(&versionFlag, "v", false, "Version information (shorthand)")
	flag.BoolVar(&debugFlag, "Xdebug", false, "Debug output")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Printf("grepj %s \n", version)
		fmt.Printf("Written by Rahul Revo, see <https://github.com/rrevo/grepj>.\n")
		os.Exit(0)
	}

	if flag.NArg() < 2 {
		fmt.Printf("usage: %s <class-name> <file 1> ... <file n> \n", filepath.Base(os.Args[0]))
		fmt.Printf("          class-name is searched in the files provided\n")
		os.Exit(2)
	}

	otherArgs := flag.Args()
	class := otherArgs[0]
	fs := otherArgs[1:]

	error, found := processFiles(class, fs)
	exit := 0
	if error {
		exit = 2
	} else if !found {
		// Not found exit code is 1
		exit = 1
	}
	debugLog("Exit Code - " + strconv.Itoa(exit))
	os.Exit(exit)
}
