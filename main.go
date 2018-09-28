package main

import (
	"fmt"

	cli "github.com/jhunt/go-cli"
	"github.com/thomasmmitchell/bosh-complete/version"
)

var opts options

type options struct {
	Debug      bool     `cli:"-d, --debug"`
	Complete   struct{} `cli:"complete"`
	BashSource struct{} `cli:"bash-source"`
	ZshSource  struct{} `cli:"zsh-source"`
	Version    struct{} `cli:"version"`
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
	case "zsh-source":
		//For my weird friends Nic Williams and Long Nguyen
		doZshSource()
	case "version":
		doVersion()
	default:
		panic("Unknown command: " + command)
	}
}

func doVersion() {
	fmt.Println(version.Version)
}
