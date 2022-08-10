package jobs

import (
	"fmt"
	"log"
	"time"

	"bytes"
	"io/ioutil"
	"net/http"

	work "github.com/aiyogg/workweixin-go"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var qyapi work.WorkWeixin

func init() {
	corpid := viper.GetString("workweixin.appinfo.corpid")
	corpsecret := viper.GetString("workweixin.appinfo.corpsecret")
	agentid := viper.GetInt("workweixin.appinfo.agentid")

	qyapi.Init(corpid, corpsecret, agentid)
}

// request 统一请求函数
func request(commitURL, jsonStrWithoutVid string, resultChan chan string) func(vid string) {
	return func(vid string) {
		jsonStr := []byte(fmt.Sprintf(jsonStrWithoutVid, vid, vid))
		log.Println("Body: ", string(jsonStr))
		req, _ := http.NewRequest("POST", commitURL, bytes.NewBuffer(jsonStr))
		req.Header.Set("Origin", "https://weread.qnmlgb.tech")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("Request Error: %s", err)
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
		log.Printf("InfiniteJob result: %s \n", result)
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
	jsonStrWithoutVid := fmt.Sprintf(`{"url": "%s&senderVid=%s&vol=%s&designId=%s_2&extra=%s"}`, wechatURL, "%s", cycleStr, cycleStr, "%s")
	reqFunc := request(commitURL, jsonStrWithoutVid, resultList)

	for _, vid := range vids {
		go reqFunc(vid)
	}

	for result := range resultList {
		log.Printf("jizanJob result: %s \n", result)
	}
}

// server酱提醒开启组队
func infinitePush() {
	sendRes := qyapi.SendText("@all", "", "", "微信读书组队提醒\n新一轮组队链接将于 1 个小时后自动提交，<a href=\"https://weread.qq.com/wrpage/infinite/lottery\">点击开启组队</a>")
	log.Printf("qyapi SendText result: %s", sendRes)
}

func wxreadJob() {
	c := cron.New()

	c.AddFunc("0 0 12 * * 6", infinitePush)
	c.AddFunc("0 0 13 * * 6", infiniteJob)
	c.AddFunc("0 0 21 * * 4", jizanJob)

	c.Start()
}
