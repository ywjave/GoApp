package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init(filename string) (err error) {
	viper.SetConfigFile(filename) // 直接指定配置文件
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml") //用于远程获取配置文件时才生效？
	//viper.AddConfigPath(filepath)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(" viper.ReadInConfig() failed", err)
		return
		//panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
