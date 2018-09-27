package main

func insertGlobalFlags() {
	insertFlag(flag{Long: "environment", Short: 'e', Complete: compEnvAliases, TakesValue: true})
	insertFlag(flag{Long: "version", Short: 'v'})
	insertFlag(flag{Long: "sha2"})
	insertFlag(flag{Long: "json"})
	insertFlag(flag{Long: "tty"})
	insertFlag(flag{Long: "no-color"})
	insertFlag(flag{Long: "non-interactive", Short: 'n'})
	insertFlag(flag{Long: "help", Short: 'h'})
}
