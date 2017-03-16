package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
	"github.com/synw/centcli/libcentcli/datatypes"
)


func GetServers() (map[string]*datatypes.Server, *terr.Trace) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.centcli")
	err := viper.ReadInConfig()
	servers := make(map[string]*datatypes.Server)
	if err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			trace := terr.New("conf.GetServers", err)
			return servers, trace
		default:
			err := errors.New("Unable to locate config file")
			trace := terr.New("conf.GetServers", err)
			return servers, trace
		}
	}
	available_servers := viper.Get("servers").([]interface{})	
	for i, _ := range available_servers {
		sv := available_servers[i].(map[string]interface{})
		name := sv["name"].(string)
		host := sv["centrifugo_host"].(string)
		port := int(sv["centrifugo_port"].(float64))
		key := sv["centrifugo_key"].(string)		
		servers[name] = &datatypes.Server{name, host, port, key}
	}
	return servers, nil
}
