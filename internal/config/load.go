package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func GetConfigPath(subPath string) (string, error) {
	path, err := xdg.ConfigFile(filepath.Join("flufu", subPath))
	if err == nil {
		return path, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot determine executable path: %w", err)
	}

	modelPath := filepath.Join(filepath.Dir(exe), subPath)
	return modelPath, nil
}
