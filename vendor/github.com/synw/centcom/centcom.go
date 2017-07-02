package centcom

import (
	"fmt"
	"github.com/synw/centcom/state"
	"strconv"
)

func New(addr string, key string) *Cli {
	return NewClient(addr, key)
}

func Disconnect(cli *Cli) {
	cli.Conn.Close()
	if state.Verbosity > 0 {
		msg := "Disconnected from " + cli.Addr
		fmt.Println(msg)
	}
	close(cli.Channels)
}

func SetVerbosity(v int) {
	state.Verbosity = v
}

func State() string {
	v := strconv.Itoa(state.Verbosity)
	msg := "- Verbosity is set to " + v
	return msg
}

func PrintState() {
	fmt.Println(State())
}
