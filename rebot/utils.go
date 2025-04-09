package rebot

import (
	"fmt"
	"strconv"
	"time"
)

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func GetDuration(Hour int, Min int, Sec int) time.Duration {
	// 获取当前时间
	now := time.Now()
	fmt.Println("现在时间", now)
	// 计算今天20点的时间
	target := time.Date(now.Year(), now.Month(), now.Day(), Hour, Min, Sec, 0, now.Location())
	fmt.Println("计算今天6点的时间", target)
	// 如果当前时间已经是下午6点之后，则计算明天6点
	if now.After(target) {
		target = target.AddDate(0, 0, 1)
	}

	// 计算倒计时
	duration := target.Sub(now) //time.Sub方法用于计算两个时间点之间的时间差
	fmt.Println("时间差是：", duration)
	return duration
}
