package main

import (
	"fmt"
	"os"
	"os/exec"
)

// This executes the supplied cmd by /bin/sh and returns an error if it returns
// unexpectedly
func System(cmd_string string) error {

	cmd := exec.Command("/bin/sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func callPreProcessingHook() {
	if AppConfig.PreProcessingHook != "" {
		fmt.Println("### Calling PreProcessingHook ...")

		err := System(AppConfig.PreProcessingHook)
		if err != nil {
			fmt.Printf("PreProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callPostProcessingHook() {
	if AppConfig.PostProcessingHook != "" {
		fmt.Println("\n### Calling PostProcessingHook ...")

		err := System(AppConfig.PostProcessingHook)
		if err != nil {
			fmt.Printf("PostProcessingHook ended with an error: %s\n", err)
		}
	}
}

func callEpisodeHook(episode_path, series_name string) {
	if AppConfig.EpisodeHook != "" {
		fmt.Println("# Calling EpisodeHook ...")

		hook_cmd := fmt.Sprintf("%s \"%s\" \"%s\"",
			AppConfig.EpisodeHook, episode_path, series_name)

		err := System(hook_cmd)
		if err != nil {
			fmt.Printf("EpisodeHook ended with an error: %s\n", err)
		}
	}
}
