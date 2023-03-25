package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	""
	"os"
)

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
	Model            string `toml:"model"`
	Personalization  string `toml:"personalization"`
	MaxTokens        int    `toml:"max_tokens"`
	Temperature      int    `toml:"temperature"`
	TopP             int    `toml:"top_p"`
	FrequencyPenalty int    `toml:"frequency_penalty"`
	PresencePenalty  int    `toml:"presence_penalty"`
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

// LoadConfig 加载配置
func LoadConfig() *Config {
	path, _ := os.Getwd()
	configPath := path + "../config.toml"

	config := &Config{}
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		fmt.Printf("Error decoding, error: %v\n", err)
	}
	return config
}

func Init() error {
	config := &Config{}
	path, _ := os.Getwd()
	configPath := path + "../config.toml"
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.
		return err
	}
}


