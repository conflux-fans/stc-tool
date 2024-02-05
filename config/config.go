package config

import (
	"fmt"
	"log"

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
	fmt.Printf("viper use config file: %v\n", viper.ConfigFileUsed())

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
	PrivateKeys  []string `yaml:"privateKeys"`
}

var (
	_config Config
)

func init() {
	InitByFile("./config.yaml")
}

func InitByFile(file string) {
	_config = *initByFile[Config](file)
}

func Get() *Config {
	return &_config
}
