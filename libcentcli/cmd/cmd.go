package cmd

import (
	"github.com/abiosoft/ishell"
	cmd_state "github.com/synw/centcli/libcentcli/cmd/state"
	"github.com/synw/centcli/libcentcli/cmd/stats"
	"github.com/synw/centcli/libcentcli/cmd/actions"
	"github.com/synw/centcli/libcentcli/cmd/chans"
)


func GetCmds(shell *ishell.Shell) *ishell.Shell {
	shell.AddCmd(cmd_state.Using())
	shell.AddCmd(cmd_state.Use())
	shell.AddCmd(stats.Channels())
	shell.AddCmd(stats.Count())
	shell.AddCmd(stats.Stats())
	shell.AddCmd(stats.Stat())
	shell.AddCmd(actions.Publish())
	shell.AddCmd(actions.Listen())
	shell.AddCmd(actions.Stop())
	shell.AddCmd(chans.Presence())
	shell.AddCmd(chans.History())
	return shell
}
