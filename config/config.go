package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
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
		PmContract       string `yaml:"pmContract"`
		StartBlockNum    int64  `yaml:"startBlockNum"`
	} `yaml:"blockChain"`
	StorageNodes   []string `yaml:"storageNodes"`
	KvNode         string   `yaml:"kvNode"`
	KvStreamId     string   `yaml:"kvStreamId"`
	ZkNode         string   `yaml:"zkNode"`
	PrivateKeyFile string   `yaml:"privateKeyFile"`
	Log            string   `yaml:"log"`
	ExtendData     struct {
		TextMaxSize int64 `yaml:"textMaxSize"`
	} `yaml:"extendData"`
}

var (
	_config      Config
	_privateKeys []string
	configPath   string = "./config.yaml"
)

const (
	DEBUG = "debug"
	INFO  = "info"
)

func SetConfigFile(path string) {
	configPath = path
}

func Init() {
	viper.SetDefault("privateKeyFile", getDefaultPrivateKeyPath())
	_config = *initByFile[Config](configPath)
	setLogger()
	_privateKeys = loadPrivateKeys()
}

func setLogger() {
	logLevel, err := logrus.ParseLevel(Get().Log)
	if err != nil {
		panic(err)
	}
	fmt.Println("log level:", logLevel)
	logger.Get().SetLevel(logLevel)
}

func Get() *Config {
	return &_config
}

func GetPrivateKeys() []string {
	return _privateKeys
}

func loadPrivateKeys() []string {
	privateKeyFile := _config.PrivateKeyFile
	logger.Get().WithField("file", privateKeyFile).Debug("Load private keys from file")
	content, err := os.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(content), "\n")
}

func getDefaultPrivateKeyPath() string {
	switch _os := runtime.GOOS; _os {
	case "windows":
		userProfile, ok := os.LookupEnv("USERPROFILE")
		if !ok {
			panic("USERPROFILE environment variable not set")
		}
		return path.Join(userProfile, ".storage-cli", "private_keys")
	case "linux", "darwin":
		home, ok := os.LookupEnv("HOME")
		if !ok {
			panic("HOME environment variable not set")
		}
		return path.Join(home, ".storage-cli", "private_keys")
	default:
		panic("Unknown system os")
	}
}
