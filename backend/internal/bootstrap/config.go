package bootstrap

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/spf13/viper"
)

func InitializeConfig(config string) *viper.Viper {
	// 初始化viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s", err))
	}
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&global.App.Config); err != nil {
		fmt.Println(err)
	}

	sanitizeConfig()

	return v
}

// sanitizeConfig cleans up loaded configuration values for safe use.
// Example: trim trailing slashes from base URLs to avoid double slashes when concatenating paths.
func sanitizeConfig() {
	global.App.Config.Notify.NoticeURL = strings.TrimRight(global.App.Config.Notify.NoticeURL, "/")
}
