package main

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func readOsVar(key string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		log.Fatal().Msg("the env var " + key + " is not given")
	}
	return tmp
}

func ExecuteOsCommand(command string) (string, error) {
	log.Info().Msgf("executing command: %s", command)
	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		return "", err
	}
	log.Info().Msgf("result from command: %s", string(output))
	return string(output), nil
}

type Kustomization struct {
	Bases           []string          `yaml:"bases"`
	PatchesJSON6902 []PatchesJSON6902 `yaml:"patchesJson6902"`
}

type PatchesJSON6902 struct {
	Target Target `yaml:"target"`
	Path   string `yaml:"path"`
}

type Target struct {
	Group   string `yaml:"group"`
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`
	Name    string `yaml:"name"`
}

type ValuesFile struct {
	Op    string `yaml:"op"`
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
}
