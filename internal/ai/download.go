package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/therealironduck/flufu/internal/config"
)

var errInvalidStatusCode = errors.New("invalid status code")

func loadModel(ctx context.Context) (string, error) {
	path, err := config.GetConfigPath(modelFileName)
	if err != nil {
		return "", fmt.Errorf("failed to load config path: %w", err)
	}

	err = ensureModel(ctx, path)
	if err != nil {
		return path, fmt.Errorf("model setup failed: %w", err)
	}

	return path, nil
}

func ensureModel(ctx context.Context, path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := downloadFile(ctx, path, modelURL); err != nil {
		return err
	}

	return nil
}

func downloadFile(ctx context.Context, dest, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not download model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d - %w", resp.StatusCode, errInvalidStatusCode)
	}

	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("could not create tmp file: %w", err)
	}

	defer func() {
		f.Close()
		os.Remove(tmp)
	}()

	const bufSize = 32 * 1024
	buf := make([]byte, bufSize)

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return fmt.Errorf("failed to write content: %w", werr)
			}
		}

		if errors.Is(readErr, io.EOF) {
			break
		}

		if readErr != nil {
			return fmt.Errorf("failed to write content: %w", readErr)
		}
	}

	f.Close()

	err = os.Rename(tmp, dest)
	if err != nil {
		return fmt.Errorf("failed to rename tmp file: %w", err)
	}

	return nil
}
