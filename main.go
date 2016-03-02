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

func (this *cliOptions) echo() {
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

func (this *cliOptions) parse() {
	flag.StringVar(&this.dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&this.from, "from", "", "The text to replace.")
	flag.StringVar(&this.to, "to", "", "The replacement text.")
	flag.StringVar(&this.file, "file", "*", "The glob file pattern to match.")
	flag.BoolVar(&this.help, "help", false, "Show help.")
	flag.Parse()
	dir, _ := homedir.Expand(this.dir)
	this.dir = dir
}

func (this *cliOptions) usage() {
	var banner string = ` ____            _
|  _ \ ___ _ __ | | __ _  ___ ___
| |_) / _ \ '_ \| |/ _' |/ __/ _ \
|  _ <  __/ |_) | | (_| | (_|  __/
|_| \_\___| .__/|_|\__'_|\___\___|
          |_|

`
	color.Cyan(banner)
	flag.Usage()
}

type fileHandler struct {
	options *cliOptions
	wg      sync.WaitGroup
	total   int64
	output  chan string
	done    chan bool
}

func (this *fileHandler) init(options *cliOptions) {
	this.options = options
	this.output = make(chan string)
	this.done = make(chan bool)
}

func (this *fileHandler) processOutput() {
	for path := range this.output {
		this.total++
		color.Magenta(path)
	}
	fmt.Printf("%d file(s) updated\n", this.total)
	this.done <- true
}

func (this *fileHandler) handlePath(fullPath string) {
	defer this.wg.Done()

	matched, err := filepath.Match(this.options.file, path.Base(fullPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	if matched {
		this.wg.Add(1)
		go this.processPath(fullPath)
	}
}

func (this *fileHandler) processPath(fullPath string) {
	defer this.wg.Done()

	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	newContents := strings.Replace(string(contents), this.options.from, this.options.to, -1)
	if newContents != string(contents) {

		err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		this.output <- fullPath
	}
}

func (this *fileHandler) walk() error {
	return filepath.Walk(this.options.dir, func(fullPath string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		this.wg.Add(1)
		go this.handlePath(fullPath)

		return nil
	})
}

func main() {

	var options cliOptions
	options.parse()

	var files fileHandler
	files.init(&options)
	go files.processOutput()

	if (!options.valid()) || options.help {
		options.usage()

	} else {
		options.echo()
		err := files.walk()

		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			return
		}

		files.wg.Wait()

		close(files.output)
		<-files.done
	}
}
