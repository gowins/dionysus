package config

import (
	"fmt"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	log.Setup()
	expectRes := []string{"4222", "g43g", "f2ff2", "vwevw"}
	getRes := []string{}
	configHandlers := []*WatchConfigHandler{
		{
			Key: "Redis.UserName",
			Func: func(valueString string) error {
				getRes = append(getRes, valueString)
				fmt.Printf("UserName valueString is %v\n", valueString)
				return nil
			},
		},
		{
			Key: "Mysql.DbName",
			Func: func(valueString string) error {
				getRes = append(getRes, valueString)
				fmt.Printf("DbName valueString is %v\n", valueString)
				panic("test panic")
				return nil
			},
		},
	}
	Setup(configHandlers...)
	viper.Set("Mysql.DbName", "4222")
	viper.WriteConfig()
	time.Sleep(time.Second * 2)
	viper.Set("Redis.UserName", "g43g")
	viper.WriteConfig()
	time.Sleep(time.Second * 2)
	viper.Set("Mysql.DbName", "f2ff2")
	viper.WriteConfig()
	time.Sleep(time.Second * 2)
	viper.Set("Redis.UserName", "vwevw")
	viper.WriteConfig()
	time.Sleep(time.Second * 2)
	for index, res := range expectRes {
		if res != getRes[index] {
			t.Errorf("want string %v, get string %v ", res, getRes[index])
			return
		}
	}
}
