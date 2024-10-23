package gm

import (
	"log"

	"github.com/spf13/viper"
)

type Cfg struct {
	FileName string
	FileType string
	FilePath []string
}

func (cfg Cfg) init() {
	viper.SetConfigName(cfg.FileName)
	viper.SetConfigType(cfg.FileType)
	for _, path := range cfg.FilePath {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("not found config file")
		} else {
			log.Fatalf("config file: %s \n", err)
		}
	}
}
