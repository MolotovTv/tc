package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var configPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configPath = usr.HomeDir + "/.tc"
}

type Config struct {
	URL      string
	UserName string
	Password string
	BuildIDs map[string]map[string]string // project, env, build_id
}

func Load() (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config := &Config{}

			fmt.Print("url: ")
			fmt.Scanln(&config.URL)

			fmt.Print("username: ")
			fmt.Scanln(&config.UserName)

			fmt.Print("password: ")
			fmt.Scanln(&config.Password)

			err = config.save()
			if err != nil {
				return nil, err
			}

			return config, nil
		}
		return nil, err
	}

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) BuildID(project, env string) string {
	if c.BuildIDs == nil {
		return ""
	}
	p, ok := c.BuildIDs[project]
	if !ok {
		return ""
	}

	return p[env]
}

func (c *Config) SetBuildID(project, env, buildID string) error {
	if c.BuildIDs == nil {
		c.BuildIDs = make(map[string]map[string]string)
	}
	_, ok := c.BuildIDs[project]
	if !ok {
		c.BuildIDs[project] = make(map[string]string)
	}
	c.BuildIDs[project][env] = buildID
	return c.save()
}

func (c *Config) BuildIDPrompt(project, env string) string {
	id := c.BuildID(project, env)
	if id != "" {
		return id
	}

	fmt.Printf("Build id for %s: ", env)
	fmt.Scanln(&id)
	if id != "" {
		err := c.SetBuildID(project, env, id)
		if err != nil {
			fmt.Println(err)
		}
	}

	return id
}

func (c *Config) save() error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, raw, 0644)
}
