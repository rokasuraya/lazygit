package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/logs"
)

var (
	// Version is the current version of lazygit, set at build time via ldflags.
	version = "development"
	// buildDate is the date the binary was built, set at build time via ldflags.
	buildDate = "unknown"
	// commitSHA is the git commit SHA the binary was built from, set at build time via ldflags.
	commitSHA = "unknown"
)

func main() {
	// If we're being run as a daemon (e.g. for a background git operation),
	// handle that case and exit early.
	if daemon.TryRunningAsDaemon() {
		os.Exit(0)
	}

	buildInfo := &config.BuildInfo{
		Version:   version,
		BuildDate: buildDate,
		Commit:    commitSHA,
	}

	// Using upstream repo URL for update checks.
	appConfig, err := config.NewAppConfig(
		"lazygit",
		version,
		commitSHA,
		buildDate,
		"https://github.com/jesseduffield/lazygit",
		env.GetGitDirEnv() != "",
		buildInfo,
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error building app config: %v", err))
	}

	// Include the build date in the log context for easier debugging of
	// issues that may be version-specific.
	logger := logs.Global.
		WithField("version", version).
		WithField("buildDate", buildDate)

	lazygitApp, err := app.NewApp(appConfig, logger)
	if err != nil {
		if lazygitApp != nil {
			_ = lazygitApp.Close()
		}
		log.Fatal(err)
	}

	defer func() {
		if err := lazygitApp.Close(); err != nil {
			log.Printf("Error closing app: %v", err)
		}
	}()

	if err := lazygitApp.Run(); err != nil {
		log.Fatal(err)
	}
}
