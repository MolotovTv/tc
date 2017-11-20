package git

import (
	"os/exec"
	"strings"
)

func Branch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
