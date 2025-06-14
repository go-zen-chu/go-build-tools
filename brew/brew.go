package brew

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/go-zen-chu/go-build-tools/github"
	"github.com/go-zen-chu/go-build-tools/util"
)

type BrewFormula struct {
	ChecksumSHA256DarwinArm64  %%string
	ChecksumSHA256DarwinX86_64 string
	ChecksumSHA256LinuxArm64   string
	ChecksumSHA256LinuxX86_64  string
}

func fileOrDirExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
}

// GenerateFormula generate a Homebrew formula string based on the provided template and artifact information.
func GenerateFormula(formulaTemplate, artifactOwner, artifactRepo, gitTag string) (string, error) {
	tmpl, err := template.New("formula").Parse(formulaTemplate)
	if err != nil {
		return "", fmt.Errorf("parse formula template: %w", err)
	}

	var httpClient http.Client
	release, err := github.GetTagRelease(&httpClient, artifactOwner, artifactRepo, gitTag)
	if err != nil {
		return "", fmt.Errorf("get latest release: %w", err)
	}
	checksumMap, err := github.GetChecksumMap(&httpClient, release)
	if err != nil {
		return "", fmt.Errorf("get checksum map: %w", err)
	}
	bf := BrewFormula{}
	for filename, checksum := range checksumMap {
		if strings.Contains(filename, "Darwin") {
			if strings.Contains(filename, "arm64") {
				bf.ChecksumSHA256DarwinArm64 = checksum
			} else if strings.Contains(filename, "x86_64") {
				bf.ChecksumSHA256DarwinX86_64 = checksum
			}
		}
		if strings.Contains(filename, "Linux") {
			if strings.Contains(filename, "arm64") {
				bf.ChecksumSHA256DarwinArm64 = checksum
			} else if strings.Contains(filename, "x86_64") {
				bf.ChecksumSHA256DarwinX86_64 = checksum
			}
		}
	}

	var bb bytes.Buffer
	err = tmpl.Execute(&bb, bf)
	if err != nil {
		return "", fmt.Errorf("execute formula template: %w", err)
	}
	return bb.String(), nil
}

// PushFormula pushes the generated Homebrew formula to the specified tap repository on GitHub.
func PushFormula(formula string, tapOwner, tapRepo, artifactName, gitTag string) error {
	if fileOrDirExists(tapRepo) {
		return fmt.Errorf("tap repo %s already exists. please remove before push", tapRepo)
	}
	tapRepoUrl := fmt.Sprintf("https://github.com/%s/%s.git", tapOwner, tapRepo)
	outMsg, errMsg, err := util.RunLongRunningCmdWithLog(fmt.Sprintf("git clone %s", tapRepoUrl))
	if err != nil {
		return fmt.Errorf("git clone: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}

	formulaFilePath := fmt.Sprintf("%s/Formula/%s.rb", tapRepo, artifactName)
	err = os.WriteFile(formulaFilePath, []byte(formula), 0644)
	if err != nil {
		return fmt.Errorf("write formula file to %s: %w", formulaFilePath, err)
	}

	commitMsg := fmt.Sprintf("update formula to tag: %s", gitTag)
	if err := github.GitHubActionPush(tapRepo, tapOwner, tapRepo, commitMsg); err != nil {
		return fmt.Errorf("push formula: %w", err)
	}
	err = os.RemoveAll(tapRepo)
	if err != nil {
		return fmt.Errorf("remove tap repo %s: %w", tapRepo, err)
	}
	return nil
}
