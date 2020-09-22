package config

import (
	"encoding/json"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"log"
	"path/filepath"
	"strings"
)

type Config struct {
	ServerConfig ServerConfig `json:"server"`
	JwtConfig    JWTConfig    `json:"jwt"`
	DBConfig     DBConfig     `json:"db"`
}

type ServerConfig struct {
	Port             int `json:"port"`
	TimeoutSecs      int `json:"timeoutSecs"`
	ReadTimeoutSecs  int `json:"readTimeoutSecs"`
	WriteTimeoutSecs int `json:"writeTimeoutSecs"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	SessionTime int    `json:"sessionTime"`
}

type DBConfig struct {
	DataSourceName string `json:"dataSourceName"`
	Migrate        bool   `json:"migrate"`
	Pool           struct {
		MaxOpen     int `json:"maxOpen"`
		MaxIdle     int `json:"maxIdle"`
		MaxLifetime int `json:"maxLifetime"`
	} `json:"pool"`
}

func (c *DBConfig) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"dataSourceName": "[PROTECTED]", // TODO : masking
		"pool": map[string]interface{}{
			"maxOpen":     c.Pool.MaxOpen,
			"maxIdle":     c.Pool.MaxIdle,
			"maxLifetime": c.Pool.MaxLifetime,
		},
		"migrate": c.Migrate,
	}
	return json.Marshal(m)
}

func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	// load from default config
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		log.Printf("failed to load default config. err: %v", err)
		return nil, err
	}

	// load from env
	err = k.Load(env.Provider("ARTICLE_SERVER_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "ARTICLE_SERVER_")), "_", ".", -1)
	}), nil)
	if err != nil {
		log.Printf("failed to load config from env. err: %v", err)
	}

	// load from config file if exist
	if configPath != "" {
		path, err := filepath.Abs(configPath)
		if err != nil {
			log.Printf("failed to get absoulute config path. configPath:%s, err: %v", configPath, err)
			return nil, err
		}
		log.Printf("load config file from %s", path)
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("failed to load config from file. err: %v", err)
			return nil, err
		}
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf. err: %v", err)
		return nil, err
	}
	return &cfg, err
}
