package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
		viper.AddConfigPath("../conf/")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf/")
	}
	// 设置配置文件格式为TOML
	viper.SetConfigType("toml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
