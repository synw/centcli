package centcom

import (
	"fmt"
	"strconv"
	"github.com/synw/centcom/state"
)


func New(host string, port int, key string) *Cli {
	return NewClient(host, port, key)
}

func Disconnect(cli *Cli) {
	cli.Conn.Close()
	if state.Verbosity > 0 {
		msg := "Disconnected from "+cli.Host
		fmt.Println(msg)
	}
	close(cli.Channels)	
}

func SetVerbosity(v int) {
	state.Verbosity = v
}

func State() string {
	v := strconv.Itoa(state.Verbosity)
	msg := "- Verbosity is set to "+v
	return msg
}

func PrintState() {
	fmt.Println(State())
}
