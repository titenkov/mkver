# 🤖 semver
[![Build Status][icon_build]][page_build]
[![License][icon_license]](LICENSE)
> Calculates semantic version based on the branch and version taken from one of the sources (environment variable, gradle version, package.json, etc.). Version is calculated and validated according to https://semver.org/.

## Installation

### Brew

```bash
$ brew install titenkov/tap/semver
```

## Usage

```bash
$ semver --help
Usage:
  semver [command]

Available Commands:
  version     print calculated semantic version (default)
  help        help about any command

Flags:
  -h, --help   help for semver
  -d, --debug  print debug into to console

  --env        calculate semantic version based on the version taken from environment variable ($VERSION)
  --gradle     calculate semantic version based on the version taken from gradle.properties
  --npm        calculate semantic version based on the version taken from package.json
  
  --git-sha    includes sha into semantic version
  --git-ref    includes git ref into semantic version
  --dirty
  --dirty-timestamp
  --build-num

Use "semver [command] --help" for more information about a command.
```
## Examples


[icon_build]:      https://travis-ci.com/titenkov/semver.svg?branch=master
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg

[page_build]:      https://travis-ci.com/titenkov/semver
[page_promo]:      https://github.com/titenkov/semver
