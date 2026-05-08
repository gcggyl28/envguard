// Package pinner provides functionality to pin (freeze) env variable values
// by recording their current state and detecting drift on subsequent runs.
package pinner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// PinnedEnv represents a frozen snapshot of specific env keys.
type PinnedEnv struct {
	PinnedAt time.Time         `json:"pinned_at"`
	Keys     map[string]string `json:"keys"`
}

// DriftResult holds the result of comparing live env values against pinned ones.
type DriftResult struct {
	Changed []DriftEntry
	Removed []string
}

// DriftEntry describes a single key whose value has drifted.
type DriftEntry struct {
	Key      string
	Pinned   string
	Current  string
}

// Pin records the values for the given keys from env and writes them to path.
func Pin(env map[string]string, keys []string, path string) (*PinnedEnv, error) {
	pinned := &PinnedEnv{
		PinnedAt: time.Now().UTC(),
		Keys:     make(map[string]string, len(keys)),
	}
	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			return nil, fmt.Errorf("pinner: key %q not found in env", k)
		}
		pinned.Keys[k] = v
	}
	data, err := json.MarshalIndent(pinned, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("pinner: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, fmt.Errorf("pinner: write %s: %w", path, err)
	}
	return pinned, nil
}

// Load reads a previously pinned env file from path.
func Load(path string) (*PinnedEnv, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pinner: read %s: %w", path, err)
	}
	var p PinnedEnv
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("pinner: unmarshal: %w", err)
	}
	return &p, nil
}

// CheckDrift compares pinned values against the current env map.
func CheckDrift(pinned *PinnedEnv, env map[string]string) DriftResult {
	var result DriftResult
	keys := make([]string, 0, len(pinned.Keys))
	for k := range pinned.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		current, ok := env[k]
		if !ok {
			result.Removed = append(result.Removed, k)
			continue
		}
		if current != pinned.Keys[k] {
			result.Changed = append(result.Changed, DriftEntry{
				Key:     k,
				Pinned:  pinned.Keys[k],
				Current: current,
			})
		}
	}
	return result
}
