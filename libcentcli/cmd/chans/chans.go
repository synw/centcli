package chans

import (
	"fmt"
	"strconv"
	"github.com/abiosoft/ishell"
	"github.com/acmacalister/skittles"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/centcli/libcentcli/state"
)


func History() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"history",
        Help: 	"History for channel: ex: history channel_name",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	}
        	if len(ctx.Args) == 0 {
				err := terr.Err("No channel provided: ex: history channel_name")
				ctx.Println(err.Error())
				return
			}
			channel := ctx.Args[0]
        	history, err := state.Cli.Http.History(channel)
    		if err != nil {
    			trace := terr.New("cmd.chans.History", err)
    			ctx.Println(trace.Formatc())
    			return
    		}
    		ctx.Println("History for channel "+channel)
    		for i, rmsg := range(history) {
    			msg, err := centcom.DecodeHttpMsg(&rmsg)
    			if err != nil {
    				trace := terr.New("cmd.chans.History", err)
    				ctx.Println(trace.Formatc())
    				return
    			}
    			out := strconv.Itoa(i)+" "+fmt.Sprintf("%s", msg.Payload)
    			ctx.Println(out)
    		}
        },
	}
	return command
}

func Presence() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"presence",
        Help: 	"Presence for channel: ex: presence channel_name",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	}
        	if len(ctx.Args) == 0 {
				err := terr.Err("No channel provided: ex: presence channel_name")
				ctx.Println(err.Error())
				return
			}
			channel := ctx.Args[0]
        	pres, err := state.Cli.Http.Presence(channel)
    		if err != nil {
    			trace := terr.New("cmd.stats.Channels", err)
    			ctx.Println(trace.Formatc())
    		}
        	ctx.Println("Presence for channel", channel)
        	var users string
        	var num int
    		for _, rmsg := range(pres) {
    			user := rmsg.User
    			if err != nil {
    				trace := terr.New("cmd.chans.Presence", err)
    				ctx.Println(trace.Formatc())
    				return
    			}
    			num++
    			users = users+user+" "
    		}
    		if users == "" {
    			ctx.Println("No users in channel")
    		} else {
    			msg := "Found "+strconv.Itoa(num)+" users: "+users
    			ctx.Println(msg)
    		}
        },
	}
	return command
}
