package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf用来保存程序的所有配置信息
var Conf = new(BlogConfig)

type BlogConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	//*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

//type RedisConfig struct {
//	Host         string `mapstructure:"host"`
//	Password     string `mapstructure:"password"`
//	Port         int    `mapstructure:"port"`
//	DB           int    `mapstructure:"db"`
//	PoolSize     int    `mapstructure:"pool_size"`
//	MinIdleConns int    `mapstructure:"min_idle_conns"`
//}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func Init(filename string) (err error) {
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("./conf/")

	//viper.SetConfigFile("./conf/config.yaml")
	viper.SetConfigFile(filename)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
		return err
	}

	//把读取到的配置信息反序化到Conf变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("Fatal error unmarshal config file: %s \n", err)
		return err
	}

	viper.WatchRemoteConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("Fatal error unmarshal config file: %s \n", err)
		}
	})

	return
}
