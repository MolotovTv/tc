package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var localPath string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	localPath = filepath.Join(wd, ".tc")
}

func readLocal() (map[string]string, error) {
	raw, err := ioutil.ReadFile(localPath)
	if err != nil {
		return map[string]string{}, err
	}
	buildTypeIDs := map[string]string{}
	json.Unmarshal(raw, &buildTypeIDs) // ignore malformed json
	return buildTypeIDs, nil
}

func writeLocal(buildTypeIDs map[string]string) error {
	raw, err := json.Marshal(buildTypeIDs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(localPath, raw, 0644)
}

func BuildTypeID(env string) (string, error) {
	buildTypeIDs, err := readLocal()
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if id, ok := buildTypeIDs[env]; ok {
		return id, nil
	}
	fmt.Printf("Build id for %s: ", env)
	var id string
	fmt.Scanln(&id)
	buildTypeIDs[env] = id
	return id, writeLocal(buildTypeIDs)
}

func SetBuildTypeID(env, id string) (err error) {
	buildTypeIDs, err := readLocal()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	buildTypeIDs[env] = id
	return writeLocal(buildTypeIDs)
}
