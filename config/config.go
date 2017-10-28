package dbconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type dbConfig struct {
	Configs []json.RawMessage `json:"config"`
}

type Config struct {
	Env      string `json:"env"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Database string `json:"database"`
}

func Read(env string, path string) (Config, error) {
	var res Config

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return res, fmt.Errorf("config read error")
	}

	var c dbConfig
	json.Unmarshal(raw, &c)

	for _, current := range c.Configs {

		var envConfig Config
		json.Unmarshal(current, &envConfig)

		if envConfig.Env == env {
			res = envConfig
			return res, nil
		}
	}

	return res, fmt.Errorf("no config for %s found", env)
}
