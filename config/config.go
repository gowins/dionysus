package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/viper"
)

var (
	DefaultFilePath   = "./etc/"
	DefaultConfigName = "config"
)

type WatchHandler func(string) error

type WatchConfigHandler struct {
	Key      string
	OldValue string
	Func     WatchHandler
}

var watchConfigHandlers = []*WatchConfigHandler{}

func Setup(configHandlers ...*WatchConfigHandler) {
	watchConfigHandlers = append(watchConfigHandlers, configHandlers...)
	viper.SetConfigName(DefaultConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(DefaultFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %w", err)
	}

	for _, watchConfigHandler := range watchConfigHandlers {
		watchConfigHandler.OldValue = viper.GetString(watchConfigHandler.Key)
	}

	viper.OnConfigChange(runHandler)
	viper.WatchConfig()
}

func runHandler(in fsnotify.Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from config handler. Err:%v ", r)
		}
	}()
	for _, watchConfigHandler := range watchConfigHandlers {
		newValue := viper.GetString(watchConfigHandler.Key)
		if watchConfigHandler.OldValue != newValue {
			watchConfigHandler.OldValue = newValue
			if err := watchConfigHandler.Func(newValue); err != nil {
				log.Errorf("handler key %v change error %v", watchConfigHandler.Key, err)
			} else {
				log.Infof("key %v change config handler success", watchConfigHandler.Key)
			}
		}
	}
}
