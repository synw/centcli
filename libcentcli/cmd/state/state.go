package state

import (
	"errors"
	"github.com/abiosoft/ishell"
	"github.com/acmacalister/skittles"
	"github.com/synw/terr"
	"github.com/synw/centcli/libcentcli/state"
)


func Using() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"using",
        Help: 	"Server actually in use",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	} 
			ctx.Println("Using server", state.Server.Name)
        },
	}
	return command
}

func Use() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"use",
        Help: 	"Use server: use server_name",
        Func: 	func(ctx *ishell.Context) {
        			var is_err bool = false
					if len(ctx.Args) == 0 {
						err := terr.Err("missing server name")
						ctx.Println(err.Error())
						is_err = true
					}
					if len(ctx.Args) > 1 {
						err := terr.Err("please use only one server at the time")
						ctx.Println(err.Error())
						is_err = true
					}
					server_name := ctx.Args[0]
					_, trace  := state.ServerExists(server_name)
					if trace != nil {
						err := errors.New("Can not find server")
						trace = terr.Add("cmd.state.Use", err, trace)
						ctx.Println(trace.Formatc())
						is_err = true
					}
					if is_err == false {
						ctx.Println("Set state")
						old := state.Server
						state.Server = state.Servers[server_name]
						// init cli and check server
						trace = state.InitServer()
			        	if trace != nil {
			        		err := errors.New("can not connect to websockets server: check your config")
			        		trace := terr.Add("cmd.state.Use", err, trace)
			        		ctx.Println(trace.Formatc())
			        		state.Server = old
			        		return
			        	} 
						msg := "Using server "+server_name
						ctx.Println(msg)
					}
				},
    }
	return command
}
