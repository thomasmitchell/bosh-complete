package main

import (
	"fmt"
	"os"

	"text/template"
)

var zshSource = fmt.Sprintf(`
#compdef {{.Bosh}}
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

_bosh_comp() {
	local output="$({{.Executable}} complete {{.Debug}} -- ${COMP_WORDS[@]:0:$COMP_CWORD} "${COMP_WORDS[$COMP_CWORD]}")"
	COMPREPLY=()
	IFS=''
  while read -r line; do
		if [[ -n "$line" ]]; then
      COMPREPLY+=("$line")
    fi
  done <<< "$output"
}

complete -o nospace -F _bosh_comp {{.Bosh}}
`)

func doZshSource() {
	tmpl := template.Must(template.New("bash_source").Parse(zshSource))
	me, err := os.Executable()
	debug := ""
	if opts.Debug {
		debug = "--debug"
	}
	if err != nil {
		panic("Could not determine executable location")
	}
	tmpl.Execute(os.Stdout, struct {
		Executable string
		Bosh       string
		Debug      string
	}{
		Executable: me,
		Bosh:       "bosh",
		Debug:      debug,
	})
}
