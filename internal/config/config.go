package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"mirasynth.stream/github-runner/internal/atlas"
)

var venv *viper.Viper

func SetupConfig() error {
	venv = viper.New()
	venv.SetConfigType(atlas.CONFIG_TYPE)
	venv.SetEnvPrefix(atlas.CONFIG_PREFIX)
	venv.SetEnvKeyReplacer(strings.NewReplacer(".", "_", " ", ""))
	venv.AutomaticEnv()

	configFilePath, err := verifyConfigFile()
	if err != nil {
		return err
	}

	venv.SetConfigFile(configFilePath)

	err = venv.ReadInConfig()
	if err != nil {
		return err
	}

	log.Debug("usign config", venv.ConfigFileUsed())

	return nil
}

func verifyConfigFile() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configDirPath := path.Join(userConfigDir, atlas.CONFIG_NAMESPACE, atlas.CONFIG_PREFIX)
	configFilePath := path.Join(configDirPath, fmt.Sprintf("%s.%s", atlas.CONFIG_FILENAME, atlas.CONFIG_TYPE))

	err = os.MkdirAll(configDirPath, 0755)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(configFilePath)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return configFilePath, err
	}

	bytes := []byte{}
	err = os.WriteFile(configFilePath, bytes, 0755)
	if err != nil {
		return "", err
	}

	return configFilePath, nil
}

func GetGitHubInstallationId() string {
	return venv.GetString("github.installationid")
}

func GetGitHubClientId() string {
	return venv.GetString("github.clientid")
}

func GetGitHubSecret() string {
	return venv.GetString("github.secret")
}

func GetGitHubWebhookSecret() string {
	return venv.GetString("github.webhook.secret")
}
