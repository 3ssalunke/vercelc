package util

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ID_LENGTH = 13
	CHARSET   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func CloneRepo(repoUrl, destination string) error {
	cmd := exec.Command("git", "clone", repoUrl, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func GenerateRandomId() string {
	b := make([]byte, ID_LENGTH)
	for i := range b {
		b[i] = CHARSET[rand.Intn(len(CHARSET))]
	}
	return string(b)
}

func GetLocalCloneFolder(projectId string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	destination := filepath.Join(currentDir, "output", projectId)
	return destination, nil
}
