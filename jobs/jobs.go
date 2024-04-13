package jobs

import (
	"fmt"

	work "github.com/aiyogg/workweixin-go"
	"github.com/spf13/viper"
)

var qyapi work.WorkWeixin

func init() {
	viper.SetConfigName("config_real")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config files: %s", err))
	}

	corpid := viper.GetString("workweixin.appinfo.corpid")
	corpsecret := viper.GetString("workweixin.appinfo.corpsecret")
	agentid := viper.GetInt("workweixin.appinfo.agentid")

	qyapi.Init(corpid, corpsecret, agentid)

}

// Startup func
func Startup() {
	weiBoSubscribeJob()
	// wxreadJob()
	// clockInJob()

	select {}
}
