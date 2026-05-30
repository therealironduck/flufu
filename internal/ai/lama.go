package ai

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func initLlama() error {
	libDir, err := resolveLibDir()
	if err != nil {
		return fmt.Errorf("cannot find lib dir: %w", err)
	}
	if err := llama.Load(libDir); err != nil {
		return fmt.Errorf("failed to load llama library from %q: %w", libDir, err)
	}

	llama.LogSet(llama.LogSilent())
	llama.Init()

	return nil
}

func resolveLibDir() (string, error) {
	if v := os.Getenv("YZMA_LIB"); v != "" {
		return v, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not load executable path: %w", err)
	}

	return filepath.Join(filepath.Dir(exe), "lib"), nil
}
