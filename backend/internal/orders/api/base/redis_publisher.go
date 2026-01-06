package base

import (
	"context"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"
)

type RedisPublisher struct{}

func NewRedisPublisher() *RedisPublisher {
	return &RedisPublisher{}
}

// PublishWithRenderType publishes message with specified render type (processlist, gh-ost, etc.)
func (p *RedisPublisher) Publish(channel string, executionID string, renderType string, data any) {
	if err := utils.Publish(
		context.Background(),
		channel,
		executionID,
		renderType,
		data,
	); err != nil {
		global.App.Log.Error(err)
	}
}
