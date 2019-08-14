package jobs

import (
	"fmt"
	"time"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config files: %s", err))
	}
}

func infiniteJob() {
	url := "https://weread.qnmlgb.tech/submit"
	jsonStr := []byte(`{"url": "https://weread.qq.com/wrpage/infinite/lottery/?vol=20190814&inviteVid=21532019"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Origin", "https://weread.qnmlgb.tech")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Respones: ", string(body))
}

func wxreadJob() {
	// c := cron.New()

	// c.AddFunc("0 12 * * 6", infiniteJob)

	now := time.Now()
	vids := viper.GetStringSlice("wxread.vids")
	wechatURL := viper.GetString("wxread.wechat.infinite")
	commitURL := viper.GetString("wxread.commit.infinite")
	cycleStr := now.Format("20060102")

	resultList := make(chan string)
	for _, vid := range vids {
		go (func(vid string) {
			jsonStr := []byte(fmt.Sprintf(`{"url": "%s?vol=%s&inviteVid=%s"}`, wechatURL, cycleStr, vid))
			fmt.Println("Body: ", string(jsonStr))
			req, _ := http.NewRequest("POST", commitURL, bytes.NewBuffer(jsonStr))
			req.Header.Set("Origin", "https://weread.qnmlgb.tech")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			resultList <- string(body)
		})(vid)
	}

	var output string
	output = <-resultList
	fmt.Printf("Result: %s", output)

	infiniteJob()
}
