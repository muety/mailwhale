package config

import (
	"flag"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/jinzhu/configor"
	"io/ioutil"
	"os"
	"strings"
)

type smtpConfig struct {
	Host     string `env:"MW_SMTP_HOST"`
	Port     uint   `env:"MW_SMTP_PORT"`
	Username string `env:"MW_SMTP_USER"`
	Password string `env:"MW_SMTP_PASS"`
	TLS      bool   `env:"MW_SMTP_TLS"`
}

type webConfig struct {
	ListenV4 string `yaml:"listen_v4" default:"127.0.0.1:3000" env:"MW_WEB_LISTEN_V4"`
}

type storeConfig struct {
	Path string `default:"data.gob.db" env:"MW_STORE_PATH"`
}

type Config struct {
	Env     string `default:"dev" env:"MW_ENV"`
	Version string
	Web     webConfig
	Smtp    smtpConfig
	Store   storeConfig
}

var cfg *Config

func Get() *Config {
	return cfg
}

func Set(config *Config) {
	cfg = config
}

func Load() *Config {
	config := &Config{}

	flag.Parse()

	if err := configor.New(&configor.Config{}).Load(config, "config.yml"); err != nil {
		logbuch.Fatal("failed to read config: %v", err)
	}

	config.Version = readVersion()

	if config.Web.ListenV4 == "" {
		logbuch.Fatal("config option 'listen4' must be specified")
	}

	Set(config)
	return Get()
}

func (c *smtpConfig) ConnStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) IsDev() bool {
	return c.Env == "dev" || c.Env == "development"
}

func readVersion() string {
	file, err := os.Open("version.txt")
	if err != nil {
		logbuch.Fatal("failed to read version: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		logbuch.Fatal(err.Error())
	}

	return strings.TrimSpace(string(bytes))
}
