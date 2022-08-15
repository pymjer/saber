package main

import (
	"context"
	"flag"
	"os"

	"prolion.top/saber/bigmap"
	"prolion.top/saber/findfiles"
	"prolion.top/saber/internal/base"
	"prolion.top/saber/internal/help"
)

func init() {
	base.Saber.Commands = []*base.Command{
		findfiles.CmdFindFiles,
		bigmap.CmdBigMap,
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
	}
	// fmt.Printf("cmd:%s \n", cmd)
	// switch cmd {
	// case "findfiles", "1":
	// 	findfiles.FindFilesMain()
	// case "bigmap", "2":
	// 	bigmap.BigMapMain()
	// case "hbaseutils", "3":
	// 	hbaseutils.HBaseUtilMain()
	// default:
	// 	log.Fatalf("未知命令: %s", cmd)
	// }
}

func invoke(cmd *base.Command, args []string) {
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
