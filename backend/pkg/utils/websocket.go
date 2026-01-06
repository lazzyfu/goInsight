package utils

import (
	"context"
	"encoding/json"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/redis/go-redis/v9"
)

const (
	RenderLogStream         = "log"
	RenderProcessListStream = "processlist"
	RenderGhostStream       = "gh-ost"
)

type PublishMSG struct {
	ExecutionID string `json:"execution_id"` // 一次执行的唯一标识，如taskid
	Type        string `json:"type"`         // gh-ost / processlist / log
	Data        any    `json:"data"`
}

// 生产消息到redis
func Publish(ctx context.Context, channel string, executionID string, renderType string, data any) error {
	msg := PublishMSG{
		ExecutionID: executionID,
		Type:        renderType,
		Data:        data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		global.App.Log.Error("publish marshal error", err)
		return err
	}

	return global.App.Redis.Publish(ctx, channel, payload).Err()
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
