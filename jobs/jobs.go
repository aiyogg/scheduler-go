package jobs

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config_real")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config files: %s", err))
	}
}

// Startup func
func Startup() {
	wxreadJob()
	clockInJob()
}
