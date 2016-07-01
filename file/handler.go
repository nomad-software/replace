package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/nomad-software/replace/cli"
)

type Handler struct {
	Options *cli.Options
	Group   sync.WaitGroup
	Output  *cli.Output
}

func (this *Handler) Init(options *cli.Options) {
	this.Options = options
	this.Output = &cli.Output{
		Console: make(chan string),
		Closed:  make(chan bool),
	}
}

func (this *Handler) handlePath(fullPath string) {
	defer this.Group.Done()

	matched, err := filepath.Match(this.Options.File, path.Base(fullPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	if matched {
		this.Group.Add(1)
		go this.processPath(fullPath)
	}
}

func (this *Handler) Walk() error {
	return filepath.Walk(this.Options.Dir, func(fullPath string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		this.Group.Add(1)
		go this.handlePath(fullPath)

		return nil
	})
}

func (this *Handler) processPath(fullPath string) {
	defer this.Group.Done()

	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	newContents := strings.Replace(string(contents), this.Options.From, this.Options.To, -1)
	if newContents != string(contents) {

		err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		this.Output.Console <- fullPath
	}
}
