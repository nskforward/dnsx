package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Addr     string            `json:"addr"`
	Upstream string            `json:"upstream"`
	TTL      int               `json:"ttl"`
	Routes   map[string]string `json:"routes"`
}

func MustLoad() Config {
	path := "config/config.json"

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			path := "/app/config.json"
			data, err = os.ReadFile(path)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	for k := range cfg.Routes {
		if k[len(k)-1] != '.' {
			panic("route domain must contain a dot at the end")
		}
	}

	return cfg
}
