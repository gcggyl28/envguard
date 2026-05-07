package filter

import (
	"regexp"
	"strings"
)

// Options controls which keys to include or exclude.
type Options struct {
	Prefix    string   // only keys with this prefix
	Suffix    string   // only keys with this suffix
	Contains  string   // only keys containing this substring
	Pattern   string   // only keys matching this regex
	Exclude   []string // explicit keys to exclude
}

// Filter returns a new map containing only the entries that satisfy opts.
func Filter(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	excluded := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excluded[k] = true
	}

	out := make(map[string]string)
	for k, v := range env {
		if excluded[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, opts.Suffix) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(k, opts.Contains) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		out[k] = v
	}
	return out, nil
}

// Summary describes the result of a filter operation.
type Summary struct {
	Total    int
	Included int
	Excluded int
}

// Summarize returns a Summary for the given before/after maps.
func Summarize(before, after map[string]string) Summary {
	return Summary{
		Total:    len(before),
		Included: len(after),
		Excluded: len(before) - len(after),
	}
}
