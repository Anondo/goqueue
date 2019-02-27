package config

import (
	"goqueue/helper"
	"os"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	helper.FailOnError(viper.ReadInConfig(), "Failed to read config file")
	helper.FailOnError(makeFiles(), "Failed to created persistance files")
}

func makeFiles() error {
	fp := viper.GetString("persistance.filepath")

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		err = os.Mkdir(viper.GetString("persistance.foldername"), 0777)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(fp)

	file.WriteString("[]")
	file.Close()

	if err != nil {
		return err
	}

	return nil
}
