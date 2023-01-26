package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/ChizarR/stats-service/pkg/logging"
)

type Config struct {
	IsDebug bool `json:"is_debug"`
	Server  struct {
		BindIP string `json:"bind_ip"`
		Port   string `json:"port"`
	} `json:"server"`
	MongoDB struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		Database    string `json:"database"`
		AuthDB      string `json:"auth_db"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Collections struct {
			Intaraction string `json:"intaraction"`
			User        string `json:"user"`
		} `json:"collections"`
	} `json:"mongo_db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("READING config file")
		instance = &Config{}
		if err := readConfig("config.json", instance); err != nil {
			logger.Fatal(err)
		}
	})
	return instance
}

func readConfig(path string, v any) error {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(config), v)
	if err != nil {
		return err
	}
	return nil
}
