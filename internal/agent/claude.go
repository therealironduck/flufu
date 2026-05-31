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

func registerHooks() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	pid := os.Getpid()

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

	// Replace any existing idle_prompt Notification hook with the current pid/exe.
	existing := s.Hooks["Stop"]
	filtered := make([]hookEntry, 0, len(existing))
	for _, entry := range existing {
		if entry.Matcher != "idle_prompt" {
			filtered = append(filtered, entry)
		}
	}

	filtered = append(filtered, hookEntry{
		Matcher: "",
		Hooks: []hookCommand{
			{
				Type:    "command",
				Command: fmt.Sprintf("%s msg %d", exe, pid),
			},
		},
	})
	s.Hooks["Stop"] = filtered

	if err := os.MkdirAll(".claude", settingsDirPerm); err != nil {
		return fmt.Errorf("create .claude dir: %w", err)
	}

	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, out, settingsFilePerm); err != nil {
		return fmt.Errorf("write settings: %w", err)
	}

	return nil
}
