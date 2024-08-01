package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PondWader/go-npm-registry/pkg"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Expected flag config to be set.")
		os.Exit(1)
	}

	pkg.StartServer(*configPath)
}
