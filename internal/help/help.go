package help

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"prolion.top/saber/internal/base"
)

func Help(w io.Writer, args []string) {
	cmd := base.Saber
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name() == arg {
				cmd = sub
				continue Args
			}
		}
		helpSuccess := "saber help"
		if i > 0 {
			helpSuccess += " " + strings.Join(args[:i], " ")
		}
		fmt.Fprintf(os.Stderr, "saber help %s: unknown help topic. Run '%s'. \n", strings.Join(args, " "), helpSuccess)
		base.SetExitStatus(2) // failed at 'saber help cmd'
		base.Exit()
	}

	if len(cmd.Commands) > 0 {
		PrintUsage(os.Stdout, cmd)
	} else {
		tmpl(os.Stdout, helpTemplate, cmd)
	}

}

var usageTemplate = `{{.Long | trim}}

Usage:

	{{.UsageLine}} <command> [arguments]

The commands are:
{{range .Commands}}{{if or (.Runnable) .Commands}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "saber help{{with .LongName}} {{.}}{{end}} <command>" for more information about a command.
{{if eq (.UsageLine) "saber"}}
Additional help topics:
{{range .Commands}}{{if and (not .Runnable) (not .Commands)}}
	{{.Name | printf "%-15s"}} {{.Short}}{{end}}{{end}}

Use "saber help{{with .LongName}} {{.}}{{end}} <topic>" for more information about that topic.
{{end}}
`

var helpTemplate = `{{if .Runnable}}usage: {{.UsageLine}}

{{end}}{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	t.Parse(text)
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func PrintUsage(w io.Writer, cmd *base.Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, cmd)
	bw.Flush()
}
