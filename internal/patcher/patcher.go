package patcher

import (
	"fmt"
	"strings"
)

// Op represents a single patch operation.
type Op struct {
	Key    string
	Value  string
	Delete bool
}

// Result holds the outcome of a patch operation.
type Result struct {
	Applied []string
	Deleted []string
	Skipped []string
}

// Patch applies a list of Ops to the given env map and returns a new map
// along with a Result describing what changed.
// If an Op has Delete=true the key is removed; otherwise it is set/overwritten.
// Keys in Skipped were requested for deletion but did not exist.
func Patch(env map[string]string, ops []Op) (map[string]string, Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var result Result
	for _, op := range ops {
		key := strings.TrimSpace(op.Key)
		if key == "" {
			continue
		}
		if op.Delete {
			if _, exists := out[key]; exists {
				delete(out, key)
				result.Deleted = append(result.Deleted, key)
			} else {
				result.Skipped = append(result.Skipped, key)
			}
		} else {
			out[key] = op.Value
			result.Applied = append(result.Applied, key)
		}
	}
	return out, result
}

// ParseOps parses a slice of "KEY=VALUE" or "!KEY" strings into Ops.
// Lines prefixed with '!' indicate a delete operation.
func ParseOps(lines []string) ([]Op, error) {
	ops := make([]Op, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			key := strings.TrimSpace(line[1:])
			if key == "" {
				return nil, fmt.Errorf("empty key in delete op: %q", line)
			}
			ops = append(ops, Op{Key: key, Delete: true})
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid patch op (expected KEY=VALUE or !KEY): %q", line)
		}
		ops = append(ops, Op{Key: strings.TrimSpace(parts[0]), Value: parts[1]})
	}
	return ops, nil
}
