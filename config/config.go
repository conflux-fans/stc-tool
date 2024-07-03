package config

import (
	"fmt"
	"log"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func initByFile[T any](configPath string) *T {
	viper.SetConfigFile(configPath)
	return loadViper[T]()
}

func loadViper[T any]() *T {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
	logger.Get().WithField("file", viper.ConfigFileUsed()).Debug("Viper use config file")

	var _config T
	if err := viper.Unmarshal(&_config, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnset = true
	}); err != nil {
		panic(err)
	}
	return &_config
}

type Config struct {
	BlockChain struct {
		URL              string `yaml:"url"`
		FlowContract     string `yaml:"flowContract"`
		TemplateContract string `yaml:"templateContract"`
	} `yaml:"blockChain"`
	StorageNodes []string `yaml:"storageNodes"`
	KvNode       string   `yaml:"kvNode"`
	ZkNode       string   `yaml:"zkNode"`
	PrivateKeys  []string `yaml:"privateKeys"`
	Log          string   `yaml:"log"`
}

var (
	_config    Config
	configPath string = "./config.yaml"
)

const (
	DEBUG = "debug"
	INFO  = "info"
)

func SetConfigFile(path string) {
	configPath = path
}

func Init() {
	_config = *initByFile[Config](configPath)
}

func Get() *Config {
	return &_config
}
