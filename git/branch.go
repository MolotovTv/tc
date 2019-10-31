package git

import (
	"os/exec"
	"strings"
)

func ShowRef(ref string) (string, error) {
	out, err := exec.Command("git", "show", "--pretty=medium", "--quiet", ref).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// CurrentBranch returns the current branch name
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CurrentCommit returns the current commit hash
func CurrentCommit() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CurrentTags returns tags pointing to current commit
func CurrentTags() ([]string, error) {
	err := exec.Command("git", "fetch", "--tags").Run()
	if err != nil {
		return nil, err
	}
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return nil, err
	}
	currentCommit := strings.TrimSpace(string(out))
	out, err = exec.Command("git", "tag", "--points-at", currentCommit).Output()
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for _, tag := range strings.Split(string(out), "\n") {
		tag = strings.TrimSpace(tag)
		if len(tag) > 0 {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

// Tags returns all tags
func Tags() ([]string, error) {
	err := exec.Command("git", "fetch", "--tags").Run()
	if err != nil {
		return nil, err
	}
	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for _, tag := range strings.Split(string(out), "\n") {
		tag = strings.TrimSpace(tag)
		if len(tag) > 0 {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func CreateTagAndPush(tag string) error {
	if err := exec.Command("git", "tag", tag).Run(); err != nil {
		return err
	}
	return exec.Command("git", "push", "--tags").Run()
}
