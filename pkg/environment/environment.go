package environment

import (
	"github.com/spf13/viper"
	"log"
)

const (
	ConfigName = "env"
	ConfigType = "yaml"
)

func ViperInit() error {
	viper.AddConfigPath("./env")

	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
		return err
	}
	return nil
}

func ViperGetEnvVariable(key string) (string, error) {
	err := ViperInit()

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatal("Invalid type assertion")
		return "", err
	}

	return value, err
}
