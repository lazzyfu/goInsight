package bootstrap

import (
	"goInsight/global"
	commonTasks "goInsight/internal/apps/common/tasks"
	dasTasks "goInsight/internal/apps/das/tasks"
	"time"

	"github.com/robfig/cron/v3"
)

func InitializeCron() {
	global.App.Cron = cron.New()

	go func() {
		_, err := global.App.Cron.AddFunc(global.App.Config.Crontab.SyncDBMetas, func() {
			global.App.Log.Info("Run SyncDBMeta At:", time.Now())
			commonTasks.SyncDBMeta()
		})
		if err != nil {
			global.App.Log.Error(err)
		}

		_, err = global.App.Cron.AddFunc("*/5 * * * *", func() {
			global.App.Log.Info("Run KillTiDBQuery At:", time.Now())
			kill := dasTasks.KillTiDBQuery{}
			kill.Run()
		})
		if err != nil {
			global.App.Log.Error(err)
		}
		global.App.Cron.Start()
		defer global.App.Cron.Stop()
		select {}
	}()
}
