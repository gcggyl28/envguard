package freezer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// FrozenEnv represents a frozen snapshot of an env map with metadata.
type FrozenEnv struct {
	FrozenAt  time.Time         `json:"frozen_at"`
	Source    string            `json:"source"`
	Keys      []string          `json:"keys"`
	Values    map[string]string `json:"values"`
}

// FreezeResult describes the outcome of a freeze operation.
type FreezeResult struct {
	Frozen  []string
	Skipped []string
	File    string
}

// Freeze writes a frozen copy of env to outFile, optionally restricting to allowKeys.
// If allowKeys is empty, all keys are frozen.
func Freeze(env map[string]string, source, outFile string, allowKeys []string) (FreezeResult, error) {
	allowSet := make(map[string]bool, len(allowKeys))
	for _, k := range allowKeys {
		allowSet[k] = true
	}

	values := make(map[string]string)
	var frozen, skipped []string

	for k, v := range env {
		if len(allowSet) > 0 && !allowSet[k] {
			skipped = append(skipped, k)
			continue
		}
		values[k] = v
		frozen = append(frozen, k)
	}

	sort.Strings(frozen)
	sort.Strings(skipped)

	fe := FrozenEnv{
		FrozenAt: time.Now().UTC(),
		Source:   source,
		Keys:     frozen,
		Values:   values,
	}

	data, err := json.MarshalIndent(fe, "", "  ")
	if err != nil {
		return FreezeResult{}, fmt.Errorf("freezer: marshal: %w", err)
	}

	if err := os.WriteFile(outFile, data, 0o644); err != nil {
		return FreezeResult{}, fmt.Errorf("freezer: write %s: %w", outFile, err)
	}

	return FreezeResult{Frozen: frozen, Skipped: skipped, File: outFile}, nil
}

// Load reads a frozen env file and returns the FrozenEnv.
func Load(path string) (FrozenEnv, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FrozenEnv{}, fmt.Errorf("freezer: read %s: %w", path, err)
	}
	var fe FrozenEnv
	if err := json.Unmarshal(data, &fe); err != nil {
		return FrozenEnv{}, fmt.Errorf("freezer: unmarshal: %w", err)
	}
	return fe, nil
}

// Thaw returns the env map from a FrozenEnv.
func Thaw(fe FrozenEnv) map[string]string {
	out := make(map[string]string, len(fe.Values))
	for k, v := range fe.Values {
		out[k] = v
	}
	return out
}
