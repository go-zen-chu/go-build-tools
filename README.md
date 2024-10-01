# go-build-tools

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/go-build-tools)](http://pkg.go.dev/github.com/go-zen-chu/go-build-tools)
[![Actions Status](https://github.com/go-zen-chu/go-build-tools/workflows/check-pr/badge.svg)](https://github.com/go-zen-chu/go-build-tools/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/go-build-tools.svg)](https://github.com/go-zen-chu/go-build-tools/issues)

Build tool repository written in Golang.

## Goal

- Use from [magefile](https://magefile.org/) and easily realize development

## Build sample

```console
$ mage -d mage 
Targets:
  dockerBuildLatest               builds the docker image with the latest tag locally
  dockerBuildPublishWithGenTag    DockerBuildPublishLatest builds and publishes the docker image with generated tag
  dockerLogin                     logs in to the docker registry
  dockerPublishLatest             publishes the docker image with the latest tag
  installDevTools                 installs required development tools for this project
  koPublish                       builds and publishes the image with ko generated tag
  koPublishLatest                 builds and publishes the image with the latest tag
  updateFormula                   updates formula with current version for homebrew tap
```

## Test

```console
go test -v ./...
```
