package cli

import (
	"fmt"

	"github.com/fatih/color"
)

type Output struct {
	Console chan string
	Closed  chan bool
}

func (this *Output) Process() {
	var total int
	for path := range this.Console {
		total++
		color.Magenta(path)
	}
	fmt.Printf("%d file(s) updated\n", total)
	this.Closed <- true
}
