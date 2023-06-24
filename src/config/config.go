package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var configConf []byte

type Config struct {
	ChatGPT struct {
		Keys      string `yaml:"keys"`
		HttpProxy string `yaml:"http_proxy"`
	} `yaml:"chat_gpt"`
	App struct {
		Addr string `yaml:"addr"`
	}
}

func GetConfig() *Config {
	config := Config{}
	if len(configConf) == 0 {
		filePath := "./config/config.yaml"
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("解析配置文件 %s 失败: %s\n", filePath, err)
		}
		configConf = content
	}

	if yaml.Unmarshal(configConf, &config) != nil {
		log.Fatalf("yaml 解析失败: %s\n", "")
	}
	return &config
}
