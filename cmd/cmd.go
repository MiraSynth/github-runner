package cmd

import (
	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/github"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"mirasynth.stream/github-runner/internal/atlas"
)

var rootCmd *cobra.Command

var configFilePath string
var githubClient *github.Client

func init() {
	rootCmd = &cobra.Command{
		Use:   "githubrunner",
		Short: atlas.GITHUBRUNNER_SHORT_DESC,
		Long:  atlas.GITHUBRUNNER_LONG_DESC,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetOutput(os.Stdout)
			log.SetLevel(log.DebugLevel)

			err := config.SetupConfig(configFilePath)
			if err != nil {
				return err
			}

			gc, err := github.GetClient(cmd.Context(), nil)
			if err != nil {
				return err
			}

			githubClient = &gc

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Sets the path to where the config file is loaded from")

	rootCmd.AddCommand(NewServerCmd())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
