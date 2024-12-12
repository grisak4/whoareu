package servconf

import "github.com/spf13/viper"

type servConf struct {
	Host string
	Port int
}

func GetServConf() servConf {
	return servConf{
		Host: viper.GetString("server.host"),
		Port: viper.GetInt("server.port"),
	}
}
