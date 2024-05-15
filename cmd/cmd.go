package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"mirasynth.stream/github-runner/internal/atlas"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "githubrunner [command]",
		Short: atlas.GITHUBRUNNER_SHORT_DESC,
		Long:  atlas.GITHUBRUNNER_LONG_DESC,
	}

	rootCmd.AddCommand(NewServerhCmd())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
