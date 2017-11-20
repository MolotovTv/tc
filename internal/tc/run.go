package tc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aestek/tc/internal/config"
)

func RunBranch(config *config.GlobalConfig, build, branch string) (string, error) {
	payload := fmt.Sprintf(`
		<build branchName="%s">
			<buildType id="%s"/>
			<comment><text>api build</text></comment>
		</build>
	`, branch, build)

	r := strings.NewReader(payload)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/app/rest/buildQueue", config.URL),
		r,
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.SetBasicAuth(config.UserName, config.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}
