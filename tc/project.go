package tc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/molotovtv/tc/internal/config"
)

type buildResponse struct {
	Builds []Build `json:"build"`
}

// BuildStatus ...
type BuildStatus string

const (
	// BuildStatusSuccess ...
	BuildStatusSuccess BuildStatus = "SUCCESS"
	// BuildStatusFailure ...
	BuildStatusFailure BuildStatus = "FAILURE"
	// BuildStatusError ...
	BuildStatusError BuildStatus = "ERROR"
)

// BuildState ...
type BuildState string

const (
	// BuildStateQueued ...
	BuildStateQueued BuildState = "queued"
	// BuildStatusRunning ...
	BuildStatusRunning BuildState = "running"
	// BuildStatusFinished ...
	BuildStatusFinished BuildState = "finished"
)

// LastBuild ...
func LastBuild(config *config.Config, build string) (*Build, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/app/rest/builds?locator=buildType:(id:%s),branch:default:any,running:any", config.URL, build),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("error making request: %s", string(body))
	}

	builds := buildResponse{}
	err = json.NewDecoder(res.Body).Decode(&builds)
	if err != nil {
		return nil, err
	}

	if len(builds.Builds) > 0 {
		return &builds.Builds[0], nil
	}

	return nil, nil
}
