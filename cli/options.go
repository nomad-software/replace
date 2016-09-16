package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

const (
	DEFAULT_DIRECTORY = "."
	DEFAULT_GLOB      = "*"
)

type Options struct {
	Case   bool
	Dir    string
	File   string
	From   string
	Help   bool
	Ignore string
	To     string
}

func ParseOptions() Options {
	var opt Options

	flag.BoolVar(&opt.Case, "case", false, "Use to switch to case sensitive matching.")
	flag.StringVar(&opt.Dir, "dir", DEFAULT_DIRECTORY, "The directory to traverse.")
	flag.StringVar(&opt.File, "file", DEFAULT_GLOB, "The glob file pattern to match.")
	flag.StringVar(&opt.From, "from", "", "The text to replace.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.StringVar(&opt.To, "to", "", "The replacement text.")
	flag.Parse()

	opt.Dir, _ = homedir.Expand(opt.Dir)

	return opt
}

func (this *Options) Valid() bool {

	if this.From == "" {
		fmt.Fprintln(os.Stderr, color.RedString("From cannot be empty."))
		return false
	}

	if this.To == "" {
		fmt.Fprintln(os.Stderr, color.RedString("To cannot be empty."))
		return false
	}

	err := this.compiles(this.Ignore)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("ignore pattern: %s", err.Error()))
		return false
	}

	return true
}

func (this *Options) Echo() {
	output := color.CyanString("replacing:   ")
	output += color.GreenString("%s\n", this.From)
	output += color.CyanString("with:        ")
	output += color.GreenString("%s\n", this.To)
	output += color.CyanString("in files:    ")
	output += color.GreenString("%s\n", this.File)
	output += color.CyanString("starting in: ")
	output += color.GreenString("%s\n", this.Dir)

	if this.Ignore != "" {
		output += color.CyanString("ignoring:    ")
		output += color.GreenString("%s\n", this.Ignore)
	}

	fmt.Print(output)
}

func (this *Options) Usage() {
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

func (this *Options) compiles(pattern string) (err error) {
	if this.Case {
		_, err = regexp.Compile(pattern)
	} else {
		_, err = regexp.Compile("(?i)" + pattern)
	}

	return err
}
