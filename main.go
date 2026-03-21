package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type Result struct {
	NewTag string
	Major  string
	Minor  string
	Part   string
}

func main() {
	defaultBump := "minor"
	vPrefix := true

	if len(os.Args) > 1 {
		defaultBump = os.Args[1]
	}
	if len(os.Args) > 2 {
		vPrefix = os.Args[2] == "true"
	}

	latestTag := os.Getenv("LATEST_TAG")
	commitMsg := os.Getenv("COMMIT_MSG")

	result := CalculateVersion(latestTag, commitMsg, defaultBump, vPrefix)

	fmt.Printf("new_tag=%s\n", result.NewTag)
	fmt.Printf("major=%s\n", result.Major)
	fmt.Printf("minor=%s\n", result.Minor)
	fmt.Printf("part=%s\n", result.Part)

	if ghOutput := os.Getenv("GITHUB_OUTPUT"); ghOutput != "" {
		f, err := os.OpenFile(ghOutput, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to GITHUB_OUTPUT: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Fprintf(f, "new_tag=%s\n", result.NewTag)
		fmt.Fprintf(f, "major=%s\n", result.Major)
		fmt.Fprintf(f, "minor=%s\n", result.Minor)
		fmt.Fprintf(f, "part=%s\n", result.Part)
	}

	if result.Part != "none" {
		fmt.Fprintf(os.Stderr, "Creating release %s (%s bump from %s)\n", result.NewTag, result.Part, latestTag)
	}
}

func CalculateVersion(latestTag, commitMsg, defaultBump string, vPrefix bool) Result {
	bump := detectBump(commitMsg, defaultBump)
	if bump == "none" {
		return Result{Part: "none"}
	}

	ver := parseVersion(latestTag)
	bumped := bumpVersion(ver, bump)

	prefix := ""
	if vPrefix {
		prefix = "v"
	}

	return Result{
		NewTag: fmt.Sprintf("%s%d.%d.%d", prefix, bumped.Major(), bumped.Minor(), bumped.Patch()),
		Major:  fmt.Sprintf("%s%d", prefix, bumped.Major()),
		Minor:  fmt.Sprintf("%s%d.%d", prefix, bumped.Major(), bumped.Minor()),
		Part:   bump,
	}
}

func parseVersion(tag string) semver.Version {
	tag = strings.TrimPrefix(tag, "v")
	if tag == "" {
		return *semver.MustParse("0.0.0")
	}
	v, err := semver.NewVersion(tag)
	if err != nil {
		return *semver.MustParse("0.0.0")
	}
	return *v
}

func detectBump(commitMsg, defaultBump string) string {
	lower := strings.ToLower(commitMsg)
	switch {
	case strings.Contains(lower, "#major"):
		return "major"
	case strings.Contains(lower, "#minor"):
		return "minor"
	case strings.Contains(lower, "#patch"):
		return "patch"
	case strings.Contains(lower, "#none"):
		return "none"
	default:
		return defaultBump
	}
}

func bumpVersion(v semver.Version, bump string) semver.Version {
	switch bump {
	case "major":
		return v.IncMajor()
	case "minor":
		return v.IncMinor()
	case "patch":
		return v.IncPatch()
	default:
		return v
	}
}
