package rebot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"time"
)

func StartTimer(tgBot *tgbotapi.BotAPI, conf Conf) {
	duration := GetDuration(conf.TgBot.Hour, conf.TgBot.Min, conf.TgBot.Sec)
	// 使用一个无限循环进行倒计时
	for {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		logrus.Infof("下班倒计时:%v %d小时%d分钟%d秒", duration, duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
		time.Sleep(1 * time.Second)
		duration = duration - 1*time.Second
		///优化方案
		// 当倒计时结束时退出for循环
		if duration < 2*time.Second { //duration < time.Second  时间还可以这样对比  duration < 0
			go sendMsg(tgBot, conf)
			StartTimer(tgBot, conf)
			break
		}
	}
}

func sendMsg(tgBot *tgbotapi.BotAPI, conf Conf) {
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Millisecond * time.Duration(3-i) * 50)
		tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", 0, 0, 3-i)))
	}
	logrus.Infof("logrus 下班时间到，全体起立，离开工位!")
	tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
}
