package rebot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func StartTimer(tgBot *tgbotapi.BotAPI, conf Conf) {
	c := cron.New(cron.WithSeconds())
	// 每天 几 点 几 分 几 秒运行任务
	spec := fmt.Sprintf("%d %d %d * * *", conf.TgBot.Sec, conf.TgBot.Min, conf.TgBot.Hour)
	_, err := c.AddFunc(spec, func() {
		sendMsg(tgBot, conf)
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	defer c.Stop()
	// 阻塞主线程，防止程序提前退出
	select {}
}

func sendMsg(tgBot *tgbotapi.BotAPI, conf Conf) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// 定时器到时，提醒下班
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Millisecond * time.Duration(3-i) * 50)
		tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", 0, 0, 3-i)))
	}
	logrus.Infof("logrus 下班时间到，全体起立，离开工位!")
	tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
}
