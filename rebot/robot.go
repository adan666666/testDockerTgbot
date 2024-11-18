package rebot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func TgRobot(config Conf) {
	var bot, err = tgbotapi.NewBotAPI(config.TgBot.Token)
	if err != nil {
		panic(err)
	}
	//bot.Debug = true

	msg := tgbotapi.NewMessage(config.TgBot.ChatID, "大佬们好，我是下班倒计时机器人")
	//发送消息
	//_, err = bot.Send(msg)
	//if err != nil {
	//	panic(err)
	//}
	u := tgbotapi.NewUpdate(0) //创建了一个新的更新对象 u，用于从 Telegram 服务器获取消息更新。参数 0 表示从最早的未读消息开始获取更新.
	u.Timeout = 60             //60秒内没有消息更新，就停止轮询，以节约资源

	//启动一个定时器 计算到下班还有多长时间
	go StartTimer(bot, config)
	// 获取一个监听管道，进行轮询监听飞机消息
	for update := range bot.GetUpdatesChan(u) {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		}
		// 打印收到的消息
		log.Infof("收到消息==>[%s] %s %v ", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

		// 检查消息是否提到了机器人 或者是命令
		if update.Message.IsCommand() || strings.Contains(update.Message.Text, "@"+bot.Self.UserName) {
			//if strings.HasPrefix(update.Message.Text, "@bx_xia_Bot") {
			// 回复消息
			responseText := "你提到我了吗？我在这里！大佬请指教！"
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			msg := tgbotapi.NewMessage(config.TgBot.ChatID, responseText)
			switch update.Message.Command() {
			case "start":
				msg.Text = "请输入1到100之间的数字 \n例如输入 /1"
			case "help":
				msg.Text = "You can control me by sending these commands:\n/start - to start the bot\n/help - to get this help message"
			default:
				if IsNumber(update.Message.Command()) {
					number, _ := strconv.Atoi(update.Message.Command())
					if number >= 1 && number <= 100 {
						duration := GetDuration(config.TgBot.Hour, config.TgBot.Min, config.TgBot.Sec)
						hour := duration / time.Hour
						mine := duration / time.Minute % 60
						second := duration / time.Second % 60
						msg.Text = fmt.Sprintf("下班倒计时: 还剩%d小时%d分钟%d秒", hour, mine, second)
					} else {
						msg.Text = "数字太大，我还在学习"
					}
				} else {
					if strings.Contains(update.Message.Text, "@bx_xia_Bot ") { //@我的(机器人)
						msg.Text = "不要@我，我很忙..."
					} else {
						msg.Text = "请重新输入..."
					}
				}
			}
			//msg.ReplyToMessageID = update.Message.MessageID  加这个是回复消息
			// 发送回复消息
			bot.Send(msg)
		} else {
			//如果是#号开头，就是我要发到群里的消息
			if strings.HasPrefix(update.Message.Text, "#") {
				msg.Text = update.Message.Text
				bot.Send(msg)
			}
		}
	}
}
