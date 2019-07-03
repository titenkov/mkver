package flags

import (
	"github.com/urfave/cli"
)

// EnvFlag allows resolving version from the env variable
var EnvFlag = cli.StringFlag{
	Name:  "env",
	Value: "VERSION",
	Usage: "Resolve version from env variable",
}

// GradleFlag allows resolving version from the gradle properties file
var GradleFlag = cli.StringFlag{
	Name:  "gradle",
	Value: "gradle.properties",
	Usage: "Resolve version from gradle",
}

// ReleaseFlag allows creating release version
// F.e. 1.0.0-SNAPSHOT -> 1.0.0
// var ReleaseFlag = cli.BoolFlag{
// 	Name:  "release",
// 	Usage: "Produce release version",
// }

// VerifyReleaseFlag allows verify the version to be a release one
// F.e. 1.0.0 -> True, 1.0.0-SNAPSHOT -> False
// var VerifyReleaseFlag = cli.BoolFlag{
// 	Name:  "verify-release",
// 	Usage: "Verify to be a release version",
// }

// BumpFlag allows to increase the patch version
// F.e. 1.0.0 -> 1.0.1
// var BumpFlag = cli.BoolFlag{
// 	Name:  "bump",
// 	Usage: "Increase version",
// }

// BumpMinorFlag allows to increase the minor version
// F.e. 1.0.0 -> 1.1.0
// var BumpMinorFlag = cli.BoolFlag{
// 	Name:  "bump-minor",
// 	Usage: "Increase minor version",
// }

// BumpMajorFlag allows to increase the patch version
// F.e. 1.0.0 -> 2.0.0
// var BumpMajorFlag = cli.BoolFlag{
// 	Name:  "bump-major",
// 	Usage: "Increase major version",
// }

// GitBuildNumFlag allows to include build number into the version while being on the release/hotfix branch
// Build number is resolved from the $BUILD_NUMBER env variable
// 1.0.0 -> 1.0.0-rc.1 (release/1.0.0 branch)
// 1.0.0 -> 1.0.0 (defect/x branch)
var GitBuildNumFlag = cli.StringFlag{
	Name:  "git-build-num",
	Value: "rc.",
	Usage: "Include build number into the version",
}

// GitBuildNumBranchFlag allows to specify branch pattern for build-num calculations
var GitBuildNumBranchFlag = cli.StringSliceFlag{
	Name:  "git-build-num-branch",
	Usage: "Specify branch for git-build-num",
}

// GitShaFlag allows to include git sha into the version
// F.e. 1.0.0-SNAPSHOT -> 1.0.0-804cb4-SNAPSHOT
var GitShaFlag = cli.BoolFlag{
	Name:  "git-sha",
	Usage: "Inlcude git sha into the version",
}

// GitRefFlag allows to include git ref into the version
// F.e. 1.0.0-SNAPSHOT -> 1.0.0-defect-x-SNAPSHOT (branch defect/x)
// or 1.0.0-SNAPSHOT -> 1.0.0-develop-SNAPSHOT (branch develop)
var GitRefFlag = cli.BoolFlag{
	Name:  "git-ref",
	Usage: "Include git ref into the version",
}

// GitRefIgnoreFlag allows to skip some of the branches in git ref calculations
var GitRefIgnoreFlag = cli.StringSliceFlag{
	Name:  "git-ref-ignore",
	Usage: "Ignore branches for git-ref",
}

// GitVerifyNonDirtyFlag allows to verify git not to be dirty (throws error)
// var GitVerifyNonDirtyFlag = cli.BoolFlag{
// 	Name:  "verify-non-dirty",
// 	Usage: "Verify non dirty git directory",
//}

// GitDirtyFlag allows to mark version in case of dirty git
// If git is dirty, version will be changed accordingly
// F.e. 1.0.0 -> 1.0.0-dirty
// var GitDirtyFlag = cli.StringFlag{
// 	Name:  "dirty",
// 	Value: "dirty",
// 	Usage: "Include 'dirty' into the version",
// }

// GitDirtyTimestampFlag allows to append timestamp if git is dirty
// F.e. 1.0.0 -> 1.0.0-dirty-122310123
// var GitDirtyTimestampFlag = cli.BoolFlag{
// 	Name:  "dirty-timestamp",
// 	Usage: "Include timestamp into the version, when git is dirty",
// }

// TimestampFlag allows to append timestamp to the version
// F.e. 1.0.0 -> 1.0.0-122310123
// var TimestampFlag = cli.BoolFlag{
// 	Name:  "timestamp",
// 	Usage: "Include timestamp into the version",
// }

// AutPilotFlag allows to use one of the predefined configs
// F.e. --auto-pilot=app === --git-ref --git-build-num --dirty --dirty-timestamp
var AutPilotFlag = cli.StringFlag{
	Name:  "auto-pilot",
	Value: "app",
	Usage: "Use pre-defined configuration",
}
