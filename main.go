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
)

func processPath(fullPath string, options *cliOptions) {
	defer wg.Done()

	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	newContents := strings.Replace(string(contents), options.from, options.to, -1)
	if newContents != string(contents) {

		err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(fullPath)
	}
}

func handlePath(fullPath string, options *cliOptions) {
	defer wg.Done()

	matched, err := filepath.Match(options.file, path.Base(fullPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if matched {
		wg.Add(1)
		go processPath(fullPath, options)
	}
}

type cliOptions struct {
	dir  string
	from string
	to   string
	file string
	help bool
}

var wg sync.WaitGroup

func main() {

	var options cliOptions

	flag.StringVar(&options.dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&options.from, "from", "", "The text to replace.")
	flag.StringVar(&options.to, "to", "", "The replacement text.")
	flag.StringVar(&options.file, "file", "*", "The glob file pattern to match.")
	flag.BoolVar(&options.help, "help", false, "Show help.")
	flag.Parse()

	if (options.from == "" && options.to == "") || options.help {
		flag.Usage()

	} else {
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
			fmt.Fprintln(os.Stderr, err)
			return
		}

		wg.Wait()
	}
}
