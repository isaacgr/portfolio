package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log/slog"
	"os"
)

const (
	SERVER_SECTION = "Server"
	BLOG_SECTION   = "Blog"
)

type ServerSection struct {
	Host        string `ini:"host"`
	Port        int    `ini:"port"`
	SMTPHost    string `ini:"smtp_host"`
	SMTPUser    string `ini:"smtp_user"`
	SMTPPass    string `ini:"smtp_pass"`
	Key         string `ini:"recaptcha_key"`
	ProjectID   string `ini:"recaptcha_project_id"`
	Action      string `ini:"recaptcha_action"`
	Credentials string `ini:"recaptcha_credentials_file"`
}

type BlogSection struct {
	BlogPath string `ini:"blog_path"`
}

type Config struct {
	ServerSection
	BlogSection
}

type ConfigProvider struct {
	Path   string
	Config *Config
	log    *slog.Logger
}

func NewConfigProvider(
	path string,
	log *slog.Logger,
) (*ConfigProvider, error) {
	config, err := parseIniConfiguration(path, log)
	if err != nil {
		return nil, err
	}

	provider := ConfigProvider{
		Path:   path,
		Config: config,
		log:    log,
	}

	return &provider, nil
}

func parseIniConfiguration(
	iniPath string,
	log *slog.Logger,
) (*Config, error) {
	if _, err := os.Stat(iniPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(
			"Unable to read config.ini. Was one created? %s",
			iniPath,
		)
	}

	log.Info("Reading config.ini.")
	parser, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, iniPath)
	if err != nil {
		return nil, err
	}

	s := &ServerSection{}
	b := &BlogSection{}

	log.Debug("Parsing 'Server' section")
	err = parser.Section(SERVER_SECTION).MapTo(s)

	if err != nil {
		return nil, err
	}

	log.Debug("Parsing 'Blog' section")
	err = parser.Section(BLOG_SECTION).MapTo(b)

	if err != nil {
		return nil, err
	}

	c := &Config{
		ServerSection: *s,
		BlogSection:   *b,
	}

	log.Info("Done reading config.ini.")

	return c, nil

}
