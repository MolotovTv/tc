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
	buildIDs := map[string]string{}
	json.Unmarshal(raw, &buildIDs) // ignore malformed json
	return buildIDs, nil
}

func writeLocal(buildIDs map[string]string) error {
	raw, err := json.Marshal(buildIDs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(localPath, raw, 0644)
}

func BuildID(env string) (string, error) {
	buildIDs, err := readLocal()
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if id, ok := buildIDs[env]; ok {
		return id, nil
	}
	fmt.Printf("Build id for %s: ", env)
	var id string
	fmt.Scanln(&id)
	buildIDs[env] = id
	return id, writeLocal(buildIDs)
}

func SetBuildID(env, id string) (err error) {
	buildIDs, err := readLocal()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	buildIDs[env] = id
	return writeLocal(buildIDs)
}
