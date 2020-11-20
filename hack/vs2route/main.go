package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"istio.io/pkg/log"

	"github.com/getyourguide/istio-config-validator/internal/pkg/envoy"
	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
)

type multiValueFlag []string

func main() {
	flag.Parse()

	istioConfigFiles := getFiles(flag.Args())

	if len(istioConfigFiles) < 1 {
		fmt.Fprintf(os.Stderr, "Missing istio config file/folder, please provide at least one istio config file or folder\n")
		flag.Usage()
		os.Exit(1)
	}
	configs, err := parser.ParseVirtualServices(istioConfigFiles)
	if err != nil {
		log.Fatal(err.Error())
	}
	envoy.Generate("", configs)

}

func getFiles(names []string) []string {
	var files []string
	for _, name := range names {
		filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err.Error())
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}
