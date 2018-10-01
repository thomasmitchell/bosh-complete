package main

import (
	"fmt"
	"os"

	"text/template"
)

var bashSource = fmt.Sprintf(`
_bosh_comp() {
	local output="$({{.Executable}} complete {{.Debug}} -- ${COMP_WORDS[@]::$COMP_CWORD} "${COMP_WORDS[$COMP_CWORD]}")"
	COMPREPLY=()
	local TMPIFS="$IFS"
	IFS=''
  while read -r line; do
		if [[ -n "$line" ]]; then
      COMPREPLY+=("$line")
    fi
	done <<< "$output"
	IFS="$TMPIFS"
}

complete -o nospace -F _bosh_comp {{.Bosh}}
`)

func doBashSource() {
	tmpl := template.Must(template.New("bash_source").Parse(bashSource))
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
