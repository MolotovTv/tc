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
	BuildIDs map[string]string `json:"build_ids"`
}

func (c *LocalConfig) BuildID(env string) string {
	return c.BuildIDs[env]
}

func (c *LocalConfig) BuildIDPromp(env string) string {
	id := c.BuildID(env)
	if id != "" {
		return id
	}

	fmt.Printf("Build id for %s: ", env)
	fmt.Scanln(&id)
	if id != "" {
		c.BuildIDs[env] = id
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
				BuildIDs: make(map[string]string),
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

	if config.BuildIDs == nil {
		config.BuildIDs = make(map[string]string)
	}

	return config, nil
}
