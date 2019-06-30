package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// CalculateSemanticVersion - calculates the semantic version depending on the gradle version (taken from the gradle.properties) and git branch
//
// Branch         gradle.version    semantic version
// develop        1.0.0-SNAPSHOT    1.0.0-SNAPSHOT
// some-branch    1.0.0-SNAPSHOT    1.0.0-some-branch-SNAPSHOT
// feature/x      1.0.0-SNAPSHOT    1.0.0-feature-x-SNAPSHOT
// defect/x       1.0.0-SNAPSHOT    1.0.0-feature-x-SNAPSHOT
// release/1.x    1.0.0             1.0.0-rc{BUILD}
// hotfix/1.x     1.0.0             1.0.0-rc{BUILD}
// master         1.0.0             1.0.0
func CalculateSemanticVersion(branch, version string) string {
	var versionRoot, versionExt, semanticVersion string

	// convert brenches like feature/XYZ to feature-xyz
	branch = strings.ToLower(strings.Replace(branch, "/", "-", -1))

	if strings.Contains(version, "-") {
		versionRoot = strings.TrimSpace(strings.Split(version, "-")[0])
		versionExt = strings.TrimSpace(strings.Split(version, "-")[1])
	} else {
		versionRoot = version
		versionExt = ""
	}

	// Calculate semantic version
	if branch == "develop" || branch == "master" {
		semanticVersion = versionRoot
	} else if strings.HasPrefix(branch, "release") || strings.HasPrefix(branch, "hotfix") {
		buildNumber := getEnvVariable("BUILD_NUMBER", "0")
		semanticVersion = fmt.Sprintf("%v-rc%v", versionRoot, buildNumber)
	} else {
		semanticVersion = fmt.Sprintf("%v-%v", versionRoot, branch)
	}

	if len(versionExt) > 0 {
		semanticVersion = fmt.Sprintf("%v-%v", semanticVersion, versionExt)
	}

	return strings.TrimSpace(semanticVersion)
}

// ResolveGradleVersion takes version property value from the gradle.properties file
func ResolveGradleVersion() (string, error) {
	out, err := exec.Command("bash", "-c", "cat gradle.properties | grep version | cut -d'=' -f2").Output()
	return strings.TrimSpace(string(out[:])), err
}

// ResolveGitBranch resolves the git branch based on either the environment variable (Jenkins build) or local git
func ResolveGitBranch() (string, error) {
	// Determine the git branch from env if running on CI, otherwise from git
	if isEnvVariableDefined("BUILD_NUMBER") { // magic jenkins variable
		if isEnvVariableDefined("CHANGE_ID") { // Are we building a PR?
			return getEnvVariable("CHANGE_BRANCH", "unknown"), nil
		}
		// Not a PR
		return getEnvVariable("BRANCH_NAME", "unknown"), nil
	}
	// Not a CI build
	out, err := exec.Command("bash", "-c", "git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown'").Output()
	return strings.TrimSpace(string(out[:])), err
}

func isEnvVariableDefined(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

func getEnvVariable(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {

	var app = cli.NewApp()

	app.Name = "mkver"
	app.Usage = "Calculates semantic version based on the branch and version taken from one of the sources (environment variable, gradle version, package.json, etc.)"
	app.Version = "0.2.0"

	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v", "--version"},
			Usage:   "Calculate semantic version",
			Action: func(c *cli.Context) {
				branch, _ := ResolveGitBranch()
				version, _ := ResolveGradleVersion()

				semanticVersion := CalculateSemanticVersion(branch, version)
				fmt.Printf("%s", semanticVersion)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
