// base包定义了命令的基础工具类
package base

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// A Command is an implementation of a tools
type Command struct {

	// Run runs the command.
	Run func(ctx context.Context, cmd *Command, args []string)

	// Usage Line is the one-line usage message.
	UsageLine string

	// Short is the short description shown in the 'saber help' output.
	Short string

	// Long is the long message shown in the 'saber help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// Commands lists
	Commands []*Command
}

// LongName returns the command's long name: all the words in the usage line between "saber" and a flag
func (c *Command) LongName() string {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}
	if name == "saber" {
		return ""
	}
	return strings.TrimPrefix(name, "saber ")
}

// Name returns the command's short name: the last word in the usage line before a flag or argument.
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// Runnable reports whether the command can be run;
func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "Run 'saber help %s' for details.\n", c.LongName())
	SetExitStatus(2)
	Exit()
}

func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
	Exit()
}

func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	SetExitStatus(1)
}

var exitStatus = 0
var exitMu sync.Mutex

func SetExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func GetExitStatus() int {
	return exitStatus
}

var atExitFuncs []func()

func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

var Saber = &Command{
	UsageLine: "saber",
	Long:      `saber is a collection tools`,
}

var Usage func()
