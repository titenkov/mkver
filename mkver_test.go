package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

// Mock exec.command
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "1a2b3c")
	os.Exit(0)
}

var tests = []struct {
	config   Config
	version  string
	branch   string
	expected string
	err      error
}{
	{Config{}, "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},
	{Config{}, "master", "1.0.0", "1.0.0", nil},

	// --git-sha tests
	{Config{gitSha: true}, "develop", "1.0.0-SNAPSHOT", "1.0.0-1a2b3c-SNAPSHOT", nil},
	{Config{gitSha: true}, "master", "1.0.0", "1.0.0-1a2b3c", nil},

	// --git-ref tests
	{Config{gitRef: true}, "develop", "1.0.0-SNAPSHOT", "1.0.0-develop-SNAPSHOT", nil},
	{Config{gitRef: true}, "develop", "1.0.0", "1.0.0-develop", nil},
	{Config{gitRef: true}, "defect/X", "1.0.0-SNAPSHOT", "1.0.0-defect-x-SNAPSHOT", nil},
	{Config{gitRef: true}, "defect/X", "1.0.0", "1.0.0-defect-x", nil},
	{Config{gitRef: true}, "feature/TEST-123", "1.0.0", "1.0.0-feature-test-123", nil},

	// --git-ref-ignore tests
	{Config{gitRef: true, gitRefIgnore: []string{"^develop$"}}, "develop", "1.0.0", "1.0.0", nil},
	{Config{gitRef: true, gitRefIgnore: []string{"^release"}}, "release/1.0.0", "1.0.0", "1.0.0", nil},
	{Config{gitRef: true, gitRefIgnore: []string{"^release"}}, "feature/x", "1.0.0", "1.0.0-feature-x", nil},

	// --git-build-num tests
	{Config{gitBuildNum: "rc."}, "develop", "1.0.0-SNAPSHOT", "1.0.0-rc.13-SNAPSHOT", nil},
	{Config{gitBuildNum: "b"}, "develop", "1.0.0", "1.0.0-b13", nil},
	{Config{gitBuildNum: "rc."}, "release/1.0.0", "1.0.0", "1.0.0-rc.13", nil},

	// --git-build-num-branch tests
	{Config{gitBuildNum: "rc.", gitBuildNumBranch: []string{"^release", "^hotfix"}}, "release/1.0.0", "1.0.0-SNAPSHOT", "1.0.0-rc.13-SNAPSHOT", nil},
	{Config{gitBuildNum: "rc.", gitBuildNumBranch: []string{"^release", "^hotfix"}}, "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},

	// --git-auto-pilot=app tests
	{DefaultConfigs["app"], "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},
	{DefaultConfigs["app"], "develop-x", "1.0.0-SNAPSHOT", "1.0.0-develop-x-SNAPSHOT", nil},
	{DefaultConfigs["app"], "feature/x", "1.0.0-SNAPSHOT", "1.0.0-feature-x-SNAPSHOT", nil},
	{DefaultConfigs["app"], "defect/XYZ-123", "1.0.0-SNAPSHOT", "1.0.0-defect-xyz-123-SNAPSHOT", nil},
	{DefaultConfigs["app"], "release/1.0.0", "1.0.0", "1.0.0-rc.13", nil},
	{DefaultConfigs["app"], "hotfix/1.1.0", "1.1.0", "1.1.0-rc.13", nil},
	{DefaultConfigs["app"], "master", "1.0.0", "1.0.0", nil},

	// --git-auto-pilot=docker tests
	{DefaultConfigs["docker"], "develop", "1.0.0-SNAPSHOT", "1.0.0-1a2b3c-SNAPSHOT", nil},
	{DefaultConfigs["docker"], "develop-x", "1.0.0-SNAPSHOT", "1.0.0-develop-x-1a2b3c-SNAPSHOT", nil},
	{DefaultConfigs["docker"], "feature/x", "1.0.0-SNAPSHOT", "1.0.0-feature-x-1a2b3c-SNAPSHOT", nil},
	{DefaultConfigs["docker"], "defect/XYZ-123", "1.0.0-SNAPSHOT", "1.0.0-defect-xyz-123-1a2b3c-SNAPSHOT", nil},
	{DefaultConfigs["docker"], "release/1.0.0", "1.0.0", "1.0.0-rc.13-1a2b3c", nil},
	{DefaultConfigs["docker"], "hotfix/1.1.0", "1.1.0", "1.1.0-rc.13-1a2b3c", nil},
	{DefaultConfigs["docker"], "master", "1.0.0", "1.0.0-1a2b3c", nil},
}

func TestCalculateSemanticVersion(t *testing.T) {
	// prepare
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	os.Setenv("BUILD_NUMBER", "13")

	// execute
	for _, test := range tests {
		got, _ := Calculate(test.config, test.branch, test.version)
		assert.Equal(t, test.expected, got)
	}
}
