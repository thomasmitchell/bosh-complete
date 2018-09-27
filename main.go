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
		log.Write("Could not init cli parser: %s", err.Error())
	}

	if opts.Debug {
		log.TurnOn()
	}

	switch command {
	case "complete":
		doComplete(args)
	case "bash-source":
		doBashSource()
	default:
		panic("Unknown command: " + command)
	}
}
