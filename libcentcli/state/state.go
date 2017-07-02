package state

import (
	"errors"
	"fmt"
	"github.com/synw/centcli/libcentcli/conf"
	"github.com/synw/centcli/libcentcli/datatypes"
	"github.com/synw/centcom"
	"github.com/synw/terr"
)

var Servers map[string]*datatypes.Server
var Server *datatypes.Server
var Cli *centcom.Cli
var Listening []string
var ListenChan chan struct{}
var User string = "cli42"

func InitState(user string) *terr.Trace {
	User = user
	servers, trace := conf.GetServers()
	if trace != nil {
		trace := terr.Pass("state.InitState", trace)
		return trace
	}
	Servers = servers
	msg := "Found servers "
	for name, _ := range Servers {
		msg = msg + name + " "
	}
	fmt.Println(msg)
	return nil
}

func InitServer() *terr.Trace {
	centcom.SetVerbosity(1)
	cli := centcom.NewClient(Server.Addr, Server.Key)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("state.InitServer", err)
		return trace
	}
	cli.IsConnected = true
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("state.InitServer", err)
		return trace
	}
	fmt.Println(terr.Ok("Http transport ready"))
	Cli = cli
	return nil
}

func SetServer(name string) *terr.Trace {
	server, trace := ServerExists(name)
	if trace != nil {
		return trace
	}
	Server = server
	return nil
}

func ServerExists(server_name string) (*datatypes.Server, *terr.Trace) {
	for name, server := range Servers {
		if server_name == name {
			return server, nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	err := errors.New(msg)
	trace := terr.New("cmd.state.serverExists", err)
	server := &datatypes.Server{}
	return server, trace
}
