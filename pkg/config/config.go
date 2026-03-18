package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		IP   string `mapstructure:"ip"`
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Pwd  string `mapstructure:"pwd"`
	}
	Mysql struct {
		DSN         string `mapstructure:"dsn"`
		MaxLifeTime int    `mapstructure:"maxLifeTime"`
		MaxOpenConn int    `mapstructure:"maxOpenConn"`
		MaxIdleConn int    `mapstructure:"maxIdleConn"`
	}
	Log struct {
		Level   string `mapstructure:"level"`
		LogPath string `mapstructure:"logPath"`
	} `mapstructure:"log"`
	Cos struct {
		SecretId  string `mapstructure:"secretId"`
		SecretKey string `mapstructure:"secretKey"`
		CDNDomain string `mapstructure:"cdnDomain"`
		BucketUrl string `mapstructure:"bucketUrl"`
	} `mapstructure:"cos"`
	DependOn struct {
		ShortUrl struct {
			Address     string `mapstructure:"address"`
			AccessToken string `mapstructure:"accessToken"`
		} `mapstructure:"shortUrl"`
		User struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"user"`
	} `mapstructure:"dependOn"`
}

var conf *Config

func InitConfig(filePath string, typ ...string) {
	v := viper.New()
	v.SetConfigFile(filePath)
	if len(typ) > 0 {
		v.SetConfigType(typ[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf = &Config{}
	err = v.Unmarshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	// 配置热更新
	v.OnConfigChange(func(in fsnotify.Event) {
		v.Unmarshal(conf)
	})
	v.WatchConfig()
}

func GetConfig() *Config {
	return conf
}
