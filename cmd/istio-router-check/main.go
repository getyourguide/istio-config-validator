package main

import (
	"fmt"
	"os"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/cmd"
)

func main() {
	cmdRoot, err := cmd.NewCmdRoot()
	if err != nil {
		fmt.Printf("failed to create command: %v", err)
		os.Exit(1)
	}
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
