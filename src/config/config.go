package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config *ConfigType

type ConfigType struct {
	Env       string `mapstructure:"env"`
	Port      int    `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	Secret    string `mapstructure:"secret"`
	Host      string `mapstructure:"host"`
	Statics   string `mapstructure:"statics"`
	Templates string `mapstructure:"templates"`
	Database  struct {
		URL        string `mapstructure:"url"`
		SqlDir     string `mapstructure:"sqldir"`
		Migrations string `mapstructure:"migrations"`
	} `mapstructure:"database"`
	Sendgrid struct {
		Disabled  bool              `mapstructure:"disabled"`
		URL       string            `mapstructure:"url"`
		ApiKey    string            `mapstructure:"apikey"`
		Templates map[string]string `mapstructure:"templates"`
	} `mapstructure:"sendgrid"`
	Wallet struct {
		Agent       string `mapstructure:"agent"`
		AgentApiKey string `mapstructure:"agent_api_key"`
		Connect     string `mapstructure:"connect"`
	} `mapstructure:"wallet"`
	Upload struct {
		Bucket      string `mapstructure:"bucket"`
		CDN         string `mapstructure:"cdn"`
		Credentials string `mapstructure:"credentials"`
	} `mapstructure:"upload"`
	Cors struct {
		Origins []string `mapstructure:"origins"`
	} `mapstructure:"cors"`
	Nats struct {
		Url   string `mapstructure:"url"`
		Token string `mapstructure:"token"`
	} `mapstructure:"nats"`
	Platforms struct {
		Accounts string `mapstructure:"accounts"`
	} `mapstructure:"platforms"`
}

func Init(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %s", err)
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}

	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}
