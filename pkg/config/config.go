/*
Package config implements a simple library for config.
*/
package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//Config interface
type Config interface {
	GetCacheConfig() *CacheConfig
	GetDBConfig() *DatabaseConfig
}

//HostConfig host
type HostConfig struct {
	Address string
	Port    string
}

//CacheConfig cache
type CacheConfig struct {
	Dialect string
	Host    string
	Port    string
	URL     string
}

//DatabaseConfig db
type DatabaseConfig struct {
	Dialect  string
	DBName   string
	UserName string
	Password string
	URL      string
}

type config struct {
	docker   string
	Host     HostConfig
	Cache    CacheConfig
	Database DatabaseConfig
}

//var configChange = make(chan int, 1)

//NewConfig instance
func NewConfig() (*config, error) {
	var err error
	var config = new(config)

	env := os.Getenv("CONTACTENV")
	if env == "dev" {
		config.Database.UserName = os.Getenv("MONGOUSERNAME")
		buf, err := ioutil.ReadFile(os.Getenv("MONGOPWD"))
		if err != nil {
			panic(errors.New("read the env var fail"))
		}
		config.Database.Password = strings.TrimSpace(string(buf))
		config.Database.URL = os.Getenv("MONGOURL")
		config.Cache.URL = os.Getenv("REDISURL")
	} else if env == "local" || env == "" {
		fileName := "local.config"
		config, err = local(fileName)
	} else if env == "remote" {
		fileName := "remote.config"
		config, err = remote(fileName)
	} else {
		panic(errors.New("env var is invalid"))
	}

	//WatchConfig(configChange)
	return config, err
}

func local(fileName string) (*config, error) {
	path := os.Getenv("CONFIGPATH")
	if path == "" {
		path = "configs"
	}
	v := viper.New()
	config := new(config)
	v.SetConfigType("json")
	v.SetConfigName(fileName)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(config)
	if err != nil {
		panic(err)
	}

	return config, err
}

func remote(fileName string) (config *config, err error) {
	path := os.Getenv("CONFIGPATH")
	if path == "" {
		path = "configs"
	}
	v := viper.New()
	v.SetConfigType("json")

	err = v.AddRemoteProvider("etcd", "http://127.0.0.1:4001", path+fileName+".json")
	if err != nil {
		panic(err)
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(config)
	if err != nil {
		panic(err)
	}

	return
}

func (c *config) GetCacheConfig() *CacheConfig {
	return &c.Cache
}
func (c *config) GetDBConfig() *DatabaseConfig {
	return &c.Database
}

/*func WatchConfig(change chan int) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("config changed: %s", e.Name)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Warnf("read config fail after changing config")
			return
		}
		change <- 1
	})
}*/
