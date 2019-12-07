package tc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunBranch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			t.FailNow()
		}
		if username != "test" || password != "test" {
			t.FailNow()
		}
		if req.URL.String() != "/app/rest/buildQueue" {
			t.FailNow()
		}
		if req.Method != http.MethodPost {
			t.FailNow()
		}
		encoder := json.NewEncoder(rw)
		err := encoder.Encode(Build{
			ID:         42,
			Status:     BuildStatusSuccess,
			BranchName: "master",
			State:      BuildStatusRunning,
		})
		if err != nil {
			t.Error(err)
		}
	}))
	// Close the server when test finishes
	defer srv.Close()

	conf := Config{
		URL:      srv.URL,
		UserName: "test",
		Password: "test",
	}
	buildID, err := RunBranch(conf, "Test_Project", "master")
	if err != nil {
		t.Error(err)
	}
	if buildID != 42 {
		t.Error("buildID must be 42")
	}
}
