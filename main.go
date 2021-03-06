package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nomad-software/replace/cli"
	"github.com/nomad-software/replace/file"
)

func main() {

	options := cli.ParseOptions()

	if options.Help {
		options.Usage()

	} else if options.Valid() {
		file := file.NewHandler(&options)

		err := file.Walk()
		if err != nil {
			fmt.Fprintln(cli.Stderr, color.RedString(err.Error()))
		}
	}
}
