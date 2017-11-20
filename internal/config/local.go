package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	localConfigPath = ".tc"
)

type LocalConfig struct {
	buildIDs map[Env]string `json:"build_ids"`
}

func (c *LocalConfig) BuildID(env Env) string {
	return c.buildIDs[env]
}

func (c *LocalConfig) BuildIDPromp(env Env) string {
	id := c.BuildID(env)
	if id != "" {
		return id
	}

	fmt.Print("Build id: ")
	fmt.Scanln(&id)
	if id != "" {
		c.buildIDs[env] = id
		err := c.save()
		if err != nil {
			fmt.Println(err)
		}
	}

	return id
}

func (c *LocalConfig) save() error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(localConfigPath, raw, 0644)
}

func Local() (*LocalConfig, error) {
	file, err := os.Open(localConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &LocalConfig{
				buildIDs: make(map[Env]string),
			}, nil
		}
		return nil, err
	}
	defer file.Close()

	config := &LocalConfig{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	if config.buildIDs == nil {
		config.buildIDs = make(map[Env]string)
	}

	return config, nil
}
