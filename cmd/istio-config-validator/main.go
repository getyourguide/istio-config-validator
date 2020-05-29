package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"istio.io/pkg/log"

	"github.com/getyourguide/istio-config-validator/internal/pkg/unit"
)

type multiValueFlag []string

func (m *multiValueFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiValueFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -t <testcases1.yml|testcasesdir1> [-t <testcases2.yml|testcasesdir2> ...] <istioconfig1.yml|istioconfigdir1> [<istioconfig2.yml|istioconfigdir2> ...]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	var testCaseParams multiValueFlag
	flag.Var(&testCaseParams, "t", "Testcase files/folders")
	flag.Parse()
	istioConfigFiles := getFiles(flag.Args())
	testCaseFiles := getFiles(testCaseParams)

	if len(testCaseFiles) < 1 {
		fmt.Fprintf(os.Stderr, "Missing testcases file/folder, please provide at least one testcases file or folder\n")
		flag.Usage()
		os.Exit(1)
	}
	if len(istioConfigFiles) < 1 {
		fmt.Fprintf(os.Stderr, "Missing istio config file/folder, please provide at least one istio config file or folder\n")
		flag.Usage()
		os.Exit(1)
	}

	err := unit.Run(testCaseFiles, istioConfigFiles)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getFiles(names []string) []string {
	var files []string
	for _, name := range names {
		filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}
