package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// This executes the supplied cmd by /bin/sh and returns an error if it returns
// unexpectedly
func System(cmdString string) error {
	return SystemV(cmdString, []string{}, os.Stderr, os.Stderr)
}

func SystemV(cmdString string, extraEnviron []string, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	cmd.Env = append(os.Environ(), extraEnviron...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func callPreProcessingHook() {
	if appConfig.PreProcessingHook != "" {
		LOG.Println("### Calling PreProcessingHook ...")

		err := System(appConfig.PreProcessingHook)
		if err != nil {
			LOG.Printf("PreProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callPostProcessingHook() {
	if appConfig.PostProcessingHook != "" {
		LOG.Println("\n### Calling PostProcessingHook ...")

		err := System(appConfig.PostProcessingHook)
		if err != nil {
			LOG.Printf("PostProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callEpisodeHook(episodePath, seriesName string) {
	if appConfig.EpisodeHook != "" {
		LOG.Println("# Calling EpisodeHook ...")

		hookCmd := fmt.Sprintf("%s \"%s\" \"%s\"",
			appConfig.EpisodeHook, episodePath, seriesName)

		err := System(hookCmd)
		if err != nil {
			LOG.Printf("EpisodeHook ended with an error: %s\n", err)
		}
	}
}
