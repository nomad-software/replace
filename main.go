package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/nomad-software/replace/cli"
	"github.com/nomad-software/replace/file"
)

func main() {

	var options cli.Options
	options.Parse()

	var file file.Handler
	file.Init(&options)
	go file.Output.Process()

	if (!options.Valid()) || options.Help {
		options.Usage()

	} else {
		options.Echo()
		err := file.Walk()

		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		file.Group.Wait()

		close(file.Output.Console)
		<-file.Output.Closed
	}
}
