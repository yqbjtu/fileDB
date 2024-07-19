package config

import (
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/klog"
	"os"
	"sync"
)

// GlobalConfig /*
type GlobalConfig struct {
	Version   string    `yaml:"version"`
	Port      int       `yaml:"port"`
	OSMConfig OSMConfig `yaml:"OSMConfig"`

	Postgres struct {
		Host      string `mapstructure:"host"`
		DBName    string `mapstructure:"dbname"`
		Port      int    `mapstructure:"port"`
		Username  string `mapstructure:"user"`
		Password  string `mapstructure:"passwd"`
		BatchSize int    `mapstructure:"batchSize"`
	} `mapstructure:"postgres"`

	LogConfig LogConfig `yaml:"logConfig"`
	//etc ..
}

type OSMConfig struct {
	DataDir string `yaml:"dataDir" mapstructure:"dataDir"`
	TempDir string `yaml:"tempDir" mapstructure:"tempDir"`
}

type MySQLConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Database string `yaml:"database" mapstructure:"database"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	IsConn   string `yaml:"isConn" mapstructure:"isConn"`
}

type LogConfig struct {
	LogPath    string `yaml:"logPath" mapstructure:"logPath"`
	ShowCaller bool   `yaml:"showCaller" mapstructure:"showCaller"`
	Level      string `yaml:"level" mapstructure:"level"`
}

var gConfig *GlobalConfig
var m sync.Mutex
var once sync.Once

func GetConfig() *GlobalConfig {
	if gConfig != nil {
		return gConfig
	} else {
		fmt.Printf("no init global conf, so use default value")
		panic("no init conf ")
	}
}

func InitConfig(confPath string) {
	once.Do(func() {
		var err error

		v := viper.New()
		v.SetConfigType("yaml")
		v.SetConfigFile(confPath) // 指定配置文件路径
		err = v.ReadInConfig()    // 读取配置信息
		if err != nil {
			dir, _ := os.Getwd()
			panic(fmt.Errorf("Fatal error conf file, err:%s, currentDir:%s\n", err, dir))
		}

		// 监控配置文件变化
		v.WatchConfig()

		err = v.Unmarshal(&gConfig)
		if err != nil {
			klog.Errorf("[viper] Unmarshal err:%s", err.Error())
			panic(err)
		}

	})
}

func (cfg *GlobalConfig) GetVersion() string {
	return cfg.Version
}

func (cfg *GlobalConfig) GetPort() int {
	return cfg.Port
}
