package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

func DeleteOldEnv(config *Config) {
	log.Info().Msg("going to delete the old environment")

	folderName := config.TemplateData.TemplateDeploymentName
	path := "./kube/pr-env/overlays/" + folderName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Info().Msgf("path %s does not exist", path)
		return
	}

	err := os.RemoveAll(path)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when deleting path %s", path)
	}
}
