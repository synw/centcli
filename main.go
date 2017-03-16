package main

import (
	"flag"
	"github.com/abiosoft/ishell"
	"github.com/synw/centcli/libcentcli/state"
	"github.com/synw/centcli/libcentcli/cmd"
)


var shell = ishell.New()
var use = flag.String("s", "unset", "Use server: -s=server_name")

func main() {
	flag.Parse()
	trace := state.InitState()
	if trace != nil {
		trace.Printc()
	}
	if *use != "unset" {
		trace = state.SetServer(*use)
		if trace != nil {
			trace.Printc()
		}
		trace = state.InitServer()
		if trace != nil {
			trace.Printc()
		}
	}
	shell.SetHomeHistoryPath(".ishell_history")
	// commands
	shell = cmd.GetCmds(shell)
	// start shell
    shell.Start()
}
