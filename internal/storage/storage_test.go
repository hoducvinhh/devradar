package storage_test

import (
	"os"
	"testing"
	"time"

	"github.com/devradar/devradar/internal/storage"
	"github.com/devradar/devradar/pkg/models"
)

func TestStorage_ConfigAndProfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "devradar-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := storage.NewStorageWithDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// 1. Config Test
	cfg, err := store.LoadConfig()
	if err != nil {
		t.Fatalf("loading empty config: %v", err)
	}
	if cfg.GitHubToken != "" {
		t.Errorf("expected empty token initially")
	}

	cfg.GitHubToken = "ghp_mock_token_12345"
	if err := store.SaveConfig(cfg); err != nil {
		t.Fatalf("saving config: %v", err)
	}

	loadedCfg, err := store.LoadConfig()
	if err != nil {
		t.Fatalf("re-loading config: %v", err)
	}
	if loadedCfg.GitHubToken != "ghp_mock_token_12345" {
		t.Errorf("expected token match, got %s", loadedCfg.GitHubToken)
	}

	// 2. Profile Test
	prof, err := store.LoadProfile()
	if err != nil {
		t.Fatalf("loading default profile: %v", err)
	}
	if prof.Level != models.DifficultyBeginner {
		t.Errorf("expected beginner default level, got %s", prof.Level)
	}

	prof.Level = models.DifficultyIntermediate
	prof.Languages = []string{"go", "rust"}
	if err := store.SaveProfile(prof); err != nil {
		t.Fatalf("saving profile: %v", err)
	}

	loadedProf, err := store.LoadProfile()
	if err != nil {
		t.Fatalf("re-loading profile: %v", err)
	}
	if loadedProf.Level != models.DifficultyIntermediate || len(loadedProf.Languages) != 2 {
		t.Errorf("unexpected profile data: %+v", loadedProf)
	}

	// 3. Reset Profile
	if err := store.ResetProfile(); err != nil {
		t.Fatalf("resetting profile: %v", err)
	}
	resetProf, _ := store.LoadProfile()
	if len(resetProf.Languages) != 0 {
		t.Errorf("expected clean profile after reset, got %+v", resetProf)
	}
}

func TestStorage_Bookmarks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "devradar-bm-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := storage.NewStorageWithDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	bm1 := models.Bookmark{
		IssueID:      101,
		Number:       1,
		Title:        "First Issue",
		HTMLURL:      "https://github.com/foo/bar/issues/1",
		RepoFullName: "foo/bar",
	}

	if err := store.SaveBookmark(bm1); err != nil {
		t.Fatalf("saving bookmark: %v", err)
	}

	bms, err := store.LoadBookmarks()
	if err != nil {
		t.Fatalf("loading bookmarks: %v", err)
	}
	if len(bms) != 1 || bms[0].Number != 1 {
		t.Fatalf("unexpected bookmarks count: %d", len(bms))
	}

	// Remove bookmark
	removed, err := store.RemoveBookmark("1")
	if err != nil || !removed {
		t.Fatalf("expected bookmark to be removed, got removed=%v, err=%v", removed, err)
	}

	bmsAfter, _ := store.LoadBookmarks()
	if len(bmsAfter) != 0 {
		t.Errorf("expected empty bookmarks, got %d", len(bmsAfter))
	}
}

func TestCache_TTL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "devradar-cache-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	c := storage.NewCache(tmpDir, 50*time.Millisecond)

	repo := models.Repository{
		FullName: "test/repo",
		Language: "Go",
	}

	c.SetRepo(repo.FullName, repo)

	cached, found := c.GetRepo(repo.FullName)
	if !found || cached.Language != "Go" {
		t.Errorf("expected cached repo immediately")
	}

	// Wait for TTL expiration
	time.Sleep(70 * time.Millisecond)
	_, foundAfter := c.GetRepo(repo.FullName)
	if foundAfter {
		t.Errorf("expected item to be expired")
	}
}
