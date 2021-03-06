package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is a color friendly pipe.
	Stdout = colorable.NewColorableStdout()

	// Stderr is a color friendly pipe.
	Stderr = colorable.NewColorableStderr()
)

type Output struct {
	Console chan string
	Closed  chan bool
}

func (this *Output) Start() {
	var total int
	for path := range this.Console {
		total++
		color.Magenta(path)
	}
	fmt.Printf("%d file(s) updated\n", total)
	this.Closed <- true
}
