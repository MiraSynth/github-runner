package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"mirasynth.stream/github-runner/cmd"
	"mirasynth.stream/github-runner/internal/config"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	err := config.SetupConfig()

	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

	cmd.Execute()
}
