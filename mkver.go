package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

// Config represents the set of arguments used for version calculation
type Config struct {
	profile           string
	env, gradle       string
	gitSha, gitRef    bool
	gitRefIgnore      []string
	gitBuildNum       string
	gitBuildNumBranch []string
}

var execCommand = exec.Command

// DefaultConfigs contain pre-configured short-cuts
var DefaultConfigs = map[string]Config{
	"gradle": Config{gitRef: true, gitRefIgnore: []string{"^develop$", "^master$", "^release", "^hotfix"}, gitBuildNum: "b", gitBuildNumBranch: []string{"^release", "^hotfix", "master"}},
	"npm":    Config{gitRef: true, gitRefIgnore: []string{"^develop$", "^master$", "^release", "^hotfix"}},
	"docker": Config{profile: "docker", gitSha: true, gitRef: true, gitRefIgnore: []string{"^develop$", "^master$", "^release", "^hotfix"}, gitBuildNum: "b"},
	"helm":   Config{gitSha: true, gitRef: true, gitRefIgnore: []string{"^develop$", "^master$", "^release", "^hotfix"}, gitBuildNum: "b"},
}

func main() {

	app := cli.NewApp()

	app.Name = "mkver"
	app.Usage = "Calculates application version by enriching the original one with various information"
	app.Version = "0.3.0"
	app.Commands = nil
	app.Flags = []cli.Flag{
		EnvFlag,
		GradleFlag,
		GitShaFlag,
		GitBuildNumFlag,
		GitBuildNumBranchFlag,
		GitRefFlag,
		GitRefIgnoreFlag,
		SnapshotFlag,
		ForFlag,
	}

	app.Action = func(ctx *cli.Context) {
		config := configure(*ctx)

		// Resolve the version, that will be used as a ground for further calculations
		// Version can be resolved from the env variable, gradle.properties or any other supported location
		version, err := resolveVersion(&config)
		if err != nil {
			log.Fatal(err)
		}

		// Resolve the git branch
		branch, err := resolveGitBranch(&config)
		if err != nil || len(branch) == 0 {
			branch = "unknown"
		}

		semanticVersion, err := Calculate(config, version, branch)

		if err != nil {
			log.Fatal(err)
		}

		if len(semanticVersion) == 0 {
			log.Fatal(errors.New("Failed to calculate version"))
		}

		fmt.Printf("%s", semanticVersion)
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// Calculate produces application version by enriching the original one with meta-informaiton based on the provided flags
func Calculate(config Config, version string, branch string) (string, error) {
	var versionBuilder strings.Builder

	// Splits original version by "-" into 2 parts: root and ext. F.e. 1.0.0-SNAPSHOT => 1.0.0 (root) and SNAPSHOT (ext)
	versionRoot, versionExt := resolveVersionRootAndExt(version)

	versionBuilder.WriteString(versionRoot)

	// Process git-ref. F.e. 1.0.0-SNAPSHOT on feature/x branch => 1.0.0-feature-x-SNAPSHOT
	processGitRef(&config, branch, &versionBuilder)

	// Process git-build-num. Will add build number taken from env variable to the result version.
	// F.e. 1.0.0 on the release/1.0.0 branch => 1.0.0-rcX (where x is a $BUILD_NUMBER env variable)
	processGitBuildNum(&config, branch, &versionBuilder)

	// Process git-sha. Will add git sha to the result version.
	// F.e. 1.0.0-SNAPSHOT => 1.0.0-ea3op1-SNAPSHOT
	processGitSha(&config, branch, &versionBuilder)

	// Appending back the version extension, which has been calculated together with a version root
	if len(versionExt) > 0 && config.profile == "gradle" {
		versionBuilder.WriteString("-" + versionExt)
	}

	return versionBuilder.String(), nil
}

//
// UTILS
//

func configure(ctx cli.Context) Config {
	var config Config

	// If "for" flag is present - apply one of the pre-defined configs
	if ctx.IsSet(ForFlag.Name) {
		name := ctx.String(ForFlag.Name)
		config = DefaultConfigs[name]
	}

	if ctx.IsSet(EnvFlag.Name) {
		config.env = ctx.String(EnvFlag.Name)
	}
	if ctx.IsSet(GradleFlag.Name) {
		config.gradle = ctx.String(GradleFlag.Name)
	}
	if ctx.IsSet(GitShaFlag.Name) {
		config.gitSha = ctx.Bool(GitShaFlag.Name)
	}
	if ctx.IsSet(GitBuildNumFlag.Name) {
		config.gitBuildNum = ctx.String(GitBuildNumFlag.Name)
	}
	if ctx.IsSet(GitBuildNumBranchFlag.Name) {
		config.gitBuildNumBranch = ctx.StringSlice(GitBuildNumBranchFlag.Name)
	}
	if ctx.IsSet(GitRefFlag.Name) {
		config.gitRef = ctx.Bool(GitRefFlag.Name)
	}
	if ctx.IsSet(GitRefIgnoreFlag.Name) {
		config.gitRefIgnore = ctx.StringSlice(GitRefIgnoreFlag.Name)
	}

	return config
}

// Resolves original version from one of the sources
func resolveVersion(cfg *Config) (string, error) {

	// Resolve from the provided env variable
	// In case, if "--env=.." flag is specified, resolve env variable with the provided name
	if len(cfg.env) > 0 {

		if val, found := os.LookupEnv(cfg.env); found {
			return val, nil
		}

		return "", fmt.Errorf("Failed to resolve version from env variable: $%s", cfg.env)
	}

	// Resolve from gradle
	if len(cfg.gradle) > 0 {
		if _, err := os.Stat(cfg.gradle); err == nil {
			gradleProperties, err := readPropertiesFile(cfg.gradle)
			if err != nil {
				return "", fmt.Errorf("Failed to resolve version from gradle properties file: $%s", cfg.gradle)
			}
			return gradleProperties["version"], nil
		}

		return "", fmt.Errorf("Failed to resolve version from gradle properties file: $%s", cfg.gradle)
	}

	// Try to auto-detect the original version source

	// Resolve from the default env variable (VERSION)
	if val, found := os.LookupEnv("VERSION"); found {
		return val, nil
	}

	// Resolve from the default gradle.properties file
	if _, err := os.Stat("gradle.properties"); err == nil {
		gradleProperties, err := readPropertiesFile("gradle.properties")
		if err != nil {
			return "", fmt.Errorf("Failed to resolve version from gradle.properties")
		}
		return gradleProperties["version"], nil
	}

	return "", errors.New("Failed to resolve version")
}

func resolveVersionRootAndExt(version string) (string, string) {
	if strings.Contains(version, "-") {
		versionParts := strings.Split(version, "-")
		return strings.TrimSpace(versionParts[0]), strings.TrimSpace(versionParts[1])
	}

	return version, ""
}

func resolveGitBranch(cfg *Config) (string, error) {
	// Determine the git branch from env if running on CI, otherwise from git
	if _, found := os.LookupEnv("BUILD_NUMBER"); found { // magic jenkins variable
		if _, found := os.LookupEnv("CHANGE_ID"); found { // Are we building a PR?
			val, _ := os.LookupEnv("CHANGE_BRANCH")
			return val, nil
		}

		val, _ := os.LookupEnv("BRANCH_NAME") // Not a PR
		return val, nil
	}

	// Not a CI build
	//TODO: change to lib?
	out, err := execCommand("bash", "-c", "git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown'").Output()
	return strings.TrimSpace(string(out[:])), err
}

func processGitRef(cfg *Config, branch string, versionBuilder *strings.Builder) {

	// Check if "--git-ref" flag is specified, otherwise - skip version processing
	if !cfg.gitRef {
		return
	}

	var ignore bool

	if len(cfg.gitRefIgnore) > 0 {
		for _, b := range cfg.gitRefIgnore {
			if match, _ := regexp.MatchString(b, branch); match {
				ignore = true
			}
		}
	}

	if !ignore {
		versionBuilder.WriteString("-" + strings.ToLower(strings.Replace(branch, "/", "-", -1)))
	}

}

func processGitBuildNum(cfg *Config, branch string, versionBuilder *strings.Builder) {
	if len(cfg.gitBuildNum) == 0 {
		return
	}

	var ignore = false

	if len(cfg.gitBuildNumBranch) > 0 {
		ignore = true
		for _, b := range cfg.gitBuildNumBranch {
			if match, _ := regexp.MatchString(b, branch); match {
				ignore = false
			}
		}
	}

	if !ignore {
		var buildNumber = "0"

		if val, found := os.LookupEnv("BUILD_NUMBER"); found {
			buildNumber = val
		}

		versionBuilder.WriteString("-" + cfg.gitBuildNum + buildNumber)
	}

}

func processGitSha(cfg *Config, branch string, versionBuilder *strings.Builder) {
	if !cfg.gitSha {
		return
	}

	out, _ := execCommand("bash", "-c", "git rev-parse --short=6 HEAD 2> /dev/null  || echo 'unknown'").Output()
	sha := strings.TrimSpace(string(out[:]))
	if "docker" == cfg.profile {
		versionBuilder.WriteString("+git." + sha)
	} else {
		versionBuilder.WriteString("-" + sha)
	}
}

func readPropertiesFile(filename string) (map[string]string, error) {
	config := map[string]string{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
