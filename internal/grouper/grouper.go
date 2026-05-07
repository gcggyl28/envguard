package grouper

import (
	"sort"
	"strings"
)

// Group represents a named collection of env keys.
type Group struct {
	Name string
	Keys []string
}

// Options controls how grouping is performed.
type Options struct {
	// ByPrefix groups keys by their prefix (text before the first underscore).
	ByPrefix bool
	// CustomGroups maps group name -> list of key prefixes or exact keys.
	CustomGroups map[string][]string
}

// GroupResult holds the output of a grouping operation.
type GroupResult struct {
	Groups   []Group
	Ungrouped []string
}

// Group partitions the provided env map according to Options.
func Group(env map[string]string, opts Options) GroupResult {
	if opts.CustomGroups != nil && len(opts.CustomGroups) > 0 {
		return groupByCustom(env, opts.CustomGroups)
	}
	if opts.ByPrefix {
		return groupByPrefix(env)
	}
	// Default: single group containing all keys alphabetically.
	keys := sortedKeys(env)
	return GroupResult{
		Groups: []Group{{Name: "all", Keys: keys}},
	}
}

func groupByPrefix(env map[string]string) GroupResult {
	prefixMap := map[string][]string{}
	var ungrouped []string

	for k := range env {
		if idx := strings.Index(k, "_"); idx > 0 {
			prefix := k[:idx]
			prefixMap[prefix] = append(prefixMap[prefix], k)
		} else {
			ungrouped = append(ungrouped, k)
		}
	}

	var groups []Group
	for name, keys := range prefixMap {
		sort.Strings(keys)
		groups = append(groups, Group{Name: name, Keys: keys})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	sort.Strings(ungrouped)
	return GroupResult{Groups: groups, Ungrouped: ungrouped}
}

func groupByCustom(env map[string]string, customGroups map[string][]string) GroupResult {
	assigned := map[string]bool{}
	var groups []Group

	groupNames := make([]string, 0, len(customGroups))
	for name := range customGroups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, name := range groupNames {
		patterns := customGroups[name]
		var matched []string
		for k := range env {
			for _, p := range patterns {
				if k == p || strings.HasPrefix(k, p+"_") {
					matched = append(matched, k)
					assigned[k] = true
					break
				}
			}
		}
		sort.Strings(matched)
		groups = append(groups, Group{Name: name, Keys: matched})
	}

	var ungrouped []string
	for k := range env {
		if !assigned[k] {
			ungrouped = append(ungrouped, k)
		}
	}
	sort.Strings(ungrouped)
	return GroupResult{Groups: groups, Ungrouped: ungrouped}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
