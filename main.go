package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type cliOptions struct {
	dir  string
	from string
	to   string
	file string
	help bool
}

func (this *cliOptions) valid() bool {
	return this.from != "" && this.to != ""
}

func (this *cliOptions) expandDir() {
	dir, _ := homedir.Expand(this.dir)
	this.dir = dir
}

func (this *cliOptions) display() {
	options := color.CyanString("replacing:   ")
	options += color.GreenString("%s\n", this.from)
	options += color.CyanString("with:        ")
	options += color.GreenString("%s\n", this.to)
	options += color.CyanString("in files:    ")
	options += color.GreenString("%s\n", this.file)
	options += color.CyanString("starting in: ")
	options += color.GreenString("%s\n", this.dir)
	fmt.Print(options)
}

var wg sync.WaitGroup
var banner string = ` ____            _
|  _ \ ___ _ __ | | __ _  ___ ___
| |_) / _ \ '_ \| |/ _' |/ __/ _ \
|  _ <  __/ |_) | | (_| | (_|  __/
|_| \_\___| .__/|_|\__'_|\___\___|
          |_|

`

func main() {

	var options cliOptions

	flag.StringVar(&options.dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&options.from, "from", "", "The text to replace.")
	flag.StringVar(&options.to, "to", "", "The replacement text.")
	flag.StringVar(&options.file, "file", "*", "The glob file pattern to match.")
	flag.BoolVar(&options.help, "help", false, "Show help.")
	flag.Parse()

	if (!options.valid()) || options.help {
		color.Cyan(banner)
		flag.Usage()

	} else {
		options.expandDir()
		options.display()

		err := filepath.Walk(options.dir, func(fullPath string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			wg.Add(1)
			go handlePath(fullPath, &options)

			return nil
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		wg.Wait()
	}
}

func processPath(fullPath string, options *cliOptions) {
	defer wg.Done()

	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	newContents := strings.Replace(string(contents), options.from, options.to, -1)
	if newContents != string(contents) {

		err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		color.Magenta(fullPath)
	}
}

func handlePath(fullPath string, options *cliOptions) {
	defer wg.Done()

	matched, err := filepath.Match(options.file, path.Base(fullPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	if matched {
		wg.Add(1)
		go processPath(fullPath, options)
	}
}
