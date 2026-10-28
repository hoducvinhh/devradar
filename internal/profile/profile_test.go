package profile_test

import (
	"os"
	"testing"

	"github.com/devradar/devradar/internal/profile"
	"github.com/devradar/devradar/internal/storage"
)

func TestProfileManager(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "devradar-prof-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := storage.NewStorageWithDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	mgr := profile.NewManager(store)

	// Initial default
	p, err := mgr.Get()
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if len(p.Languages) != 0 {
		t.Errorf("expected empty languages initially")
	}

	// Update profile
	updated, err := mgr.Update("intermediate", []string{"Go", "Python", "go"}, []string{"Backend", "CLI"}, 100)
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}

	if updated.Level != "intermediate" {
		t.Errorf("expected level intermediate, got %s", updated.Level)
	}
	// Verify deduplication and lowercase
	if len(updated.Languages) != 2 {
		t.Errorf("expected 2 deduplicated languages, got %d", len(updated.Languages))
	}
	if updated.MinStars != 100 {
		t.Errorf("expected min stars 100, got %d", updated.MinStars)
	}

	// Test helper ParseCommaSeparated
	parsed := profile.ParseCommaSeparated("  go , rust, Rust , python ")
	if len(parsed) != 3 {
		t.Errorf("expected 3 deduplicated items, got %v", parsed)
	}

	// Reset
	if err := mgr.Reset(); err != nil {
		t.Fatalf("reset profile: %v", err)
	}
	resetP, _ := mgr.Get()
	if len(resetP.Languages) != 0 {
		t.Errorf("expected clean profile after reset")
	}
}
