package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/robfig/cron"
	viper "github.com/spf13/viper"
)

type weiBoIndex struct {
	Ok   int `json:"ok"`
	Data struct {
		Cards []struct {
			CardType int `json:"card_type"`
			Mblog    struct {
				ID           string `json:"id"`
				Mid          string `json:"mid"`
				Text         string `json:"text"`
				CreatedAt    string `json:"created_at"`
				RegionName   string `json:"region_name"`
				ThumbnailPic string `json:"thumbnail_pic"`
			} `json:"mblog"`
		} `json:"cards"`
	} `json:"data"`
}

type channel struct {
	uid         string
	containerid string
	res         json.RawMessage
	nickname    string
	receiver    string
}

type uidStore struct {
	uid string
	id  string
}

func (lwi *uidStore) save() {
	file, err := os.Create(fmt.Sprintf("/tmp/latest_weibo_id_%s.txt", lwi.uid))
	if err != nil {
		log.Fatalf("Create Error: %s", err)
	}
	defer file.Close()
	file.WriteString(lwi.id)
}
func (lwi *uidStore) get() string {
	file, err := os.OpenFile(fmt.Sprintf("/tmp/latest_weibo_id_%s.txt", lwi.uid), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Open Error: %s", err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Read Error: %s", err)
	}
	return string(buf)
}

func getWeiBoIndex(uid, containerid, nickname, receiver string, ch chan channel) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?type=uid&value=%s&containerid=%s", uid, containerid), nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request Error: %s", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	ch <- channel{
		uid:         uid,
		containerid: containerid,
		res:         json.RawMessage(body),
		nickname:    nickname,
		receiver:    receiver,
	}
}

func checkLatestWeiBo() {
	users := viper.Get("WeiBoSubscribe.users")
	ch := make(chan channel)
	if users != nil {
		usersArray := users.([]interface{})
		for _, user := range usersArray {
			userMap := user.(map[interface{}]interface{})
			uid := userMap["uid"].(string)
			containerid := userMap["containerid"].(string)
			nickname := userMap["nickname"].(string)
			receiver := userMap["receiver"].(string)

			go getWeiBoIndex(uid, containerid, nickname, receiver, ch)
		}

		for result := range ch {
			var weiBoIndex weiBoIndex
			err := json.Unmarshal(result.res, &weiBoIndex)
			if err != nil {
				log.Fatalf("Unmarshal Error: %s", err)
			}
			if weiBoIndex.Ok != 1 {
				log.Fatalf("WeiBoSubscribe Error: %d", weiBoIndex.Ok)
			}

			latestCard := weiBoIndex.Data.Cards[0]
			uidStore := &uidStore{uid: result.uid}
			oldId := uidStore.get()
			if latestCard.Mblog.ID != oldId {
				uidStore.id = latestCard.Mblog.ID
				uidStore.save()
				createTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", latestCard.Mblog.CreatedAt)
				fmtTime := createTime.Format("2006.01.02 15:04:05")
				reBreakLine := regexp.MustCompile(`<br />`)
				plainText := reBreakLine.ReplaceAllString(latestCard.Mblog.Text, "\n")
				reEmoji := regexp.MustCompile(`<.*?>`)
				plainText = reEmoji.ReplaceAllString(plainText, "")
				place := latestCard.Mblog.RegionName
				link := "https://m.weibo.cn/detail/" + latestCard.Mblog.Mid
				sendRes := qyapi.SendText(result.receiver, "", "", fmt.Sprintf("%s %s \n %s  \n%s \n<a href='%s'>直达现场 >></a>", result.nickname, place, fmtTime, plainText, link))
				log.Printf("qyapi SendText result: %s", sendRes)
			} else {
				log.Printf("%s(UID: %s): No new WeiBo, latest one was created at: %s", result.nickname, result.uid, latestCard.Mblog.CreatedAt)
			}
		}
	}
}

func weiBoSubscribeJob() {
	c := cron.New()

	c.AddFunc("0 */2 * * * *", checkLatestWeiBo)

	c.Start()
}
