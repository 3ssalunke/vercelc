package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/3ssalunke/vercelc/shared/util"
)

func BuildProject(projectId string) error {
	folderPath, err := util.GetPathForFolder(fmt.Sprintf("build/output/%s", projectId))
	if err != nil {
		return fmt.Errorf("failed to get path for build folder: %v", err)
	}

	_, err = os.Stat(folderPath)
	if err != nil {
		return fmt.Errorf("path for build folder does not exist: %v", err)
	}

	packageJSONPath := filepath.Join(folderPath, "package.json")
	_, err = os.Stat(packageJSONPath)
	if err != nil {
		return fmt.Errorf("package.json does not exist: %v", err)
	}

	err = os.Chdir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	cmd := exec.Command("npm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed execute install command: %v", err)
	}

	cmd = exec.Command("npm", "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed execute build command: %v", err)
	}

	log.Printf("code built successfully for project id %s", projectId)
	return nil
}
