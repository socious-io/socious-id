package config

import (
	"os"

	"gopkg.in/yaml.v2"
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

func Init(filename string) (*ConfigType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	conf := new(ConfigType)
	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}
	Config = conf
	return conf, err
}
