package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testDocker/rebot"
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(200,
			gin.H{
				"code": 0,
				"msg":  "看到消息就说明布置成功了。",
				"data": gin.H{
					"token":  conf.TgBot.Token,
					"name":   conf.System.Name,
					"chatId": conf.TgBot.ChatID,
					"second": conf.TgBot.Sec,
				},
			})
	})
	fmt.Println("启动机器人...")
	go rebot.TgRobot(conf)
	r.Run(":5000")
}
