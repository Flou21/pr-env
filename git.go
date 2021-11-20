package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
)

var repository *git.Repository

func CloneRepository() error {

	var err error
	username := os.Getenv("GIT_USERNAME")
	token := os.Getenv("GIT_TOKEN")
	url := "https://github.com/Coflnet/kube.git"
	auth := &http.BasicAuth{Username: username, Password: token}
	repository, err = git.PlainClone("./kube", false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal().Err(err).Msgf("can not clone the kube repository")
		return err
	}

	return nil
}

func PushEnv() error {
	if repository == nil {
		log.Error().Msgf("the repository variable is nil")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		log.Error().Err(err).Msgf("error when getting the worktree")
		return err
	}

	worktree.Pull(&git.PullOptions{})

	worktree.Add(".")

	_, err = worktree.Commit("[CI] update temp env ", &git.CommitOptions{
		All: true,
		Committer: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
		Author: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Error().Err(err).Msgf("something went wrong when committing")
		return err
	}

	username := os.Getenv("GIT_USERNAME")
	token := os.Getenv("GIT_TOKEN")
	auth := &http.BasicAuth{Username: username, Password: token}
	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth:       auth,
	})

	if err != nil {
		log.Error().Err(err).Msgf("something went wrong when pushing")
		return err
	}

	return nil
}

func AddToRepository(filename string, content string, config *Config) error {

	log.Info().Msgf("git repository: %s", config.GitRepository)
	log.Info().Msgf("git username: %s", config.GitUsername)

	repoURL := fmt.Sprintf("https://%s:%s@%s", config.GitUsername, config.GitPassword, config.GitRepository)

	const repoPath = "/tmp/repo_create"
	log.Info().Msgf("clone repository")
	repository, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:        repoURL,
		Progress:   os.Stdout,
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	log.Info().Msgf("getting the worktree from the git repository")
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	// pull to make sure the newest files are here
	worktree.Pull(&git.PullOptions{})

	log.Info().Msgf("writing content to file in git repository")
	combinedFilename := filepath.Join(repoPath, filename)
	err = os.WriteFile(combinedFilename, []byte(content), 0644)
	if err != nil {
		return err
	}
	log.Info().Msgf("wrote the files to " + combinedFilename)

	log.Info().Msgf("going to add the new file to git")
	log.Info().Msgf("adding file" + filename)
	worktree.Add(filename)

	log.Info().Msgf("goingt to commit the new file to git")

	_, err = worktree.Commit("[CI] update add file: "+filename, &git.CommitOptions{
		All: true,
		Committer: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
		Author: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	log.Info().Msgf("going to push the new commit")
	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}

func RemoveFromRepository(filename string, config *Config) error {
	repoURL := fmt.Sprintf("https://%s:%s@%s", config.GitUsername, config.GitPassword, config.GitRepository)
	const repoPath = "/tmp/repo_delete"
	log.Info().Msgf("clone repository")
	repository, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:        repoURL,
		Progress:   os.Stdout,
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	log.Info().Msgf("getting the worktree from the git repository")
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	log.Info().Msgf("writing content to file in git repository")
	combinedFilename := filepath.Join(repoPath, filename)

	_, err = os.Stat(combinedFilename)
	if os.IsNotExist(err) {
		log.Info().Msgf("file does not exist there don't delete it")
		return nil
	}

	err = os.Remove(combinedFilename)
	if err != nil {
		return err
	}

	log.Info().Msgf("going to add the new file to git")
	log.Info().Msgf("adding file" + filename)
	worktree.Add(filename)

	log.Info().Msgf("goingt to commit the new file to git")

	_, err = worktree.Commit("[CI] update remove file: "+filename, &git.CommitOptions{
		All: true,
		Committer: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
		Author: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	username := os.Getenv("GIT_USERNAME")
	token := os.Getenv("GIT_TOKEN")
	auth := &http.BasicAuth{Username: username, Password: token}
	log.Info().Msgf("going to push the new commit")
	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth:       auth,
	})

	if err != nil {
		return err
	}

	return nil
}
