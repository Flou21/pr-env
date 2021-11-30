package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	Open        = "opened"
	Closed      = "closed"
	Synchronize = "synchronize"
)

func main() {
	err := godotenv.Load()
	log.Info().Msgf("hello creating or deleting pr env now")
	if err != nil {
		log.Warn().Err(err).Msgf("error reading .env file")
	}

	config := Config{}

	readConfig(&config)

	if config.State == Open {
		CloneRepository()
		CreateNewEnv(&config)
		err := PushEnv()
		if err != nil {
			log.Fatal().Err(err).Msgf("something went wrong when uploading changes to git")
		}
	}

	if config.State == Closed {
		CloneRepository()
		DeleteOldEnv(&config)
		PushEnv()
	}

	if config.State == Synchronize {
		CloneRepository()
		UpdateEnv()
	}

	log.Info().Msgf("finished pr env stuff")
}

func readConfig(config *Config) {
	config.State = readOsVar("PULL_REQUEST_NEW_STATE")
	config.GitUsername = readOsVar("GIT_USERNAME")
	config.GitPassword = readOsVar("GIT_TOKEN")

	tmpNumber := os.Getenv("PULL_REQUEST_NUMBER")

	num, err := strconv.Atoi(tmpNumber)

	if err != nil {
		log.Error().Err(err).Msgf("could not parse number")
	}

	config.Number = num

	readTemplateData(config)
}

func readTemplateData(config *Config) {

	config.TemplateData = &TemplateData{}

	containerImage := readOsVar("CONTAINER_IMAGE")
	containerPort := readOsVar("CONTAINER_PORT")
	namespace := readOsVar("K8S_NAMESPACE")
	pullRequestNumber := readOsVar("PULL_REQUEST_NUMBER")

	config.TemplateData.TemplateDeploymentName = "pr-env-" + namespace + "-" + pullRequestNumber
	config.TemplateData.TemplateNamespace = namespace
	config.TemplateData.TemplateContainerImage = containerImage
	config.TemplateData.TemplateContainerTag = pullRequestNumber
	config.TemplateData.TemplateContainerName = "pr-env-" + namespace + "-" + pullRequestNumber
	config.TemplateData.TemplatePort = containerPort

}

func BuildManifestFilename(config *Config) string {
	return config.TemplateData.TemplateNamespace + "/" + "pr-env-" + config.TemplateData.TemplateNamespace + "-" + config.TemplateData.TemplateContainerTag + ".yaml"
}

type Config struct {
	Number        int
	State         string
	Template      string
	TemplateData  *TemplateData
	GitRepository string
	GitUsername   string
	GitPassword   string
}
