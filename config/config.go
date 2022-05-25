package config

import (
	"flag"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/jinzhu/configor"
	"github.com/muety/mailwhale/types"
	"io/ioutil"
	"os"
	"strings"
)

const (
	KeyUser   = "user"
	KeyClient = "client"
)

type EmailPasswordTuple struct {
	Email    string
	Password string
}

type mailConfig struct {
	Domain          string `yaml:"domain" env:"MW_MAIL_DOMAIN"`
	SystemSenderTpl string `yaml:"system_sender" env:"MW_MAIL_SYSTEM_SENDER" default:"MailWhale System <system@{0}>"`
}

type smtpConfig struct {
	Host          string `env:"MW_SMTP_HOST"`
	Port          uint   `env:"MW_SMTP_PORT"`
	Username      string `env:"MW_SMTP_USER"`
	Password      string `env:"MW_SMTP_PASS"`
	TLS           bool   `env:"MW_SMTP_TLS"`
	SkipVerifyTLS bool   `yaml:"skip_verify_tls" env:"MW_SMTP_SKIP_VERIFY_TLS"`
}

type webConfig struct {
	ListenV4    string   `yaml:"listen_v4" env:"MW_WEB_LISTEN_V4"` // deprecated, use ListenAddr
	ListenAddr  string   `yaml:"listen_addr" default:"127.0.0.1:3000" env:"MW_WEB_LISTEN_ADDR"`
	CorsOrigins []string `yaml:"cors_origins" env:"MW_WEB_CORS_ORIGINS"`
	PublicUrl   string   `yaml:"public_url" default:"http://localhost:3000" env:"MW_WEB_PUBLIC_URL"`
}

type storeConfig struct {
	Path string `default:"data.json.db" env:"MW_STORE_PATH"`
}

type securityConfig struct {
	Pepper        string `env:"MW_SECURITY_PEPPER"`
	AllowSignup   bool   `yaml:"allow_signup" env:"MW_SECURITY_ALLOW_SIGNUP" default:"true"`
	VerifySenders bool   `yaml:"verify_senders" default:"true" env:"MW_SECURITY_VERIFY_SENDERS"`
	VerifyUsers   bool   `yaml:"verify_users" default:"true" env:"MW_SECURITY_VERIFY_USERS"`
}

type Config struct {
	Env      string `default:"dev" env:"MW_ENV"`
	Version  string
	Mail     mailConfig
	Web      webConfig
	Smtp     smtpConfig
	Store    storeConfig
	Security securityConfig
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

	if config.Web.ListenV4 != "" {
		config.Web.ListenAddr = config.Web.ListenV4 // for backwards-compatbility
	}

	if config.Web.ListenAddr == "" {
		logbuch.Fatal("config option 'listen_addr' must be specified")
	}

	if !config.Mail.SystemSender().Valid() {
		logbuch.Fatal("system sender address is invalid")
	}

	logbuch.Info("---")
	logbuch.Info("This instance is assumed to be publicly accessible at: %v", config.Web.GetPublicUrl())
	logbuch.Info("User registration enabled: %v", config.Security.AllowSignup)
	logbuch.Info("Account activation required: %v", config.Security.VerifyUsers)
	logbuch.Info("Sender address verification required: %v", config.Security.VerifySenders)
	logbuch.Info("---")

	Set(config)
	return Get()
}

func (c *webConfig) GetPublicUrl() string {
	return strings.TrimSuffix(c.PublicUrl, "/")
}

func (c *smtpConfig) ConnStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *mailConfig) SystemSender() types.MailAddress {
	return types.MailAddress(strings.Replace(c.SystemSenderTpl, "{0}", c.Domain, -1))
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
