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

	"github.com/titenkov/mkver/flags"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()

	app.Name = "mkver"
	app.Usage = "Calculates semantic version by enriching the original one with various aspects"
	app.Version = "0.3.0"
	app.Commands = nil
	app.Flags = []cli.Flag{
		flags.EnvFlag,
		flags.GradleFlag,
		flags.GitShaFlag,
		flags.GitBuildNumFlag,
		flags.GitBuildNumBranchFlag,
		flags.GitRefFlag,
		flags.GitRefIgnoreFlag,
		flags.AutPilotFlag,
	}

	app.Action = func(ctx *cli.Context) {
		semanticVersion, err := Calculate(*ctx)

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

// Resolves original version from one of the sources
func resolveVersion(ctx cli.Context) (string, error) {

	// Resolve from the provided env variable
	// In case, if "--env=.." flag is specified, resolve env variable with the provided name
	if ctx.IsSet(flags.EnvFlag.Name) {
		envFlag := ctx.GlobalString(flags.EnvFlag.Name)

		if len(envFlag) > 0 {
			if val, found := os.LookupEnv(envFlag); found {
				return val, nil
			}

			return "", fmt.Errorf("Failed to resolve version from env variable: $%s", envFlag)
		}
	}

	// Resolve from the default env variable (VERSION)
	if val, found := os.LookupEnv("VERSION"); found {
		return val, nil
	}

	// Resolve from gradle
	if ctx.IsSet(flags.GradleFlag.Name) {
		gradleFlag := ctx.GlobalString(flags.GradleFlag.Name)
		if len(gradleFlag) > 0 {
			if _, err := os.Stat(gradleFlag); err == nil {
				gradleProperties, err := readPropertiesFile(gradleFlag)
				if err != nil {
					return "", fmt.Errorf("Failed to resolve version from gradle properties file: $%s", gradleFlag)
				}
				return gradleProperties["version"], nil
			}

			return "", fmt.Errorf("Failed to resolve version from gradle properties file: $%s", gradleFlag)
		}
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

func resolveGitBranch(ctx cli.Context) (string, error) {
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
	out, err := exec.Command("bash", "-c", "git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown'").Output()
	return strings.TrimSpace(string(out[:])), err
}

// Calculate produces the semantic version according to the provided flags
func Calculate(ctx cli.Context) (string, error) {
	var versionBuilder strings.Builder
	var versionRoot, versionExt string

	// Resolve the version, that will be used as a ground for further calculations
	version, err := resolveVersion(ctx)

	if err != nil {
		return "", err
	}

	if strings.Contains(version, "-") {
		versionRoot = strings.TrimSpace(strings.Split(version, "-")[0])
		versionExt = strings.TrimSpace(strings.Split(version, "-")[1])
	} else {
		versionRoot = version
		versionExt = ""
	}

	versionBuilder.WriteString(versionRoot)

	// Resolve the git branch
	branch, err := resolveGitBranch(ctx)

	if err != nil || len(branch) == 0 {
		branch = "unknown"
	}

	if ctx.IsSet(flags.GitRefFlag.Name) && ctx.Bool(flags.GitRefFlag.Name) {
		var ignore bool

		if ctx.IsSet(flags.GitRefIgnoreFlag.Name) {
			var ignoreBranches = ctx.StringSlice(flags.GitRefIgnoreFlag.Name)

			for _, b := range ignoreBranches {
				if match, _ := regexp.MatchString(b, branch); match {
					ignore = true
				}
			}
		}

		if !ignore {
			versionBuilder.WriteString("-" + strings.ToLower(strings.Replace(branch, "/", "-", -1)))
		}
	}

	if ctx.IsSet(flags.GitBuildNumFlag.Name) {
		var ignore = false

		if ctx.IsSet(flags.GitBuildNumBranchFlag.Name) {
			var branches = ctx.StringSlice(flags.GitBuildNumBranchFlag.Name)
			ignore = true

			for _, b := range branches {
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

			versionBuilder.WriteString("-" + ctx.String(flags.GitBuildNumFlag.Name) + buildNumber)
		}
	}

	if ctx.IsSet(flags.GitShaFlag.Name) && ctx.Bool(flags.GitShaFlag.Name) {
		out, _ := exec.Command("bash", "-c", "git rev-parse --short=6 HEAD 2> /dev/null  || echo 'unknown'").Output()
		sha := strings.TrimSpace(string(out[:]))
		versionBuilder.WriteString("-" + sha)
	}

	if len(versionExt) > 0 {
		versionBuilder.WriteString("-" + versionExt)
	}

	return versionBuilder.String(), nil
}

//
// UTILS
//

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
