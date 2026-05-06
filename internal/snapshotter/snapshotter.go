package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's key-value pairs.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot of the given env map to the specified file path as JSON.
func Save(env map[string]string, source, dest string) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Env:       env,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal failed: %w", err)
	}

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return fmt.Errorf("snapshotter: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshotter: unmarshal failed: %w", err)
	}

	return &snap, nil
}

// Compare returns keys that were added, removed, or changed between two snapshots.
func Compare(base, current *Snapshot) (added, removed, changed []string) {
	for k, v := range current.Env {
		if bv, ok := base.Env[k]; !ok {
			added = append(added, k)
		} else if bv != v {
			changed = append(changed, k)
		}
	}

	for k := range base.Env {
		if _, ok := current.Env[k]; !ok {
			removed = append(removed, k)
		}
	}

	return added, removed, changed
}
