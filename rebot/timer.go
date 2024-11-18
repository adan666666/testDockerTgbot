package rebot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"time"
)

func StartTimer(tgBot *tgbotapi.BotAPI, conf Conf) {
	//获取当前时间，返回一个表示当前时间的 time.Time 对象
	now := time.Now()
	// 构造一个指定时间，通过传入年、月、日、时、分、秒、纳秒和时区来创建一个 time.Time 对象。
	target := time.Date(now.Year(), now.Month(), now.Day(), conf.TgBot.Hour, conf.TgBot.Min, 0, 0, now.Location())
	// 如果当前时间已经是今天晚上6点之后，则要等到明天
	if now.After(target) {
		target = target.AddDate(0, 0, 1)
	}
	// 计算定时器延时(到下班还需要多久)
	duration := target.Sub(now)
	// 启动一个定时器下
	timer := time.NewTimer(duration)
	// 阻塞到定时器到时
	<-timer.C

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.Infof("下班倒计时:%v %d小时%d分钟%d秒", duration, duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
	// 定时器到时，提醒下班
	sendMsg(tgBot, conf)
}

func sendMsg(tgBot *tgbotapi.BotAPI, conf Conf) {
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Millisecond * time.Duration(3-i) * 50)
		tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", 0, 0, 3-i)))
	}
	logrus.Infof("logrus 下班时间到，全体起立，离开工位!")
	//tgBot.Send(tgbotapi.NewMessage(conf.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
	StartTimer(tgBot, conf)
}
