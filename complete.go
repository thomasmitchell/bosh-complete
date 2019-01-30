package main

import (
	"fmt"
	"os"
	"strings"
)

//Keeps command completion from automatically adding a space to move
// to the next token
var dontAddSpace bool

//Keeps command completion from filtering out command completions that
// don't start with the current token
var dontFilterPrefix bool

type compContext struct {
	CurrentToken string
	Command      string
	Args         []string
	//Long flag string of the current flag that is being completed
	CurrentFlag string
	//Long flag string to value(s)
	Flags    map[string][]string
	Switches []string
}

func (c *compContext) InsertIfEnvvar(envvar, flag string) {
	val := os.Getenv(envvar)
	if val != "" {
		c.Flags[flag] = append(c.Flags[flag], val)
	}
}

func (c compContext) Complete() ([]string, error) {
	var compFn compFunc

	log.Write("Current Token: %s", c.CurrentToken)

	//determine what we're completing
	if c.CurrentFlag != "" {
		log.Write("Checking current flag: %s", c.CurrentFlag)
		flag, found := flags[c.CurrentFlag]
		if found {
			log.Write("Completing flag value for flag %s", c.CurrentFlag)
			compFn = flag.Complete
		}
	} else if strings.HasPrefix(c.CurrentToken, "-") {
		log.Write("Completing flag names")
		compFn = compFlagNames
	} else if c.Command == "" {
		log.Write("Completing command names")
		compFn = compCommandNames
	} else {
		if cmd, found := commands.Find(c.Command); found {
			if len(c.Args) < len(cmd.Args) {
				position := len(c.Args)
				log.Write("Completing positional arg for command `%s' position %d (0 indexed)", c.Command, position)
				compFn = cmd.Args[position]
			}
		}
	}

	if compFn == nil {
		log.Write("No completion registered")
		compFn = compNoop
	}

	candidates, err := compFn(c)
	if err != nil {
		return nil, err
	}

	//log.Write("Completion candidates: \n---START---\n%s\n---END---\n", strings.Join(candidates, "\n"))

	ret := []string{}
	for _, val := range candidates {
		if strings.ContainsAny(val, " \t\n\r") {
			val = fmt.Sprintf(`"%s"`, val)
		}
		if dontFilterPrefix || strings.HasPrefix(val, c.CurrentToken) {
			ret = append(ret, val)
		}
	}

	if len(ret) == 1 && !dontAddSpace {
		ret[0] = fmt.Sprintf("%s ", ret[0])
	}

	log.Write("Completion return: \n---START---\n%s\n---END---\n", strings.Join(ret, "\n"))
	return ret, nil
}

type compFunc func(compContext) ([]string, error)

func doComplete(boshArgs []string) {
	log.Write("in complete")
	argsString := ""
	if len(boshArgs) > 0 {
		argsString = fmt.Sprintf(`'%s'`, strings.Join(boshArgs, `', '`))
	}
	log.Write("Bosh args: [%s]", argsString)

	insertGlobalFlags()
	commands.Populate()

	compContext := parseContext(boshArgs)
	results, err := compContext.Complete()
	if err != nil {
		log.Write("Completion error: %s", err.Error())
		return
	}

	response := strings.Join(results, "\n")
	fmt.Print(response)
}

func parseContext(args []string) compContext {
	if len(args) < 2 {
		os.Exit(0)
	}

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") && strings.Contains(args[i], "=") {
			spl := strings.SplitN(args[i], "=", 2)
			log.Write("Split `%s' into parts `%s' and `%s'", args[i], spl[0], spl[1])
			before := args[:i]
			rest := []string{spl[0], spl[1]}
			if i < len(args)-1 {
				rest = append(rest, args[i+1:]...)
			}
			tmp := append(before, rest...)
			args = tmp
			i++
		}
	}

	if len(args) > 0 {
		log.Write("args after split: [`%s']", strings.Join(args, "', `"))
	}

	ret := compContext{
		CurrentToken: args[len(args)-1],
		Flags:        map[string][]string{},
	}

	//loop over all but last token - the last one is the token
	// we're suggesting changes to.
	for i := 0; i < len(args)-1; i++ {
		token := args[i]
		if strings.HasPrefix(token, "-") && ret.CurrentFlag == "" {
			//Check if value or not
			f := flags[token]
			ret.CurrentFlag = "--" + f.Long
			if f.Complete == nil {
				log.Write("current flag is switch: %s", ret.CurrentFlag)
				ret.Switches = append(ret.Switches)
				ret.CurrentFlag = ""
			}
		} else {
			if ret.CurrentFlag != "" {
				//This is the value to a flag
				ret.Flags[ret.CurrentFlag] = append(ret.Flags[ret.CurrentFlag], token)
				ret.CurrentFlag = ""
			} else if ret.Command == "" {
				//This is a command name
				if cmd, found := commands.Find(token); found {
					ret.Command = token
					cmd.InsertFlags()
				}
			} else {
				//This is a positional argument
				ret.Args = append(ret.Args, token)
			}
		}
	}

	//Flags override environment variables, so put in env vars last... they would
	// become the second flag value, which is typically ignored in the code
	ret.InsertIfEnvvar("BOSH_ENVIRONMENT", "--environment")
	ret.InsertIfEnvvar("BOSH_DEPLOYMENT", "--deployment")
	ret.InsertIfEnvvar("BOSH_CLIENT", "--client")
	ret.InsertIfEnvvar("BOSH_CLIENT_SECRET", "--client-secret")
	ret.InsertIfEnvvar("BOSH_NON_INTERACTIVE", "--non-interactive")
	ret.InsertIfEnvvar("BOSH_CA_CERT", "--ca-cert")

	return ret
}
