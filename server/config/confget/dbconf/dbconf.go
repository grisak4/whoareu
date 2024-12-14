package dbconf

import "github.com/spf13/viper"

type dbConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func GetDBConf() dbConf {
	return dbConf{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
	}
}
