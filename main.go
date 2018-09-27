package main

import (
	cli "github.com/jhunt/go-cli"
)

var opts options

type options struct {
	Debug      bool     `cli:"-d, --debug"`
	Complete   struct{} `cli:"complete"`
	BashSource struct{} `cli:"bash-source"`
}

func main() {
	command, args, err := cli.Parse(&opts)
	if err != nil {
		panic("Could not init cli parser: " + err.Error())
	}

	if opts.Debug {
		log.TurnOn()
	}

	log.Write("")

	switch command {
	case "complete":
		doComplete(args)
	case "bash-source":
		doBashSource()
	default:
		panic("Unknown command: " + command)
	}
}
