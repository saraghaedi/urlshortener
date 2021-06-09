package main

import (
	"os"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
