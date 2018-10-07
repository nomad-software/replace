package file

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/nomad-software/replace/cli"
)

const MAX_NUMBER_OF_WORKERS = 100

type WorkerQueue struct {
	Closed chan bool
	Group  *sync.WaitGroup
	Input  chan UnitOfWork
	Output *cli.Output
}

type UnitOfWork struct {
	File string
	From string
	To   string
}

func (this *WorkerQueue) Start() {
	go this.Output.Start()

	life := make(chan bool)

	for i := 0; i <= MAX_NUMBER_OF_WORKERS; i++ {
		go this.worker(life)
	}

	for i := 0; i <= MAX_NUMBER_OF_WORKERS; i++ {
		<-life
	}

	close(this.Output.Console)
	<-this.Output.Closed

	this.Closed <- true
}

func (this *WorkerQueue) Close() {
	close(this.Input)
	<-this.Closed
}

func (this *WorkerQueue) worker(death chan<- bool) {
	for work := range this.Input {

		contents, err := ioutil.ReadFile(work.File)
		if err != nil {
			fmt.Fprintln(cli.Stderr, color.RedString(err.Error()))
			this.Group.Done()
			continue
		}

		newContents := strings.Replace(string(contents), work.From, work.To, -1)

		if newContents != string(contents) {
			err = ioutil.WriteFile(work.File, []byte(newContents), 0)
			if err != nil {
				fmt.Fprintln(cli.Stderr, color.RedString(err.Error()))
				this.Group.Done()
				continue
			}

			this.Output.Console <- work.File
		}

		this.Group.Done()
	}

	death <- true
}
