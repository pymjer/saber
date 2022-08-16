package envcmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"prolion.top/saber/internal/base"
	"prolion.top/saber/internal/cfg"
)

var CmdEnv = &base.Command{
	UsageLine: "saber env [-u] [-w] [var ...]",
	Short:     "print saber environment information",
	Long: `
Env prints Saber environment information.

The -u flag requires one or more arguments and unsets
the default setting for the named environment variables,
if one has been set with 'saber env -w'.

The -w flag requires one or more arguments of the
form NAME=VALUE and changes the default settings
of the named environment variables to the given values.
	`,
}

func init() {
	CmdEnv.Run = runEnv
}

var (
	envU = CmdEnv.Flag.Bool("u", false, "unset variable")
	envW = CmdEnv.Flag.Bool("w", false, "set variable")
)

func runEnv(ctx context.Context, cmd *base.Command, args []string) {
	if *envW && *envU {
		base.Fatalf("saber env: cannot use -u with -w")
	}

	if *envW {
		runEnvW(args)
		return
	}

	if *envU {
		runEnvU(args)
		return
	}

	env := cfg.CmdEnv
	if len(args) > 0 {
		for _, name := range args {
			fmt.Printf("%s\n", findEnv(env, name))
		}
		return
	}

	PrintEnv(os.Stdout, env)
}

func MkEnv() []cfg.EnvVar {
	envFile, _ := cfg.EnvFile()
	env := []cfg.EnvVar{
		{Name: "SABERENV", Value: envFile},
		{Name: "SABERVERSION", Value: "1.0"},
		{Name: "ZK", Value: cfg.Getenv("ZK")},
	}

	return env
}

func PrintEnv(w io.Writer, env []cfg.EnvVar) {
	for _, e := range env {
		fmt.Fprintf(w, "%s=\"%s\"\n", e.Name, e.Value)
	}
}

func findEnv(env []cfg.EnvVar, name string) string {
	for _, e := range env {
		if e.Name == name {
			return e.Value
		}
	}
	return ""
}

func runEnvW(args []string) {
	if len(args) == 0 {
		base.Fatalf("saber env -w: no KEY=VALUE arguments given")
	}

	add := make(map[string]string)
	for _, arg := range args {
		i := strings.Index(arg, "=")
		if i < 0 {
			base.Fatalf("saber env -w: arguments must be KEY=VALUE: invalid argument: %s", arg)
		}
		key, val := arg[:i], arg[i+1:]
		add[key] = val
	}
	updateEnvFile(add, nil)
}

func runEnvU(args []string) {
	if len(args) == 0 {
		base.Fatalf("saber env -u: no arguments given")
	}
	del := make(map[string]bool)
	for _, arg := range args {
		del[arg] = true
	}
	updateEnvFile(nil, del)
}

func updateEnvFile(add map[string]string, del map[string]bool) {
	file, err := cfg.EnvFile()
	if file == "" {
		base.Fatalf("saber env: cannot find go env config: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil && (!os.IsNotExist(err) || len(add) == 0) {
		base.Fatalf("saber env: reading saber env config: %v", err)
	}

	lines := strings.SplitAfter(string(data), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	} else {
		lines[len(lines)-1] += "\n"
	}
	// Delete all but last copy of any duplicated variables,
	// since the last copy is the one that takes effect.
	prev := make(map[string]int)
	for l, line := range lines {
		if key := lineToKey(line); key != "" {
			if p, ok := prev[key]; ok {
				lines[p] = ""
			}
			prev[key] = l
		}
	}

	// Add variables (go env -w). Update existing lines in file if present, add to end otherwise.
	for key, val := range add {
		if p, ok := prev[key]; ok {
			lines[p] = key + "=" + val + "\n"
			delete(add, key)
		}
	}
	for key, val := range add {
		lines = append(lines, key+"="+val+"\n")
	}

	// Delete requested variables (go env -u).
	for key := range del {
		if p, ok := prev[key]; ok {
			lines[p] = ""
		}
	}

	// Sort runs of KEY=VALUE lines
	// (that is, blocks of lines where blocks are separated
	// by comments, blank lines, or invalid lines).
	start := 0
	for i := 0; i <= len(lines); i++ {
		if i == len(lines) || lineToKey(lines[i]) == "" {
			sortKeyValues(lines[start:i])
			start = i + 1
		}
	}

	data = []byte(strings.Join(lines, ""))
	err = os.WriteFile(file, data, 0666)
	if err != nil {
		// Try creating directory.
		os.MkdirAll(filepath.Dir(file), 0777)
		err = os.WriteFile(file, data, 0666)
		if err != nil {
			base.Fatalf("saber env: writing go env config: %v", err)
		}
	}
}

// lineToKey returns the KEY part of the line KEY=VALUE or else an empty string.
func lineToKey(line string) string {
	i := strings.Index(line, "=")
	if i < 0 || strings.Contains(line[:i], "#") {
		return ""
	}
	return line[:i]
}

// sortKeyValues sorts a sequence of lines by key.
// It differs from sort.Strings in that keys which are GOx where x is an ASCII
// character smaller than = sort after GO=.
// (There are no such keys currently. It used to matter for GO386 which was
// removed in Go 1.16.)
func sortKeyValues(lines []string) {
	sort.Slice(lines, func(i, j int) bool {
		return lineToKey(lines[i]) < lineToKey(lines[j])
	})
}
