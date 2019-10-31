package tc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Config struct {
	URL      string
	UserName string
	Password string
	BuildIDs map[string]map[string]string
}

var localPath string

func findCurrentProject() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Base(wd), nil
}

func LoadConfig() (Config, error) {
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

func (c Config) Save() error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	path := usr.HomeDir + "/.tc"
	return ioutil.WriteFile(path, raw, 0644)
}

func (c Config) BuildTypeID(env string) (string, error) {
	currentProject, err := findCurrentProject()
	if err != nil {
		return "", err
	}
	if c.BuildIDs == nil {
		c.BuildIDs = make(map[string]map[string]string)
	}
	if envs, ok := c.BuildIDs[currentProject]; ok {
		if id, ok := envs[env]; ok {
			return id, nil
		}
	} else {
		c.BuildIDs[currentProject] = make(map[string]string)
	}

	fmt.Printf("Set buildID for project %s in env %s: ", currentProject, env)
	var id string
	fmt.Scanln(&id)
	c.BuildIDs[currentProject][env] = id
	if err := c.Save(); err != nil {
		return "", err
	}
	return id, nil
}
