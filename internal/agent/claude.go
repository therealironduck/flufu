package agent

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	settingsPath     = ".claude/settings.local.json"
	settingsDirPerm  = 0o755
	settingsFilePerm = 0o644
)

type hookCommand struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type hookEntry struct {
	Matcher string        `json:"matcher,omitempty"`
	Hooks   []hookCommand `json:"hooks"`
}

type settingsPermissions struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

type settingsFile struct {
	Permissions *settingsPermissions   `json:"permissions,omitempty"`
	Hooks       map[string][]hookEntry `json:"hooks,omitempty"`
}

func selfCmd() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}

	return fmt.Sprintf("%s msg %d", exe, os.Getpid()), nil
}

func removeHooks() {
	ownCommand, err := selfCmd()
	if err != nil {
		return
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return
	}

	var s settingsFile
	if err := json.Unmarshal(data, &s); err != nil {
		return
	}

	filtered := make([]hookEntry, 0, len(s.Hooks["Stop"]))
	for _, entry := range s.Hooks["Stop"] {
		isOwn := false
		for _, h := range entry.Hooks {
			if h.Command == ownCommand {
				isOwn = true
				break
			}
		}
		if !isOwn {
			filtered = append(filtered, entry)
		}
	}

	s.Hooks["Stop"] = filtered

	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(settingsPath, out, settingsFilePerm)
}

func registerHooks() error {
	cmd, err := selfCmd()
	if err != nil {
		return err
	}

	var s settingsFile
	data, err := os.ReadFile(settingsPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read settings: %w", err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("parse settings: %w", err)
		}
	}

	if s.Hooks == nil {
		s.Hooks = make(map[string][]hookEntry)
	}

	// Replace any existing Stop hook with the current pid/exe.
	s.Hooks["Stop"] = []hookEntry{{
		Hooks: []hookCommand{{Type: "command", Command: cmd}},
	}}

	if err := os.MkdirAll(".claude", settingsDirPerm); err != nil {
		return fmt.Errorf("error creating .claude dir: %w", err)
	}

	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, out, settingsFilePerm); err != nil {
		return fmt.Errorf("error writing settings: %w", err)
	}

	return nil
}
