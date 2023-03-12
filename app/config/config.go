package config

import (
	"errors"
	"os"

	dotenv "github.com/golobby/dotenv"
	"github.com/golobby/env/v2"
)

func Init() (Conf, error) {
	conf := Conf{}
	var err error
	//Load Environment file into config
	if _, err := os.Stat(".env"); err != nil {
		// If .env not exists on local, it will use env var from OS
		err := env.Feed(&conf)
		if err != nil {
			return conf, errors.New("Error loading .env file from OS")
		}
	} else {
		// Use env from local development
		file, err := os.Open(".env")
		if err != nil {
			return conf, errors.New("Error loading .env file from file")
		} else {
			err = dotenv.NewDecoder(file).Decode(&conf)
			if err != nil {
				return conf, errors.New("Cannot decode .env file")
			}
		}
	}
	return conf, err
}
