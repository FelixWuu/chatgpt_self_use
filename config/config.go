package config

import (
	"github.com/BurntSushi/toml"
	"github.com/FelixWuu/chatgpt_self_use/utils/logger"
	"os"
)

var config *Config

type Config struct {
	Api     `toml:"api"`
	Bot     `toml:"bot"`
	Auth    `toml:"auth"`
	Address `toml:"address"`
}

type Api struct {
	ApiKey string `toml:"api_key"`
	ApiUrl string `toml:"api_url"`
}

type Bot struct {
	Model            string  `toml:"model"`
	Personalization  string  `toml:"personalization"`
	MaxTokens        int     `toml:"max_tokens"`
	Temperature      int     `toml:"temperature"`
	TopP             float32 `toml:"top_p"`
	FrequencyPenalty float32 `toml:"frequency_penalty"`
	PresencePenalty  float32 `toml:"presence_penalty"`
}

type Auth struct {
	AuthUser     string `toml:"auth_user"`
	AuthPassword string `toml:"auth_password"`
}

type Address struct {
	Port   int    `toml:"port"`
	Listen string `toml:"listen"`
	Proxy  string `toml:"proxy"`
}

func Init() error {
	config = &Config{}
	path, _ := os.Getwd()
	configPath := path + "../config.toml"
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		logger.Errorf("Init config failed, error: %v\n", err)
		return err
	}
	return nil
}

func Inst() *Config {
	return config
}
