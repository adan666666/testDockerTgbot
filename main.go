package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testDockerTgbot/rebot"
	"time"
)

// https://www.bilibili.com/video/BV1HH4y1W7H4/?spm_id_from=333.337.search-card.all.click&vd_source=55f7073cc1049edc8b91cea83217e7b6 视频
// https://www.fengfengzhidao.com/article/dtyibo4BEG4v2tWkcxXp 文档
func main() {
	r := gin.Default()
	file, err := os.ReadFile("settings.yaml")
	if err != nil {
		return
	}
	var conf rebot.Conf
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}
	//curl 127.0.0.1:5000
	r.GET("/", func(c *gin.Context) {
		c.JSON(200,
			gin.H{
				"code": 0,
				"msg":  "看到消息就说明布置成功了。",
				"data": gin.H{
					"tiem":   time.Now().Format("2006-01-02 15:04:05"),
					"token":  conf.TgBot.Token,
					"name":   conf.System.Name,
					"chatId": conf.TgBot.ChatID,
					"second": conf.TgBot.Sec,
				},
			})
	})
	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: &tr}
	fmt.Println(client)
	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	fmt.Println("启动机器人...")
	go rebot.TgRobot(conf)
	r.Run(":5000")
}
