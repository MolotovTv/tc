package tc

import (
	"encoding/json"
	"fmt"
	io "io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// RunBranch ...
func RunBranch(config Config, buildType, branch string) (int, error) {
	payload := fmt.Sprintf(`
		<build branchName="%s">
			<buildType id="%s"/>
			<comment><text>api build</text></comment>
		</build>
	`, branch, buildType)

	r := strings.NewReader(payload)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/app/rest/buildQueue", config.URL),
		r,
	)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/xml")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return 0, fmt.Errorf("error making request: %s", string(body))
	}

	var build Build
	return build.ID, errors.WithStack(json.NewDecoder(res.Body).Decode(&build))
}
