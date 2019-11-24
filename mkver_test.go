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
	name     string
	config   Config
	version  string
	branch   string
	expected string
	err      error
}{
	{"default", Config{}, "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},
	{"default", Config{}, "master", "1.0.0", "1.0.0", nil},

	// --git-sha
	{"--git-sha", Config{gitSha: true}, "develop", "1.0.0-SNAPSHOT", "1.0.0-1a2b3c-SNAPSHOT", nil},
	{"--git-sha", Config{gitSha: true}, "master", "1.0.0", "1.0.0-1a2b3c", nil},

	// --git-ref tests
	{"--git-ref", Config{gitRef: true}, "develop", "1.0.0-SNAPSHOT", "1.0.0-develop-SNAPSHOT", nil},
	{"--git-ref", Config{gitRef: true}, "develop", "1.0.0", "1.0.0-develop", nil},
	{"--git-ref", Config{gitRef: true}, "defect/X", "1.0.0-SNAPSHOT", "1.0.0-defect-x-SNAPSHOT", nil},
	{"--git-ref", Config{gitRef: true}, "defect/X", "1.0.0", "1.0.0-defect-x", nil},
	{"--git-ref", Config{gitRef: true}, "feature/TEST-123", "1.0.0", "1.0.0-feature-test-123", nil},

	// --git-ref-ignore tests
	{"--git-ref-ignore", Config{gitRef: true, gitRefIgnore: []string{"^develop$"}}, "develop", "1.0.0", "1.0.0", nil},
	{"--git-ref-ignore", Config{gitRef: true, gitRefIgnore: []string{"^release"}}, "release/1.0.0", "1.0.0", "1.0.0", nil},
	{"--git-ref-ignore", Config{gitRef: true, gitRefIgnore: []string{"^release"}}, "feature/x", "1.0.0", "1.0.0-feature-x", nil},

	// --git-build-num tests
	{"--git-build-num", Config{gitBuildNum: "rc."}, "develop", "1.0.0-SNAPSHOT", "1.0.0-rc.13-SNAPSHOT", nil},
	{"--git-build-num", Config{gitBuildNum: "b"}, "develop", "1.0.0", "1.0.0-b13", nil},
	{"--git-build-num", Config{gitBuildNum: "rc."}, "release/1.0.0", "1.0.0", "1.0.0-rc.13", nil},

	// --git-build-num-branch tests
	{"--git-build-num-branch", Config{gitBuildNum: "rc.", gitBuildNumBranch: []string{"^release", "^hotfix"}}, "release/1.0.0", "1.0.0-SNAPSHOT", "1.0.0-rc.13-SNAPSHOT", nil},
	{"--git-build-num-branch", Config{gitBuildNum: "rc.", gitBuildNumBranch: []string{"^release", "^hotfix"}}, "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},

	//
	// Profiles
	//

	// --for=gradle
	// It's important to support SNAPSHOT versioning for Java artifacts
	// {"--for=gradle", DefaultConfigs["gradle"], "develop", "1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "develop-x", "1.0.0-SNAPSHOT", "1.0.0-develop-x-SNAPSHOT", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "feature/x", "1.0.0-SNAPSHOT", "1.0.0-feature-x-SNAPSHOT", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "defect/XYZ-123", "1.0.0-SNAPSHOT", "1.0.0-defect-xyz-123-SNAPSHOT", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "defect/XYZ-123", "1.0.0", "1.0.0-defect-xyz-123-SNAPSHOT", nil}, // FAIL: Should probably be a snapshot?
	// {"--for=gradle", DefaultConfigs["gradle"], "release/1.0.0", "1.0.0", "1.0.0-b13", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "hotfix/1.1.0", "1.1.0", "1.1.0-b13", nil},
	// {"--for=gradle", DefaultConfigs["gradle"], "master", "1.0.0", "1.0.0-b13", nil},

	// --for=npm
	// {"--for=npm", DefaultConfigs["npm"], "develop", "1.0.0", "1.0.0", nil},
	// {"--for=npm", DefaultConfigs["npm"], "develop-x", "1.0.0", "1.0.0-develop-x", nil},
	// {"--for=npm", DefaultConfigs["npm"], "feature/x", "1.0.0", "1.0.0-feature-x", nil},
	// {"--for=npm", DefaultConfigs["npm"], "defect/XYZ-123", "1.0.0", "1.0.0-defect-xyz-123", nil},
	// {"--for=npm", DefaultConfigs["npm"], "release/1.0.0", "1.0.0", "1.0.0", nil},
	// {"--for=npm", DefaultConfigs["npm"], "hotfix/1.1.0", "1.1.0", "1.1.0", nil},
	// {"--for=npm", DefaultConfigs["npm"], "master", "1.0.0", "1.0.0", nil},

	// --for=docker tests
	{"--for=docker", DefaultConfigs["docker"], "develop", "1.0.0", "1.0.0-b13+git.1a2b3c", nil},
	{"--for=docker", DefaultConfigs["docker"], "develop", "1.0.0-SNAPSHOT", "1.0.0-b13+git.1a2b3c", nil}, // Should ignore snapshot suffix
	{"--for=docker", DefaultConfigs["docker"], "develop-x", "1.0.0", "1.0.0-develop-x-b13+git.1a2b3c", nil},
	{"--for=docker", DefaultConfigs["docker"], "feature/x", "1.0.0-SNAPSHOT", "1.0.0-feature-x-b13+git.1a2b3c", nil},
	{"--for=docker", DefaultConfigs["docker"], "defect/XYZ-123", "1.0.0-SNAPSHOT", "1.0.0-defect-xyz-123-1a2b3c-SNAPSHOT", nil},
	{"--for=docker", DefaultConfigs["docker"], "release/1.0.0", "1.0.0", "1.0.0-rc.13-1a2b3c", nil},
	{"--for=docker", DefaultConfigs["docker"], "hotfix/1.1.0", "1.1.0", "1.1.0-rc.13-1a2b3c", nil},
	{"--for=docker", DefaultConfigs["docker"], "master", "1.0.0", "1.0.0-1a2b3c", nil},

	// --for=helm tests
	// {"--for=helm", DefaultConfigs["docker"], "develop", "1.0.0-SNAPSHOT", "1.0.0-1a2b3c-SNAPSHOT", nil},
	// {"--for=helm", DefaultConfigs["docker"], "develop-x", "1.0.0-SNAPSHOT", "1.0.0-develop-x-1a2b3c-SNAPSHOT", nil},
	// {"--for=helm", DefaultConfigs["docker"], "feature/x", "1.0.0-SNAPSHOT", "1.0.0-feature-x-1a2b3c-SNAPSHOT", nil},
	// {"--for=helm", DefaultConfigs["docker"], "defect/XYZ-123", "1.0.0-SNAPSHOT", "1.0.0-defect-xyz-123-1a2b3c-SNAPSHOT", nil},
	// {"--for=helm", DefaultConfigs["docker"], "release/1.0.0", "1.0.0", "1.0.0-rc.13-1a2b3c", nil},
	// {"--for=helm", DefaultConfigs["docker"], "hotfix/1.1.0", "1.1.0", "1.1.0-rc.13-1a2b3c", nil},
	// {"--for=helm", DefaultConfigs["docker"], "master", "1.0.0", "1.0.0-1a2b3c", nil},
}

func TestMkver(t *testing.T) {
	// prepare
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	os.Setenv("BUILD_NUMBER", "13")

	// execute
	for _, test := range tests {
		got, _ := Calculate(test.config, test.branch, test.version)
		assert.Equal(t, test.expected, got, "failed while testing "+test.name)
	}
}
