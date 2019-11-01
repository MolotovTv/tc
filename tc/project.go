package tc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
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

// CancelBuild ...
func CancelBuild(config Config, buildID int) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/app/rest/builds/id:%d", config.URL, buildID),
		strings.NewReader("<buildCancelRequest comment='build cancelled by api' readdIntoQueue='false' />"),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("error making request: %s", string(body))
	}
	return nil
}

// LastBuild ...
func LastBuild(config Config, buildTypeID string) (Build, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/app/rest/builds?locator=buildType:(id:%s),branch:default:any,running:any,defaultFilter:false,count:1", config.URL, buildTypeID),
		nil,
	)
	if err != nil {
		return Build{}, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Build{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return Build{}, fmt.Errorf("error making request: %s", string(body))
	}

	builds := buildResponse{}
	if err := json.NewDecoder(res.Body).Decode(&builds); err != nil {
		return Build{}, errors.WithStack(err)
	}

	if len(builds.Builds) > 0 {
		return builds.Builds[0], nil
	}
	return Build{}, nil
}

func LastBuildSuccessByBranch(config Config, buildTypeID string, branch string) (Build, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/app/rest/builds?locator=buildType:(id:%s),branch:%s,running:false,status:success,defaultFilter:false", config.URL, buildTypeID, branch),
		nil,
	)
	if err != nil {
		return Build{}, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Build{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return Build{}, fmt.Errorf("error making request: %s", string(body))
	}

	builds := buildResponse{}
	if err := json.NewDecoder(res.Body).Decode(&builds); err != nil {
		return Build{}, errors.WithStack(err)
	}

	if len(builds.Builds) > 0 {
		return builds.Builds[0], nil
	}
	return Build{}, nil
}

func GetBuild(config Config, buildID int) (DetailedBuild, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/app/rest/builds/id:%d", config.URL, buildID),
		nil,
	)
	if err != nil {
		return DetailedBuild{}, errors.WithStack(err)
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return DetailedBuild{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return DetailedBuild{}, fmt.Errorf("error making request: %s", string(body))
	}

	build := DetailedBuild{}
	return build, errors.WithStack(json.NewDecoder(res.Body).Decode(&build))
}
