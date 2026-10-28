package discovery

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/devradar/devradar/internal/classification"
	"github.com/devradar/devradar/internal/github"
	"github.com/devradar/devradar/internal/ranking"
	"github.com/devradar/devradar/internal/repository"
	"github.com/devradar/devradar/internal/storage"
	"github.com/devradar/devradar/pkg/models"
)

// Service coordinates issue discovery, classification, and ranking.
type Service struct {
	client     *github.Client
	storage    *storage.Storage
	cache      *storage.Cache
	classifier *classification.Classifier
	analyzer   *repository.Analyzer
	ranker     *ranking.Ranker
}

// NewService creates a new discovery service instance.
func NewService(
	client *github.Client,
	store *storage.Storage,
	cache *storage.Cache,
) *Service {
	return &Service{
		client:     client,
		storage:    store,
		cache:      cache,
		classifier: classification.NewClassifier(),
		analyzer:   repository.NewAnalyzer(),
		ranker:     ranking.NewRanker(),
	}
}

// Discover finds and ranks issues according to filters and user profile.
func (s *Service) Discover(ctx context.Context, opts *models.FilterOptions) ([]models.Issue, error) {
	if opts == nil {
		opts = models.DefaultFilterOptions()
	}

	// 1. Load user profile for personalized matching
	var prof *models.Profile
	if s.storage != nil {
		p, err := s.storage.LoadProfile()
		if err == nil {
			prof = p
		}
	}
	if prof == nil {
		prof = models.DefaultProfile()
	}

	// Apply profile language if none specified in filter
	effectiveLang := opts.Language
	if effectiveLang == "" && len(prof.Languages) > 0 {
		effectiveLang = prof.Languages[0]
	}

	// Apply profile difficulty if none specified in filter
	effectiveLevel := opts.Level
	if effectiveLevel == "" && prof.Level != "" {
		effectiveLevel = prof.Level
	}

	// 2. Fetch raw issues from GitHub REST API
	searchOpts := *opts
	if searchOpts.Language == "" && effectiveLang != "" {
		searchOpts.Language = effectiveLang
	}

	issues, err := s.client.SearchIssues(ctx, &searchOpts)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return []models.Issue{}, nil
	}

	// 3. Classify and enrich issues
	// Distinct repository names to fetch metadata efficiently
	uniqueRepos := make(map[string]bool)
	for _, issue := range issues {
		if issue.RepoFullName != "" {
			uniqueRepos[issue.RepoFullName] = true
		}
	}

	repoData := s.fetchReposMetadata(ctx, uniqueRepos)

	// 4. Score and Rank each issue
	for i := range issues {
		iss := &issues[i]

		// Difficulty Classification
		diff, diffReason := s.classifier.Classify(iss)
		iss.Difficulty = diff
		iss.DifficultyReason = diffReason

		// Repository Health
		var repoHealth float64 = 70.0
		var repoHealthReason string = "Repository metadata unavailable"

		if repo, exists := repoData[iss.RepoFullName]; exists && repo != nil {
			iss.RepoStars = repo.StargazersCount
			iss.RepoForks = repo.ForksCount
			iss.RepoPushedAt = repo.PushedAt
			iss.RepoArchived = repo.Archived
			if iss.RepoLanguage == "" {
				iss.RepoLanguage = repo.Language
			}

			repoHealth, repoHealthReason = s.analyzer.Analyze(repo)
		}

		// Multi-factor Ranking
		s.ranker.Rank(iss, prof, effectiveLang, effectiveLevel, repoHealth, repoHealthReason)
	}

	// 5. Filter by level if explicitly requested by user
	if opts.Level != "" && opts.Level != models.DifficultyUnknown {
		var filtered []models.Issue
		for _, iss := range issues {
			if strings.EqualFold(string(iss.Difficulty), string(opts.Level)) {
				filtered = append(filtered, iss)
			}
		}
		// If we had issues but filter reduced to 0, keep all so user sees alternatives with lower rank
		if len(filtered) > 0 {
			issues = filtered
		}
	}

	// 6. Sort results by total ranking score descending
	sort.Slice(issues, func(i, j int) bool {
		scoreI := 0
		scoreJ := 0
		if issues[i].Ranking != nil {
			scoreI = issues[i].Ranking.TotalScore
		}
		if issues[j].Ranking != nil {
			scoreJ = issues[j].Ranking.TotalScore
		}
		return scoreI > scoreJ
	})

	return issues, nil
}

// fetchReposMetadata fetches repo metadata with caching and concurrency control.
func (s *Service) fetchReposMetadata(ctx context.Context, repos map[string]bool) map[string]*models.Repository {
	result := make(map[string]*models.Repository)
	var mu sync.Mutex

	// Rate-limit safety: fetch up to 5 repos concurrently
	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for repoFullName := range repos {
		parts := strings.Split(repoFullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, name := parts[0], parts[1]

		// Check cache first
		if s.cache != nil {
			if cached, ok := s.cache.GetRepo(repoFullName); ok {
				result[repoFullName] = cached
				continue
			}
		}

		// Fetch from GitHub
		wg.Add(1)
		go func(o, n, full string) {
			defer wg.Done()
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			repo, err := s.client.GetRepo(ctx, o, n)
			if err == nil && repo != nil {
				mu.Lock()
				result[full] = repo
				mu.Unlock()

				if s.cache != nil {
					s.cache.SetRepo(full, *repo)
				}
			}
		}(owner, name, repoFullName)
	}

	// Wait with a short timeout to prevent blocking search if network stalls
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		// Proceed with whatever repo data was fetched
	}

	return result
}
