# DevRadar CLI Usage Guide

## Getting Started

### 1. Installation

#### Download Binary
Download pre-built binaries from GitHub Releases for Linux, macOS, or Windows.

#### Build from Source
```bash
git clone https://github.com/devradar/devradar.git
cd devradar
go build -o devradar ./cmd/devradar
```

---

## Command Reference

### `devradar find`
Search and rank open-source issues across GitHub.

```bash
# Find beginner issues in Go
devradar find --language go --level beginner

# Find unassigned issues with label "good first issue"
devradar find --label "good first issue" --unassigned --limit 15

# Search inside a specific organization or repo
devradar find --org kubernetes --level intermediate
devradar find --repo gin-gonic/gin --unassigned

# Output formats
devradar find --language rust --format json
devradar find --language python --format csv > results.csv
```

### `devradar view`
Inspect an issue in detail, review match explanations, and read contribution tips.

```bash
devradar view https://github.com/gin-gonic/gin/issues/3937
devradar view gin-gonic/gin#3937
```

### `devradar profile`
Configure your developer profile locally to receive personalized rankings.

```bash
# Set preferences
devradar profile set --level beginner --languages go,python --interests backend,cli

# Show preferences
devradar profile show

# Reset preferences
devradar profile reset
```

### `devradar repo`
Analyze repository maintenance activity and contributor friendliness.

```bash
devradar repo gin-gonic/gin
devradar repo golang/go
```

### `devradar save`, `saved`, `unsave`
Track interesting issues without leaving your terminal.

```bash
# Save an issue
devradar save https://github.com/gin-gonic/gin/issues/3937

# List bookmarked issues
devradar saved

# Remove bookmark
devradar unsave 3937
```

### `devradar auth`
Manage your GitHub API token to expand your rate limit to 5,000 requests/hour.

```bash
# Save token
devradar auth set-token ghp_your_personal_access_token

# Check quota
devradar auth status

# Clear token
devradar auth clear
```
