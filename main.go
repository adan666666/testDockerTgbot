package main

import "github.com/gin-gonic/gin"

// https://www.bilibili.com/video/BV1HH4y1W7H4/?spm_id_from=333.337.search-card.all.click&vd_source=55f7073cc1049edc8b91cea83217e7b6 视频
// https://www.fengfengzhidao.com/article/dtyibo4BEG4v2tWkcxXp 文档
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"200": "看到消息就说明布置成功了"})
	})
	r.Run(":5000")
}
