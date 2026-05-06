package merger

import (
	"fmt"
	"sort"
)

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// StrategyBase keeps the base value on conflict.
	StrategyBase Strategy = iota
	// StrategyOverride replaces base value with override on conflict.
	StrategyOverride
)

// Result holds the merged environment and metadata about the operation.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
	Added     []string // keys present only in override
}

// Conflict records a key whose value differed between base and override.
type Conflict struct {
	Key           string
	BaseValue     string
	OverrideValue string
	Resolved      string
}

// Merge combines base and override env maps using the given strategy.
// Keys in override that are absent from base are always added.
func Merge(base, override map[string]string, strategy Strategy) Result {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []Conflict
	var added []string

	for k, ov := range override {
		bv, exists := merged[k]
		if !exists {
			merged[k] = ov
			added = append(added, k)
			continue
		}
		if bv != ov {
			resolved := bv
			if strategy == StrategyOverride {
				resolved = ov
				merged[k] = ov
			}
			conflicts = append(conflicts, Conflict{
				Key:           k,
				BaseValue:     bv,
				OverrideValue: ov,
				Resolved:      resolved,
			})
		}
	}

	sort.Strings(added)
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Merged: merged, Conflicts: conflicts, Added: added}
}

// StrategyFromString parses a strategy name ("base" or "override").
func StrategyFromString(s string) (Strategy, error) {
	switch s {
	case "base":
		return StrategyBase, nil
	case "override":
		return StrategyOverride, nil
	default:
		return StrategyBase, fmt.Errorf("unknown merge strategy %q: use \"base\" or \"override\"", s)
	}
}
