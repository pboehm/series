package main

import (
	"fmt"
	"os"
	"os/exec"
)

// This executes the supplied cmd by /bin/sh and returns an error if it returns
// unexpectedly
func System(cmdString string) error {

	cmd := exec.Command("/bin/sh", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func callPreProcessingHook() {
	if appConfig.PreProcessingHook != "" {
		fmt.Println("### Calling PreProcessingHook ...")

		err := System(appConfig.PreProcessingHook)
		if err != nil {
			fmt.Printf("PreProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callPostProcessingHook() {
	if appConfig.PostProcessingHook != "" {
		fmt.Println("\n### Calling PostProcessingHook ...")

		err := System(appConfig.PostProcessingHook)
		if err != nil {
			fmt.Printf("PostProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callEpisodeHook(episodePath, seriesName string) {
	if appConfig.EpisodeHook != "" {
		fmt.Println("# Calling EpisodeHook ...")

		hookCmd := fmt.Sprintf("%s \"%s\" \"%s\"",
			appConfig.EpisodeHook, episodePath, seriesName)

		err := System(hookCmd)
		if err != nil {
			fmt.Printf("EpisodeHook ended with an error: %s\n", err)
		}
	}
}
