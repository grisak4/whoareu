package jwtsec

import "github.com/spf13/viper"

func GetJwtToken() []byte {
	return []byte(viper.GetString("jwt.secret"))
}
