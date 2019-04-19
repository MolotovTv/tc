package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type Config struct {
	URL      string
	UserName string
	Password string
}

func Load() (Config, error) {
	usr, err := user.Current()
	if err != nil {
		return Config{}, err
	}
	path := usr.HomeDir + "/.tc"
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, err
		}
		config := Config{}
		fmt.Print("url: ")
		fmt.Scanln(&config.URL)

		fmt.Print("username: ")
		fmt.Scanln(&config.UserName)

		fmt.Print("password: ")
		fmt.Scanln(&config.Password)

		raw, err := json.Marshal(config)
		if err != nil {
			return Config{}, err
		}
		return config, ioutil.WriteFile(path, raw, 0644)
	}
	config := Config{}
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, f.Close()
}
