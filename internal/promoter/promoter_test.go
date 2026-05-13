package promoter_test

import (
	"testing"

	"github.com/user/envguard/internal/promoter"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{}
	_, sum := promoter.Promote(src, dst, promoter.Options{})
	if sum.Promoted != 2 {
		t.Fatalf("expected 2 promoted, got %d", sum.Promoted)
	}
	if dst["FOO"] != "bar" || dst["BAZ"] != "qux" {
		t.Fatal("unexpected dst values")
	}
}

func TestPromote_DenyList(t *testing.T) {
	src := map[string]string{"SECRET": "s3cr3t", "APP_NAME": "myapp"}
	dst := map[string]string{}
	_, sum := promoter.Promote(src, dst, promoter.Options{DenyKeys: []string{"SECRET"}})
	if sum.Promoted != 1 || sum.Skipped != 1 {
		t.Fatalf("unexpected counts: %+v", sum)
	}
	if _, ok := dst["SECRET"]; ok {
		t.Fatal("SECRET should not have been promoted")
	}
}

func TestPromote_AllowList(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	_, sum := promoter.Promote(src, dst, promoter.Options{AllowKeys: []string{"A", "C"}})
	if sum.Promoted != 2 {
		t.Fatalf("expected 2 promoted, got %d", sum.Promoted)
	}
	if _, ok := dst["B"]; ok {
		t.Fatal("B should not have been promoted")
	}
}

func TestPromote_NoOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}
	_, sum := promoter.Promote(src, dst, promoter.Options{Overwrite: false})
	if sum.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", sum.Skipped)
	}
	if dst["KEY"] != "old" {
		t.Fatal("existing value should not be overwritten")
	}
}

func TestPromote_WithOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}
	_, sum := promoter.Promote(src, dst, promoter.Options{Overwrite: true})
	if sum.Promoted != 1 {
		t.Fatalf("expected 1 promoted, got %d", sum.Promoted)
	}
	if dst["KEY"] != "new" {
		t.Fatal("value should have been overwritten")
	}
}

func TestPromote_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"X": "1"}
	dst := map[string]string{}
	promoter.Promote(src, dst, promoter.Options{})
	if len(src) != 1 {
		t.Fatal("src should not be mutated")
	}
}
