package bootstrap

import (
	"fmt"
	"goInsight/global"
	commonTasks "goInsight/internal/app/common/tasks"
	dasTasks "goInsight/internal/app/das/tasks"
	"time"

	"github.com/robfig/cron/v3"
)

func InitializeCron() {
	global.App.Cron = cron.New()

	go func() {
		_, err := global.App.Cron.AddFunc("*/1 * * * *", func() {
			fmt.Println("Run SyncDBMeta At:", time.Now())
			commonTasks.SyncDBMeta()
		})
		if err != nil {
			global.App.Log.Error(err)
		}

		_, err = global.App.Cron.AddFunc("*/5 * * * *", func() {
			fmt.Println("Run KillTiDBQuery At:", time.Now())
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
