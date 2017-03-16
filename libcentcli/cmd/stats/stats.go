package stats

import (
	"fmt"
	"errors"
	"strconv"
	"strings"
	"github.com/centrifugal/gocent"
	"github.com/abiosoft/ishell"
	"github.com/acmacalister/skittles"
	"github.com/synw/terr"
	"github.com/synw/centcli/libcentcli/state"
)


func Channels() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"chans",
        Help: 	"Channels on the server",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        	} else {
        		channels, err := state.Cli.Http.Channels()
        		if err != nil {
        			trace := terr.New("cmd.stats.Channels", err)
        			ctx.Println(trace.Formatc())
        		}
	        	ctx.Println("Available channels:", channels)
        	}
        },
	}
	return command
}

func Stat() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"stat",
        Help: 	"Get a server statistic: ex: stat node_num_clients",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	}
        	if len(ctx.Args) == 0 {
				err := terr.Err("No arguments provided: ex: stat node_num_clients")
				ctx.Println(err.Error())
				return
			} else if len(ctx.Args) > 1 {
				err := terr.Err("Only one argument is allowed: ex: stat node_num_clients")
				ctx.Println(err.Error())
				return
			}
			stat := ctx.Args[0]
			stats, err := state.Cli.Http.Stats()
			if err != nil {
				tr := terr.Err(err.Error())
				ctx.Println(tr.Error())
				return
			}
			for _, node := range(stats.Nodes) {
				msg := msgForNode(&node)
				ctx.Println(msg)
				nmsg, trace := getStat(&node, stat)
				if trace != nil {
					ctx.Println(trace.Formatc())
					return
				}
				ctx.Println(nmsg)
			}
        },
	}
	return command
}

func Stats() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"stats",
        Help: 	"Server statistics",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	}
        	if len(ctx.Args) == 0 {
				err := terr.Err("No arguments provided: ex: stats node")
				ctx.Println(err.Error())
				return
			} else if len(ctx.Args) > 1 {
				err := terr.Err("Only one argument is allowed: ex: stats node")
				ctx.Println(err.Error())
				return
			}
    		stats, err := state.Cli.Http.Stats()
    		if err != nil {
    			trace := terr.New("cmd.stats.Stats", err)
    			ctx.Println(trace.Formatc())
    		}
    		if ctx.Args[0] == "all" {
    			for _, node := range(stats.Nodes) {
	    			msg := statsForNode(&node, "all")
	    			ctx.Println(msg)
    			}
    		} else if ctx.Args[0] == "node" {
    			for _, node := range(stats.Nodes) {
    				msg := statsForNode(&node, "node")
    				ctx.Println(msg)
    			}
    		} else if ctx.Args[0] == "http" {
    			for _, node := range(stats.Nodes) {
    				msg := statsForNode(&node, "http")
    				ctx.Println(msg)
    			}
    		} else if ctx.Args[0] == "client" {
    			for _, node := range(stats.Nodes) {
    				msg := statsForNode(&node, "client")
    				ctx.Println(msg)
    			}
    		}
        },
	}
	return command
}

func Count() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"count",
        Help: 	"Count things on the server: ex: count chans",
        Func: 	func(ctx *ishell.Context) {
        	if state.Server == nil {
        		ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
        		return
        	}
    		if len(ctx.Args) == 0 {
				err := terr.Err("missing item to count: ex: count chans")
				ctx.Println(err.Error())
				return
			}
			if ctx.Args[0] == "chans" {
				channels, err := state.Cli.Http.Channels()
	    		if err != nil {
	    			trace := terr.New("cmd.stats.Count", err)
	    			ctx.Println(trace.Formatc())
	    		}
	    		num := strconv.Itoa(len(channels))
	    		expr := "channels"
	    		if len(channels) == 1 {
	    			expr = "channel"
	    		}
	        	ctx.Println("Found "+num+" "+expr)
	        } else {
	        	err := terr.Err("Unknown keyword: type help count to see the valid keywords")
	        	ctx.Println(err.Error())
	        }
        },
	}
	return command
}

// internal methods

func getStat(node *gocent.NodeInfo, name string) (string, *terr.Trace) {
	metrics := node.Metrics
	var msg string
	for k, v := range(metrics) {
		if k == name {
			msg = msg+" "+k+" : "+formatNum(v)
			return msg, nil
		}
	}
	msg = "Invalid metric "+name
	err := errors.New(msg)
	trace := terr.New("cmd.stat.getStat", err)
	return "", trace
}

func statsForNode(node *gocent.NodeInfo, mode string) string {
	msg := msgForNode(node)
	metrics := node.Metrics
	for k, v := range(metrics) {
		if mode == "all" {
			msg = msg+"\n - "+k+" : "+formatNum(v)
		} else if mode == "node" {
			if strings.HasPrefix(k, "node") {
				msg = msg+"\n - "+k+" : "+formatNum(v)
			} 
		} else if mode == "http" {
			if strings.HasPrefix(k, "http") {
				msg = msg+"\n - "+k+" : "+formatNum(v)
			}
		} else if mode == "client" {
			if strings.HasPrefix(k, "client") {
				msg = msg+"\n - "+k+" : "+formatNum(v)
			}
		}
	}
	return msg
}

func formatNum(num int64) string {
	str := strconv.FormatInt(num, 10)
	n := skittles.BoldWhite(str)
	msg := fmt.Sprintf("%s", n)
	return msg
}

func msgForNode(node *gocent.NodeInfo) string {
	addr := state.Cli.Host+":"+strconv.Itoa(state.Cli.Port)
	msg := "-------------------------------------------\n"
	msg = msg+"Stats for node "+node.Name+" ("+addr+")"
	msg = msg+"\n-------------------------------------------"
	return msg
}
