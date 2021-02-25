package initializers

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(fmt.Sprintf(`./configs/%s.env`, os.Getenv(`GO_ENV`)))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}
