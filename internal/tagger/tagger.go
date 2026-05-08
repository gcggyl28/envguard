package tagger

import (
	"sort"
	"strings"
)

// Tag represents a label assigned to an env key.
type Tag struct {
	Key   string
	Value string
	Tags  []string
}

// Options controls tagging behaviour.
type Options struct {
	// AutoTag applies built-in heuristic tags (e.g. "secret", "url", "flag").
	AutoTag bool
	// CustomTags maps a tag name to a list of key substrings that trigger it.
	CustomTags map[string][]string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		AutoTag:    true,
		CustomTags: map[string][]string{},
	}
}

var secretHints = []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "APIKEY", "API_KEY", "PRIVATE"}
var urlHints = []string{"URL", "HOST", "ENDPOINT", "ADDR", "ADDRESS", "DSN"}
var flagHints = []string{"ENABLE", "DISABLE", "FLAG", "FEATURE", "TOGGLE"}

// Tag assigns tags to each key in env based on options.
func Tag(env map[string]string, opts Options) []Tag {
	results := make([]Tag, 0, len(env))

	for k, v := range env {
		upper := strings.ToUpper(k)
		tags := []string{}

		if opts.AutoTag {
			if matchesAny(upper, secretHints) {
				tags = appendUniq(tags, "secret")
			}
			if matchesAny(upper, urlHints) {
				tags = appendUniq(tags, "url")
			}
			if matchesAny(upper, flagHints) {
				tags = appendUniq(tags, "flag")
			}
		}

		for tag, hints := range opts.CustomTags {
			if matchesAny(upper, hints) {
				tags = appendUniq(tags, tag)
			}
		}

		sort.Strings(tags)
		results = append(results, Tag{Key: k, Value: v, Tags: tags})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

func matchesAny(upper string, hints []string) bool {
	for _, h := range hints {
		if strings.Contains(upper, h) {
			return true
		}
	}
	return false
}

func appendUniq(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return tags
		}
	}
	return append(tags, tag)
}
