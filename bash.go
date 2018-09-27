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
  while read -r line; do
		if [[ -n $line ]]; then
			if [[ "$line" =~ \ |\' ]]; then  # if has spaces
				line="\"$line\""
			fi
      COMPREPLY+=("$line")
    fi
  done <<< "$output"
  if [[ ${#COMPREPLY[@]} -eq 1 ]]; then
    COMPREPLY[0]="${COMPREPLY[0]} "
  fi
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
