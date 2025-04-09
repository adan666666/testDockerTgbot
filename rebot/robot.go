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

	// 存储用户的选择
	userSelections := make(map[int64][]string)

	u := tgbotapi.NewUpdate(0) //创建了一个新的更新对象 u，用于从 Telegram 服务器获取消息更新。参数 0 表示从最早的未读消息开始获取更新.
	u.Timeout = 60             //60秒内没有消息更新，就停止轮询，以节约资源
	//u.Offset = -1              // 跳过旧的更新

	//启动一个定时器 计算到下班还有多长时间
	go StartTimer(bot, config)
	// 获取一个监听管道，进行轮询监听飞机消息
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		//if update.Message == nil { // 忽略任何非消息更新 //会把update.CallbackQuery消息过滤掉
		//	continue
		//}
		// 打印收到的消息
		if update.Message != nil {
			log.Infof("收到消息==>[userName=%s/From.String=%s/ID=%d] [消息是=%s}] [Chat.ID=%v] ", update.Message.From.UserName, update.Message.From.String(), update.Message.From.ID, update.Message.Text, update.Message.Chat.ID) //如果没有设置userName,From.String()可以取到name
		}

		if update.CallbackQuery != nil { // 用户点击按钮
			log.Printf("CallbackQuery received: %+v", update.CallbackQuery)
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID
			data := callback.Data

			if data == "confirm" { // 用户点击确认
				selection := userSelections[chatID]
				if len(selection) == 0 {
					reply := "你未选择任何选项！"
					bot.Send(tgbotapi.NewMessage(chatID, reply))
				} else {
					reply := fmt.Sprintf("你选择了: %s", strings.Join(selection, ", "))
					bot.Send(tgbotapi.NewMessage(chatID, reply))
				}
				userSelections[chatID] = nil // 清空用户选择
			} else { // 记录选择
				userSelections[chatID] = toggleSelection(userSelections[chatID], data)
				//reply := fmt.Sprintf("当前选择: %s", strings.Join(userSelections[chatID], ", "))
				//bot.Send(tgbotapi.NewMessage(chatID, reply))
			}

			// 回答 CallbackQuery
			bot.Send(tgbotapi.NewCallback(callback.ID, "操作已记录"))
		}

		if update.Message != nil && update.Message.IsCommand() { // 用户发送命令
			msg := update.Message
			chatID := msg.Chat.ID

			switch msg.Command() {
			case "start":
				// 发送选项按钮
				buttons := []tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData("选项1", "option1"),
					tgbotapi.NewInlineKeyboardButtonData("选项2", "option2"),
					tgbotapi.NewInlineKeyboardButtonData("选项3", "option3"),
				}
				confirmButton := tgbotapi.NewInlineKeyboardButtonData("确认", "confirm")

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(buttons...),
					tgbotapi.NewInlineKeyboardRow(confirmButton),
				)

				reply := tgbotapi.NewMessage(chatID, "请选择选项：")
				reply.ReplyMarkup = keyboard

				_, err = bot.Send(reply)
				if err != nil {
					panic(err)
				}
			}
		}

		// 检查消息是否提到了机器人 或者是命令
		if update.Message != nil && (update.Message.IsCommand() || strings.Contains(update.Message.Text, "@"+bot.Self.UserName)) {
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
			if update.Message != nil && strings.HasPrefix(update.Message.Text, "#") {
				msg.Text = update.Message.Text
				bot.Send(msg)
			}
		}

	}
}

// toggleSelection 用于更新用户的选择
func toggleSelection(selection []string, item string) []string {
	for i, v := range selection {
		if v == item {
			// 如果已经选中，取消选择
			return append(selection[:i], selection[i+1:]...)
		}
	}
	// 如果未选中，添加到选择中
	return append(selection, item)
}
