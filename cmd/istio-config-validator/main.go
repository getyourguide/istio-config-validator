package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -t <testcases1.yml> [-t <testcases2.yml> ...] <istioconfig1.yml> [<istioconfig2.yml> ...]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	var testCaseFiles multiValueFlag
	flag.Var(&testCaseFiles, "t", "Testcase files")
	flag.Parse()
	istioConfigFiles := flag.Args()

	if len(testCaseFiles) < 1 {
		fmt.Fprintf(os.Stderr, "Missing testcases file, please provide at least one testcases file\n")
		flag.Usage()
		os.Exit(1)
	}
	if len(istioConfigFiles) < 1 {
		fmt.Fprintf(os.Stderr, "Missing istio config file, please provide at least one istio config file\n")
		flag.Usage()
		os.Exit(1)
	}

	// TODO: instead of printing out the file names, call the parsing&validation functions
	fmt.Println("Received the following test case files:")
	for _, file := range testCaseFiles {
		fmt.Println(file)
	}
	fmt.Println("Received the following istio config files:")
	for _, file := range istioConfigFiles {
		fmt.Println(file)
	}
}
