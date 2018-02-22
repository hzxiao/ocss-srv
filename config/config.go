package config

import (
	"fmt"
	"github.com/hzxiao/goutil/util"
	"github.com/spf13/viper"
)

func InitConfig(configName, configPath string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func PrintAll() {
	allConfigMap := viper.AllSettings()
	fmt.Println("--------config----------")
	for k, v := range allConfigMap {
		fmt.Printf("[%v]\n", k)
		one := util.MapV(v)
		for kk, vv := range one {
			fmt.Printf("%v = %v \n", kk, vv)
		}
		fmt.Println()
	}
	fmt.Println("------------------------")
}
