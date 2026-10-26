package repository_test

import (
	"testing"
	"time"

	"github.com/devradar/devradar/internal/repository"
	"github.com/devradar/devradar/pkg/models"
)

func TestAnalyzer_AllBranches(t *testing.T) {
	a := repository.NewAnalyzer()

	// 1. Disabled
	disabled := &models.Repository{FullName: "foo/bar", Disabled: true}
	s, _ := a.Analyze(disabled)
	if s != 0.0 {
		t.Errorf("expected 0 for disabled repo")
	}

	// 2. Archived
	arch := &models.Repository{FullName: "foo/bar", Archived: true}
	s, _ = a.Analyze(arch)
	if s != 0.0 {
		t.Errorf("expected 0 for archived repo")
	}

	// 3. Pushed within 14 days
	r14 := &models.Repository{PushedAt: time.Now().Add(-5 * 24 * time.Hour)}
	s, _ = a.Analyze(r14)
	if s != 100.0 {
		t.Errorf("expected 100, got %f", s)
	}

	// 4. Pushed within 30 days
	r30 := &models.Repository{PushedAt: time.Now().Add(-20 * 24 * time.Hour)}
	s, _ = a.Analyze(r30)
	if s != 90.0 {
		t.Errorf("expected 90, got %f", s)
	}

	// 5. Pushed within 60 days
	r60 := &models.Repository{PushedAt: time.Now().Add(-45 * 24 * time.Hour)}
	s, _ = a.Analyze(r60)
	if s != 75.0 {
		t.Errorf("expected 75, got %f", s)
	}

	// 6. Pushed within 90 days
	r90 := &models.Repository{PushedAt: time.Now().Add(-80 * 24 * time.Hour)}
	s, _ = a.Analyze(r90)
	if s != 60.0 {
		t.Errorf("expected 60, got %f", s)
	}

	// 7. Pushed within 180 days
	r180 := &models.Repository{PushedAt: time.Now().Add(-150 * 24 * time.Hour)}
	s, _ = a.Analyze(r180)
	if s != 40.0 {
		t.Errorf("expected 40, got %f", s)
	}

	// 8. Stale (>180 days)
	rstale := &models.Repository{PushedAt: time.Now().Add(-250 * 24 * time.Hour)}
	s, _ = a.Analyze(rstale)
	if s != 10.0 {
		t.Errorf("expected 10, got %f", s)
	}

	// 9. Stars bonus (>500 stars)
	rStars := &models.Repository{
		PushedAt:        time.Now().Add(-20 * 24 * time.Hour),
		StargazersCount: 600,
	}
	s, _ = a.Analyze(rStars)
	if s != 95.0 {
		t.Errorf("expected 95 with star bonus, got %f", s)
	}

	// 10. Nil
	snil, _ := a.Analyze(nil)
	if snil != 50.0 {
		t.Errorf("expected 50 for nil")
	}
}
