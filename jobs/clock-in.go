package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var h Holiday

func init() {
	h = Holiday{}
	h.Init()
}

type Holiday struct {
	Workday []string `json:"workday"`
	Holiday []string `json:"holiday"`
}

func (h *Holiday) Init() {
	confFilePath := fmt.Sprintf("/tmp/%s_holiday.json", time.Now().Format("2006"))
	conf, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(conf, h)
	log.Printf("节假日信息：%+v", h)
}

func (h *Holiday) IsHoliday(data time.Time) bool {
	weekStr := data.Weekday().String()
	nowStr := data.Format("20060102")
	for _, v := range h.Holiday {
		if v == nowStr {
			return true
		}
	}
	if weekStr == "Saturday" || weekStr == "Sunday" {
		for _, v := range h.Workday {
			if v == nowStr {
				return false
			}
		}
		return true
	}
	return false
}

func workWxNotify(which string) func() {
	return func() {
		now := time.Now()
		if h.IsHoliday(now) {
			return
		}

		var text string
		if which == "上班" {
			text = "元气满满的一天开始啦，记得打卡哦！"
		} else if which == "下班" {
			text = "又到了愉快的下班时间啦，打卡下班！"
		}
		userids := viper.GetStringSlice("clockin.userids")
		touser := strings.Join(userids, "|")
		qyapi.SendText(touser, "", "", fmt.Sprintf("%s打卡提醒\n%s", which, text))
		log.Printf("%s打卡提醒\n%s", which, text)
	}
}

func clockInJob() {
	c := cron.New()

	c.AddFunc("0 55 8 * * *", workWxNotify("上班"))
	c.AddFunc("0 0 18 * * *", workWxNotify("下班"))

	c.Start()
}
