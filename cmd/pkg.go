package cmd

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/internal/git"
	"github.com/molotovtv/tc/tc"
	"github.com/spf13/cobra"
)

var (
	tagRegex = regexp.MustCompile("^v([0-9]+)\\.([0-9]+)\\.([0-9]+)$")
)

func init() {
	RootCmd.AddCommand(pkgCmd)
}

type version struct {
	Tag   string
	Major int
	Minor int
	Patch int
}

func versionFromTag(tag string) *version {
	matches := tagRegex.FindAllStringSubmatch(tag, -1)
	if len(matches) > 0 && len(matches[0]) == 4 {
		major, err := strconv.Atoi(matches[0][1])
		if err != nil {
			return nil
		}
		minor, err := strconv.Atoi(matches[0][2])
		if err != nil {
			return nil
		}
		patch, err := strconv.Atoi(matches[0][3])
		if err != nil {
			return nil
		}
		return &version{
			Tag:   tag,
			Major: major,
			Minor: minor,
			Patch: patch,
		}
	}
	return nil
}

func suggestFromTags(tags []string) string {
	// filter only valid tags
	validVersions := []*version{}
	for _, tag := range tags {
		if v := versionFromTag(tag); v != nil {
			validVersions = append(validVersions, v)
		}
	}

	if len(validVersions) == 0 {
		return "v1.0.0" // suggest v1.0.0 if no other idea
	}

	// sort by version
	sort.Slice(validVersions, func(i, j int) bool {
		if validVersions[i].Major < validVersions[j].Major {
			return true
		}
		if validVersions[i].Major > validVersions[j].Major {
			return false
		}
		if validVersions[i].Minor < validVersions[j].Minor {
			return true
		}
		if validVersions[i].Minor > validVersions[j].Minor {
			return false
		}
		if validVersions[i].Patch < validVersions[j].Patch {
			return true
		}
		if validVersions[i].Patch > validVersions[j].Patch {
			return false
		}
		return false
	})

	return fmt.Sprintf("v%d.%d.%d", validVersions[len(validVersions)-1].Major, validVersions[len(validVersions)-1].Minor, validVersions[len(validVersions)-1].Patch+1)
}

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Package a tag",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tag := ""
		if len(args) > 0 {
			tag = args[0]
		} else {
			tags, err := git.CurrentTags()
			if err != nil {
				log.Fatalf("%+v", err)
				return
			}
			if len(tags) == 0 {
				tags, err := git.Tags()
				if err != nil {
					log.Fatalf("%+v", err)
					return
				}
				suggestedTag := suggestFromTags(tags)
				fmt.Printf("No tag found in HEAD. Suggested tag : %s ; create ?", suggestedTag)
				var ok string
				fmt.Scanln(&ok)
				if err := git.CreateTagAndPush(suggestedTag); err != nil {
					log.Fatalf("%+v", err)
					return
				}
				tag = suggestedTag
			} else if len(tags) > 1 {
				log.Fatalf("more than one tag available, precise which one to package: %+v", tags)
				return
			} else {
				tag = tags[0]
			}
		}

		if !tagRegex.Match([]byte(tag)) {
			log.Fatalf("tag %s does not validate, tags for packaging must be vx.y.z", tag)
			return
		}
		c, err := config.Load()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		buildTypeID, err := config.BuildTypeID("packaging")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		fmt.Printf("Starting packaging version %s...\n", color.New(color.FgGreen).SprintFunc()(tag))
		buildID, err := tc.RunBranch(c, buildTypeID, tag)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		if runSilent {
			return
		}

		buildStatus(c, buildID)
	},
}
