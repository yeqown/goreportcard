package types

import (
	"os"
	"path/filepath"

	"github.com/gojp/goreportcard/internal/repository"
	vcshelper "github.com/gojp/goreportcard/internal/vcs-helper"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var (
	_cfg           *Config
	_defaultConfig = &Config{
		Port: 8000,
		DB:   repository.Badger,
		VCS:  vcshelper.BuiltinTool,
		VCSOptions: []*vcshelper.VCSOption{
			{
				Host:           "github.com",
				Prefix:         "git",
				PrivateKeyPath: genPrivateKeyPath(),
			},
		},
		RepoRoot: "goreportcard-repos/",
		Domain:   "http://localhost:8000",
		SkipDirs: []string{},
		URIFormatRules: []uriFormatRule{
			{
				Prefix:    "github.com",
				URIFormat: "https://%s/blob/%s/%s",
			},
		},
	}
)

func init() {
	home, _ := os.UserHomeDir()
	_defaultConfig.RepoRoot = filepath.Join(home, _defaultConfig.RepoRoot)
}

// GetConfig get global config
func GetConfig() *Config {
	if _cfg == nil {
		_cfg = _defaultConfig
	}
	return _cfg
}

// Init .
// load config from file into the built-in _cfg variable
// if path is empty, then write with default config
func Init(confPath string) error {
	_cfg = new(Config)

	if _, err := toml.DecodeFile(confPath, _cfg); err != nil {
		if os.IsNotExist(err) {
			_cfg = _defaultConfig
			if err = writeConfig(confPath, _cfg); err != nil {
				return errors.Wrap(err, "types.Init.DecodeFile.writeConfig")
			}
		}

		return errors.Wrap(err, "types.Init.DecodeFile")
	}

	return nil
}

func writeConfig(confPath string, cfg *Config) error {
	fd, err := os.OpenFile(confPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return toml.NewEncoder(fd).Encode(cfg)
}

type Config struct {
	// server options
	Port       int                    `toml:"port"`
	DB         repository.DBType      `toml:"db"`
	VCS        vcshelper.VCSType      `toml:"vcs"`
	VCSOptions []*vcshelper.VCSOption `toml:"vcs_options"`
	RepoRoot   string                 `toml:"repoRoot"`
	Domain     string                 `toml:"domain"`

	// lint options
	SkipDirs []string `toml:"skipDirs"`

	// report options
	URIFormatRules []uriFormatRule `toml:"uriFormatRules"`
}

type uriFormatRule struct {
	Prefix    string `toml:"prefix"`
	URIFormat string `toml:"uriFormat"`
}

// genPrivateKeyPath get default private key path
func genPrivateKeyPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh", "id_rsa")
}
