//go:build mage
// +build mage

package main

import (
	"fmt"
	"log/slog"
	"os"

	gbt "github.com/go-zen-chu/go-build-tools"
)

const currentVersion = "1.0.0"
const currentTagVersion = "v" + currentVersion

const imageRegistry = "amasuda"
const repository = "go-build-tools"
const dockerFileLocation = "."

const tapOwner = "go-zen-chu"
const tapRepo = "homebrew-tools"

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

/*=======================
setup
=======================*/

// InstallDevTools installs required development tools for this project
func InstallDevTools() error {
	outMsg, errMsg, err := gbt.RunLongRunningCmdWithLog("go install go.uber.org/mock/mockgen@latest")
	if err != nil {
		return fmt.Errorf("installing mockgen: %w\nstdout: %s\nstderr: %s\n", err, outMsg, errMsg)
	}
	outMsg, errMsg, err = gbt.RunLongRunningCmdWithLog("brew install mage")
	if err != nil {
		return fmt.Errorf("installing mockgen: %w\nstdout: %s\nstderr: %s\n", err, outMsg, errMsg)
	}
	return nil
}

/*=======================
workflow
=======================*/

// KoPublish builds and publishes the image with ko generated tag
func KoPublish() error {
	return gbt.KoPublish(".")
}

// KoPublishLatest builds and publishes the image with the latest tag
func KoPublishLatest() error {
	return gbt.KoPublishLatest(".")
}

/*=======================
examples. not used in this repo
=======================*/

// DockerLogin logs in to the docker registry
func DockerLogin() error {
	return gbt.DockerLogin()
}

// DockerBuildLatest builds the docker image with the latest tag locally
func DockerBuildLatest() error {
	return gbt.DockerBuildLatest(imageRegistry, repository, dockerFileLocation)
}

// DockerPublishLatest publishes the docker image with the latest tag
func DockerPublishLatest() error {
	return gbt.DockerPublishLatest(imageRegistry, repository)
}

// DockerBuildPublishLatest builds and publishes the docker image with generated tag
func DockerBuildPublishWithGenTag() error {
	return gbt.DockerBuildPublishGeneratedImageTag(imageRegistry, repository, dockerFileLocation)
}

const formulaTemplate = `class GoBuildTools < Formula
    desc "Go build tools of go-zen-chu"
    homepage "https://github.com/go-zen-chu/go-build-tools"
    version "%[1]s"
    
    on_macos do
        if Hardware::CPU.arm?
            url "https://github.com/go-zen-chu/go-build-tools/releases/download/v%[1]s/go-build-tools_Darwin_arm64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinArm64}}"
        else
            url "https://github.com/go-zen-chu/go-build-tools/releases/download/v%[1]s/go-build-tools_Darwin_x86_64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinX86_64}}"
        end
    end

    def install
        bin.install "go-build-tools"
    end

    test do
        system "#{bin}/go-build-tools", "--help"
    end
end
`

// UpdateFormula updates formula with current version for homebrew tap
func UpdateFormula() error {
	ft := fmt.Sprintf(formulaTemplate, currentVersion)
	return gbt.GenerateFormula(ft,
		tapOwner,
		tapRepo,
		"go-zen-chu",
		repository,
		currentTagVersion,
	)
}
