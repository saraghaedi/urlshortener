package main

import (
	"fmt"
	"os"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()
	fmt.Println("test drone ci")

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
