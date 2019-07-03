# 🤖 mkver
[![Build Status][icon_build]][page_build]
[![License][icon_license]](LICENSE)
> Calculates application version by enriching the original one with various information

## Installation

### Brew

```bash
$ brew install titenkov/tap/mkver
```

## Usage

```bash
$ mkver --help
Usage:
  mkver [flags]

Flags:
  -h, --help                help for mkver

  --env                     resolve version from env variable
  --gradle                  resolve version from gradle properties
  
  --git-sha                 include git sha into the version
  --git-ref                 include git ref into the version
  --git-ref-ignore          exclude branches using regexp from git ref calculation
  --git-build-num           include build number into the version
  --git-build-num-branch    specify branches using regexp for build num calculation
  
```
## Examples

```bash
mkver --git-ref --git-sha --git-ref-ignore=^develop$ --git-ref-ignore=^master$ --git-ref-ignore=^release --git-build-num=rc. --git-build-num-branch=^release.+$ 
```

[icon_build]:      https://travis-ci.com/titenkov/mkver.svg?branch=master
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg

[page_build]:      https://travis-ci.com/titenkov/mkver
[page_promo]:      https://github.com/titenkov/mkver
