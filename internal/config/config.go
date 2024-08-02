package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Addr     string `json:"addr"`
	Domain   string `json:"domain"`
	Upstream string `json:"upstream"`
	TLS      struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
	} `json:"tls"`
}

func MustLoad() Config {
	path := "config/config.json"
	slog.Info("try to read config", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		slog.Warn("read config failed", "error", err)
		if os.IsNotExist(err) {
			res, _ := os.Executable()
			path := filepath.Join(filepath.Dir(res), "config.json")
			slog.Info("try to read config", "path", path)
			data, err = os.ReadFile(path)
			if err != nil {
				slog.Warn("read config failed", "error", err)
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	slog.Info("success read config")

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
