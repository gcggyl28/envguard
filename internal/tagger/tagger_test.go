package tagger

import (
	"testing"
)

func TestTag_NoHints(t *testing.T) {
	env := map[string]string{"APP_NAME": "envguard", "VERSION": "1.0"}
	results := Tag(env, DefaultOptions())
	for _, r := range results {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags for %s, got %v", r.Key, r.Tags)
		}
	}
}

func TestTag_SecretKey(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	results := Tag(env, DefaultOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !contains(results[0].Tags, "secret") {
		t.Errorf("expected 'secret' tag, got %v", results[0].Tags)
	}
}

func TestTag_URLKey(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://localhost/db"}
	results := Tag(env, DefaultOptions())
	if !contains(results[0].Tags, "url") {
		t.Errorf("expected 'url' tag, got %v", results[0].Tags)
	}
}

func TestTag_FlagKey(t *testing.T) {
	env := map[string]string{"ENABLE_METRICS": "true"}
	results := Tag(env, DefaultOptions())
	if !contains(results[0].Tags, "flag") {
		t.Errorf("expected 'flag' tag, got %v", results[0].Tags)
	}
}

func TestTag_AutoTagDisabled(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t", "API_TOKEN": "abc"}
	opts := Options{AutoTag: false, CustomTags: map[string][]string{}}
	results := Tag(env, opts)
	for _, r := range results {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags when AutoTag=false, got %v for %s", r.Tags, r.Key)
		}
	}
}

func TestTag_CustomTag(t *testing.T) {
	env := map[string]string{"STRIPE_KEY": "sk_live_xxx", "UNRELATED": "val"}
	opts := Options{
		AutoTag: false,
		CustomTags: map[string][]string{
			"payment": {"STRIPE", "PAYPAL"},
		},
	}
	results := Tag(env, opts)
	for _, r := range results {
		if r.Key == "STRIPE_KEY" && !contains(r.Tags, "payment") {
			t.Errorf("expected 'payment' tag for STRIPE_KEY, got %v", r.Tags)
		}
		if r.Key == "UNRELATED" && len(r.Tags) != 0 {
			t.Errorf("expected no tags for UNRELATED, got %v", r.Tags)
		}
	}
}

func TestTag_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	results := Tag(env, DefaultOptions())
	if results[0].Key != "A_KEY" || results[1].Key != "M_KEY" || results[2].Key != "Z_KEY" {
		t.Errorf("expected sorted output, got %v %v %v", results[0].Key, results[1].Key, results[2].Key)
	}
}

func contains(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
