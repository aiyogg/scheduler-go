package jobs

import (
	"fmt"
	"time"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron"
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

// request 统一请求函数
func request(commitURL, jsonStrWithoutVid string, resultChan chan string) func(vid string) {
	return func(vid string) {
		jsonStr := []byte(fmt.Sprintf(jsonStrWithoutVid, vid, vid))
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
		resultChan <- string(body)
	}
}

// infiniteJob 无限卡
func infiniteJob() {
	now := time.Now()
	fmt.Println("infiniteJob run at: ", now.Format("2006-01-02 15:04:05"))

	vids := viper.GetStringSlice("wxread.vids")
	wechatURL := viper.GetString("wxread.wechat.infinite")
	commitURL := viper.GetString("wxread.commit.infinite")
	cycleStr := now.Format("20060102")

	resultList := make(chan string)
	jsonStrWithoutVid := fmt.Sprintf(`{"url": "%s?collageId=%s_%s&shareVid=%s"}`, wechatURL, "%s", cycleStr, "%s")
	reqFunc := request(commitURL, jsonStrWithoutVid, resultList)

	for _, vid := range vids {
		go reqFunc(vid)
	}

	for result := range resultList {
		fmt.Printf("InfiniteJob result: %s \n", result)
	}
}

// jizanJob 联名卡
func jizanJob() {
	now := time.Now()
	fmt.Println("jizanJob run at: ", now.Format("2006-01-02 15:04:05"))

	vids := viper.GetStringSlice("wxread.vids")
	wechatURL := viper.GetString("wxread.wechat.jizan")
	commitURL := viper.GetString("wxread.commit.jizan")
	cycleStr := now.Format("20060102")

	resultList := make(chan string)
	jsonStrWithoutVid := fmt.Sprintf(`{"url": "%s?isAnimateNavBarBackground=1&senderVid=%s&vol=%s&designId=%s_2&from=timeline"}`, wechatURL, "%s", cycleStr, "%s")
	reqFunc := request(commitURL, jsonStrWithoutVid, resultList)

	for _, vid := range vids {
		go reqFunc(vid)
	}

	for result := range resultList {
		fmt.Printf("jizanJob result: %s \n", result)
	}
}

// flipJob 翻一翻
func flipJob() {
	now := time.Now()
	fmt.Println("flipJob run at: ", now.Format("2006-01-02 15:04:05"))

	vids := viper.GetStringSlice("wxread.vids")
	wechatURL := viper.GetString("wxread.wechat.flip")
	commitURL := viper.GetString("wxread.commit.flip")
	cycleStr := now.Format("20060102")

	resultList := make(chan string)

	jsonStrWithoutVid := fmt.Sprintf(`{"url": "%s?vol=%s&inviteVid=%s&u=%s"}`, wechatURL, cycleStr, "%s", "%s")
	reqFunc := request(commitURL, jsonStrWithoutVid, resultList)

	for _, vid := range vids {
		go reqFunc(vid)
	}

	for result := range resultList {
		fmt.Printf("flipJob result: %s \n", result)
	}
}

func wxreadJob() {
	c := cron.New()

	c.AddFunc("0 0 12 * * 6", infiniteJob)
	c.AddFunc("0 0 21 * * 4", jizanJob)
	c.AddFunc("0 0 12 * * 2", flipJob)
	c.Start()

	select {}
}
