package main

func insertGlobalFlags() {
	insertFlag(flag{Long: "version", Short: 'v'})
	insertFlag(flag{Long: "config", Complete: compFiles})
	insertFlag(flag{Long: "environment", Short: 'e', Complete: compEnvAliases})
	insertFlag(flag{Long: "ca-cert", Complete: compFiles})
	insertFlag(flag{Long: "sha2"})
	insertFlag(flag{Long: "parallel", Complete: compNoop})
	insertFlag(flag{Long: "client", Complete: compNoop})
	insertFlag(flag{Long: "client-secret", Complete: compNoop})
	//TODO: --deployment
	insertFlag(flag{Long: "deployment", Complete: compNoop})
	insertFlag(flag{Long: "column", Complete: compNoop})
	insertFlag(flag{Long: "json"})
	insertFlag(flag{Long: "tty"})
	insertFlag(flag{Long: "no-color"})
	insertFlag(flag{Long: "non-interactive", Short: 'n'})
	insertFlag(flag{Long: "help", Short: 'h'})
}
