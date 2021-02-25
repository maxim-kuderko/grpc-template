package initializers

import (
	"github.com/spf13/viper"
	"os"
)

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(os.Getenv(`GO_ENV`))
	v.SetConfigType(`env`)
	v.AddConfigPath(`../configs`)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	v.AutomaticEnv()
	return v
}
