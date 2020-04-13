package main

import (
	"flag"
	"log"

	"github.com/afking/starlarkrepl"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

func run() error {
	flag.Parse()

	thread := &starlark.Thread{Load: repl.MakeLoad()}
	globals := make(starlark.StringDict)
	options := starlarkrepl.Options{AutoComplete: true}

	return starlarkrepl.Run(thread, globals, options)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
