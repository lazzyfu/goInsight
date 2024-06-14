package utils

import (
	"context"
	"encoding/json"
	"goInsight/global"

	"github.com/redis/go-redis/v9"
)

type PublishMSG struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 生产消息到redis
func Publish(c context.Context, channel string, data interface{}, renderType string) error {
	var msg PublishMSG
	msg.Type = renderType
	msg.Data = data
	jsonData, err := json.Marshal(msg)
	if err != nil {
		global.App.Log.Error(err)
	}
	return global.App.Redis.Publish(c, channel, jsonData).Err()
}

// 订阅redis消息
func Subscribe(sub *redis.PubSub, ctx context.Context, channel string, ch chan string) {
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		global.App.Log.Error("订阅redis消息 Err: ", err)
		return
	}
	ch <- msg.Payload
}
