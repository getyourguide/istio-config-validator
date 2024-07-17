package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-s] -t <testcases1.yml|testcasesdir1> [-t <testcases2.yml|testcasesdir2> ...] <istioconfig1.yml|istioconfigdir1> [<istioconfig2.yml|istioconfigdir2> ...]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	var testCaseParams multiValueFlag
	flag.Var(&testCaseParams, "t", "Testcase files/folders")
	summaryOnly := flag.Bool("s", false, "show only summary of tests (in case of failures full details are shown)")
	strict := flag.Bool("strict", false, "fail on unknown fields")

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

	summary, details, err := unit.Run(testCaseFiles, istioConfigFiles, *strict)
	if err != nil {
		fmt.Println(strings.Join(details, "\n"))
		log.Fatal(err.Error())
	}
	if !*summaryOnly {
		fmt.Println(strings.Join(details, "\n"))
		fmt.Println("")
	}
	fmt.Println(strings.Join(summary, "\n"))
}

func getFiles(names []string) []string {
	var files []string
	for _, name := range names {
		err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err.Error())
			}
			if !info.IsDir() && isYaml(info) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return files
}

func isYaml(info os.FileInfo) bool {
	extension := filepath.Ext(info.Name())
	return extension == ".yaml" || extension == ".yml"
}
