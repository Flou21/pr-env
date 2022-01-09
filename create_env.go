package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func CreateNewEnv(config *Config) {
	log.Info().Msg("going to create a new environment")

	log.Info().Msg("generating new manifest")
	err := GenerateNewValuesFile(config.Template, config)

	if err != nil {
		log.Error().Err(err).Msgf("error when generating new values file")
		log.Fatal().Msgf("stop program")
	}
}

func GenerateNewValuesFile(template string, config *Config) error {

	folderName := config.TemplateData.TemplateDeploymentName
	path := "./kube/pr-env/overlays/" + folderName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Error().Err(err).Msgf("could not create directory %s", folderName)
			return err
		}
	}

	err := CreateKustomization(path, folderName)
	if err != nil {
		log.Error().Err(err).Msgf("could not create kustomization.yaml file")
	}

	err = CreateValuesFile(path, folderName, config)
	if err != nil {
		log.Error().Err(err).Msgf("could not create values file")
	}

	log.Info().Msgf("finsihed")
	return nil
}

func CreateKustomization(path string, folderName string) error {
	filename := path + "/" + "kustomization.yaml"
	kustomization := Kustomization{
		Bases: []string{"../../base"},
		PatchesJSON6902: []PatchesJSON6902{
			{
				Path: folderName + "-deployment.yaml",
				Target: Target{
					Group:   "apps",
					Version: "v1",
					Kind:    "Deployment",
					Name:    "deployment",
				},
			},
			{
				Path: folderName + "-service.yaml",
				Target: Target{
					Version: "v1",
					Kind:    "Service",
					Name:    "service",
				},
			},
			{
				Path: folderName + "-ingress.yaml",
				Target: Target{
					Version: "v1",
					Kind:    "Ingress",
					Name:    "ingress",
				},
			},
		},
	}

	output, err := yaml.Marshal(&kustomization)
	if err != nil {
		log.Error().Err(err).Msgf("error when converting kustomization struct to yaml")
		return err
	}

	err = os.WriteFile(filename, output, 0777)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when saving kustomization.yaml file %s", filename)
		return err
	}

	return nil
}

func CreateValuesFile(path string, folderName string, config *Config) error {
	err := CreateDeploymentFile(path, folderName, config)
	if err != nil {
		log.Error().Err(err).Msgf("there happened an error when creating the deployment patch")
	}

	err = CreateServiceFile(path, folderName, config)
	if err != nil {
		log.Error().Err(err).Msgf("there happened an error when creating the service patch")
	}

	err = CreateIngressFile(path, folderName, config)
	if err != nil {
		log.Error().Err(err).Msgf("there happened an error when creating the ingress patch")
	}

	return nil
}

func CreateServiceFile(path string, folderName string, config *Config) error {

	filename := path + "/" + folderName + ".yaml"

	valuesFile := []ValuesFile{
		{
			Op:    "replace",
			Path:  "/metadata/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/metadata/namespace",
			Value: config.TemplateData.TemplateNamespace,
		},
		{
			Op:    "replace",
			Path:  "/spec/selector/app",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/ports/0/port",
			Value: config.TemplateData.TemplatePort,
		},
		{
			Op:    "replace",
			Path:  "/spec/ports/0/targetPort",
			Value: config.TemplateData.TemplatePort,
		},
	}

	output, err := yaml.Marshal(&valuesFile)
	if err != nil {
		log.Error().Err(err).Msgf("error when converting kustomization struct to yaml")
		return err
	}

	err = os.WriteFile(filename, output, 0777)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when saving kustomization.yaml file %s", filename)
		return err
	}

	return nil
}

func CreateDeploymentFile(path string, folderName string, config *Config) error {

	filename := path + "/" + folderName + ".yaml"

	valuesFile := []ValuesFile{
		{
			Op:    "replace",
			Path:  "/metadata/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/metadata/namespace",
			Value: config.TemplateData.TemplateNamespace,
		},
		{
			Op:    "replace",
			Path:  "/metadata/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/metadata/labels/app",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/selector/matchLabels/app",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/template/metadata/labels/app",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/template/spec/containers/0/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/template/spec/containers/0/image",
			Value: config.TemplateData.TemplateContainerImage,
		},
	}

	output, err := yaml.Marshal(&valuesFile)
	if err != nil {
		log.Error().Err(err).Msgf("error when converting kustomization struct to yaml")
		return err
	}

	err = os.WriteFile(filename, output, 0777)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when saving kustomization.yaml file %s", filename)
		return err
	}

	return nil
}

func CreateIngressFile(path string, folderName string, config *Config) error {

	filename := path + "/" + folderName + ".yaml"

	valuesFile := []ValuesFile{
		{
			Op:    "replace",
			Path:  "/metadata/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/metadata/namespace",
			Value: config.TemplateData.TemplateNamespace,
		},
		{
			Op:    "replace",
			Path:  "/spec/rules/0/http/paths/0/path",
			Value: "/" + folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/rules/0/http/paths/0/backend/service/name",
			Value: folderName,
		},
		{
			Op:    "replace",
			Path:  "/spec/rules/0/http/paths/0/backend/service/port/number",
			Value: config.TemplateData.TemplatePort,
		},
	}

	output, err := yaml.Marshal(&valuesFile)
	if err != nil {
		log.Error().Err(err).Msgf("error when converting kustomization struct to yaml")
		return err
	}

	err = os.WriteFile(filename, output, 0777)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when saving kustomization.yaml file %s", filename)
		return err
	}

	return nil
}
