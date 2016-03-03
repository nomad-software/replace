package cli

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type Options struct {
	Dir  string
	From string
	To   string
	File string
	Help bool
}

func (this *Options) Valid() bool {
	return this.From != "" && this.To != ""
}

func (this *Options) Echo() {
	options := color.CyanString("replacing:   ")
	options += color.GreenString("%s\n", this.From)
	options += color.CyanString("with:        ")
	options += color.GreenString("%s\n", this.To)
	options += color.CyanString("in files:    ")
	options += color.GreenString("%s\n", this.File)
	options += color.CyanString("starting in: ")
	options += color.GreenString("%s\n", this.Dir)
	fmt.Print(options)
}

func (this *Options) Parse() {
	flag.StringVar(&this.Dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&this.From, "from", "", "The text to replace.")
	flag.StringVar(&this.To, "to", "", "The replacement text.")
	flag.StringVar(&this.File, "file", "*", "The glob file pattern to match.")
	flag.BoolVar(&this.Help, "help", false, "Show help.")
	flag.Parse()
	dir, _ := homedir.Expand(this.Dir)
	this.Dir = dir
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
