package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var globalConfigPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	globalConfigPath = usr.HomeDir + "/.tc"
}

type GlobalConfig struct {
	URL      string
	UserName string
	Password string
}

func Global() (*GlobalConfig, error) {
	file, err := os.Open(globalConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			config := &GlobalConfig{}

			fmt.Print("url: ")
			fmt.Scanln(&config.URL)

			fmt.Print("username: ")
			fmt.Scanln(&config.UserName)

			fmt.Print("password: ")
			fmt.Scanln(&config.Password)

			raw, err := json.Marshal(config)
			if err != nil {
				return nil, err
			}

			err = ioutil.WriteFile(globalConfigPath, raw, 0644)
			if err != nil {
				return nil, err
			}

			return config, nil
		}
		return nil, err
	}

	config := &GlobalConfig{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
