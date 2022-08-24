package main

import (
	"context"
	"flag"
	"os"

	"prolion.top/saber/bigmap"
	"prolion.top/saber/findfiles"
	"prolion.top/saber/hbaseutils"
	"prolion.top/saber/internal/base"
	"prolion.top/saber/internal/cfg"
	"prolion.top/saber/internal/envcmd"
	"prolion.top/saber/internal/help"
	"prolion.top/saber/replace"
)

func init() {
	base.Saber.Commands = []*base.Command{
		findfiles.CmdFindFiles,
		bigmap.CmdBigMap,
		hbaseutils.CmdHBaseUtil,
		envcmd.CmdEnv,
		replace.CmdReplaceFile,
	}
}

func main() {
	flag.Usage = base.Usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		base.Usage()
		return
	}

	if args[0] == "help" {
		help.Help(os.Stdout, args[1:])
		return
	}

	for bigCmd := base.Saber; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if !cmd.Runnable() {
				continue
			}
			invoke(cmd, args)
			base.Exit()
			return
		}
		base.Fatalf("no command[%s] found.", args[0])
		return
	}
}

func invoke(cmd *base.Command, args []string) {

	cfg.CmdEnv = envcmd.MkEnv()
	for _, env := range cfg.CmdEnv {
		if os.Getenv(env.Name) != env.Value {
			os.Setenv(env.Name, env.Value)
		}
	}

	cmd.Flag.Usage = func() { cmd.Usage() }
	cmd.Flag.Parse(args[1:])
	args = cmd.Flag.Args()

	cmd.Run(context.Background(), cmd, args)
}

// golang allowing multiple init in one file
func init() {
	base.Usage = mainUsage
}

func mainUsage() {
	help.PrintUsage(os.Stderr, base.Saber)
	os.Exit(2)
}
