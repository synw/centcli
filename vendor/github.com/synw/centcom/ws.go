package centcom

import (
	"fmt"
	"errors"
	"time"	
	"strconv"
	"encoding/json"
	"github.com/centrifugal/gocent"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	color "github.com/acmacalister/skittles"
	"github.com/synw/centcom/state"
)


type Msg struct {
	UID	string
	Channel string
	Payload interface{}
}

type Cli struct {
	Host string
	Port int
	Key string
	Http *gocent.Client
	Conn centrifuge.Centrifuge
	SubEvents *centrifuge.SubEventHandler
	Subs map[string]centrifuge.Sub
	Channels chan *Msg
	HttpOk bool
	IsConnected bool
}

func (cli Cli) CheckHttp() error {
	url := "http://"+cli.Host+":"+strconv.Itoa(cli.Port)
	cli.Http = gocent.NewClient(url, cli.Key, 5*time.Second)
	// test the http connection
	err := checkHttpConnection(&cli)
	if err != nil {
		return err
	}
	cli.HttpOk = true;
	return nil
}

func (cli Cli) Subscribe(channel string) error {
	sub, err := cli.Conn.Subscribe(channel, cli.SubEvents)
	if err != nil {
		return  err
	}
	cli.Subs[channel] = sub
	if state.Verbosity > 1 {
		msg := "Suscribed to channel "+channel
		fmt.Println(ok(msg))
	}
	return nil
}

func (cli Cli) Unsubscribe(channel string) error {
	sub, err := getSubscription(&cli, channel)
	if err != nil {
		return err		
	}
	err = sub.Unsubscribe()
	if err != nil {
		return err
	}
	delete(cli.Subs, channel)
	if state.Verbosity > 1 {
		msg := "Unsuscribed to channel "+channel
		fmt.Println(ok(msg))
	}	
	return nil
}

func (cli Cli) Publish(channel string, payload interface{}) error {
	if len(cli.Subs) == 0 {
		msg := "No subscription found for channel "+channel
		err := newErr(msg)
		return err
	}
	sub, err := getSubscription(&cli, channel)
	if err != nil {
		return err
	}
	dataBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	sub.Publish(dataBytes)
	if err != nil {
		return err		
	}
	if state.Verbosity > 2 {
		msg := "Message sent into channel "+channel
		fmt.Println(msg)
	}
	return nil
}

// constructors

func NewClient(host string, port int, key string) *Cli {
	addr := "http://"+host+":"+strconv.Itoa(port)
	http := gocent.NewClient(addr, key, 5*time.Second)
	var ws centrifuge.Centrifuge
	var subevents *centrifuge.SubEventHandler
	subs := make(map[string]centrifuge.Sub)
	c := make(chan *Msg)
	cl := Cli{host, port, key, http, ws, subevents, subs, c, false, false}
	return &cl
}

func NewMsg(uid string, channel_name string, payload interface{}) *Msg {
	msg := &Msg{uid, channel_name, payload}
	return msg
}

// initialization

func Connect(cli *Cli) (error) {
	// Never show secret to client of your application. Keep it on your application backend only.
	secret := cli.Key
	// Application user ID.
	user := "cli"
	// Current timestamp as string.
	timestamp := centrifuge.Timestamp()
	// Empty info.
	info := ""
	// Generate client token so Centrifugo server can trust connection parameters received from client.
	token := auth.GenerateClientToken(secret, user, timestamp, info)
	
	creds := &centrifuge.Credentials{
		User:      user,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}

	onMessage := func(sub centrifuge.Sub, rawmsg centrifuge.Message) error {
		//fmt.Println(fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), rawmsg))
		channel := fmt.Sprintf("%s", sub.Channel())	
		msg, err := decodeCentrifugeMsg(channel, &rawmsg)
		if err != nil {	
			msg := "Error decoding Centrifuge message: "+err.Error()
			err := newErr(msg)
			fmt.Println(err.Error())
		}
		cli.Channels <- msg
		return nil
	}

	onJoin := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		if state.Verbosity > 2 {
			fmt.Println(fmt.Sprintf("User %s joined channel %s with client ID %s", msg.User, sub.Channel(), msg.Client))
		}
		return nil
	}

	onLeave := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		if state.Verbosity > 2 {		
			fmt.Println(fmt.Sprintf("User %s left channel %s with client ID %s", msg.User, sub.Channel(), msg.Client))
		}
		return nil
	}
	
	onPrivateSub := func(c centrifuge.Centrifuge, req *centrifuge.PrivateRequest) (*centrifuge.PrivateSign, error) {
		info := ""
		sign := auth.GenerateChannelSign(cli.Key, req.ClientID, req.Channel, info)
		privateSign := &centrifuge.PrivateSign{Sign: sign, Info: info}
		return privateSign, nil
	}
	
	events := &centrifuge.EventHandler{
		OnPrivateSub: onPrivateSub,
	}

	subevents := &centrifuge.SubEventHandler{
		OnMessage: onMessage,
		OnJoin:    onJoin,
		OnLeave:   onLeave,
	}
	
	wsURL := "ws://"+cli.Host+":"+strconv.Itoa(cli.Port)+"/connection/websocket"
	conn := centrifuge.NewCentrifuge(wsURL, creds, events, centrifuge.DefaultConfig)

	err := conn.Connect()
	if err != nil {
		msg := "Error connecting to "+wsURL+" : "+err.Error()
		err = newErr(msg)
		return err
	}
	
	cli.Conn = conn
	cli.SubEvents = subevents
	if state.Verbosity > 0 {
		msg := "Connected to "+wsURL
		fmt.Println(ok(msg))
	}
	return nil

}

func DecodeHttpMsg(raw *gocent.Message) (*Msg, error) {
	msg, err := decodeRawMessage(raw.Channel, raw.Data)
	if err != nil {
		return msg, err
	}
	msg.UID = raw.UID
	return msg, nil
}

// internal methods

func getSubscription(cli *Cli, channel_name string) (centrifuge.Sub, error) {
	for name, sub := range(cli.Subs) {
		if name == channel_name {
			return sub, nil
		}
	}
	msg := "Can not find channel "+channel_name+" in server subscriptions"
	err := newErr(msg)
	var s centrifuge.Sub
	return s, err
}

func decodeRawMessage(channel string, raw *json.RawMessage) (*Msg, error) {
	msg := &Msg{}
	msg.Channel = channel
	byte, err := json.Marshal(raw)
	if err != nil {
		return msg, newErr(err.Error())
	}
	err = json.Unmarshal(byte, &msg.Payload)
	if err != nil {
		return msg, newErr(err.Error())
	}
	return msg, nil
}

func decodeCentrifugeMsg(channel string, centmsg *centrifuge.Message) (*Msg, error) {
	msg := &Msg{}
	msg.Channel = channel
	msg.UID = centmsg.UID
	var err error
	msg.Payload, err = decodeRawMessage(channel, centmsg.Data)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func checkHttpConnection(cli *Cli) error  {
	_, err := cli.Http.Publish("$devnull", []byte(`{"test": "test"}`))
	if err != nil {
		return newErr(err.Error())
	}
	return nil
}

func ok(msg string) string {
	msg = "["+color.Green("ok")+"] "+msg
	return msg
}

func newErr(msg string) error {
	err := errors.New(msg)
	return err
}

