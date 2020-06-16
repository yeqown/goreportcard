package model

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var (
	_cfg *Config

	_defaultConfig = &Config{
		Domain:   "goreportcard.com",
		SkipDirs: []string{},
		RepoRoot: "_repos/src/",
		URIFormatRules: []uriFormatRule{
			{
				Prefix:    "github.com",
				URIFormat: "https://%s/blob/%s/%s",
			},
		},
	}
)

// Init load config from file
func Init(confPath string) error {
	_cfg = new(Config)

	if _, err := toml.DecodeFile(confPath, _cfg); err != nil {
		return errors.Wrap(err, "model.Init.DecodeFile")
	}

	return nil
}

// Config type contains params to start web application
type Config struct {
	Domain         string          `toml:"domain"`
	SkipDirs       []string        `toml:"skipDirs"`
	RepoRoot       string          `toml:"repoRoot"`
	URIFormatRules []uriFormatRule `toml:"uriFormatRules"`
}

type uriFormatRule struct {
	Prefix    string `toml:"prefix"`
	URIFormat string `toml:"uriFormat"`
}

// GetConfig get global config
func GetConfig() *Config {
	if _cfg == nil {
		_cfg = _defaultConfig
	}
	return _cfg
}
