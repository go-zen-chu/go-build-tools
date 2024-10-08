package gbt

import "fmt"

func KoInstall() error {
	return RunCmdWithLog("go install github.com/google/ko@latest")
}

// Build & publish image with ko
func KoPublish(importPath string) error {
	imageTag, err := GenerateImageTag()
	if err != nil {
		return fmt.Errorf("error generating image tag: %w", err)
	}
	// make sure you are logged in to the container registry
	// TIPS: ko build would add `-<md5>`` to image name without --base-import-paths flag
	koBuildCmd := fmt.Sprintf("ko build --base-import-paths --tags %s %s", imageTag, importPath)
	out, errMsg, err := RunLongRunningCmdWithLog(koBuildCmd)
	if err != nil {
		return fmt.Errorf("building image with ko: %w, stdout log: %s, stderr log: %s", err, out, errMsg)
	}
	return nil
}

// Build & publish latest tag with ko
func KoPublishLatest(importPath string) error {
	koBuildCmd := fmt.Sprintf("ko build --base-import-paths --tags latest %s", importPath)
	out, errMsg, err := RunLongRunningCmdWithLog(koBuildCmd)
	if err != nil {
		return fmt.Errorf("building image with ko: %w, stdout log: %s, stderr log: %s", err, out, errMsg)
	}
	return nil
}
