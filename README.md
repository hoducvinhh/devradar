# DevRadar 📡

> **Find your next Open Source contribution.**  
> Smart CLI tool connecting developers with open-source GitHub issues tailored to their skills, experience level, and interests.

[![CI](https://github.com/hoducvinhh/devradar/actions/workflows/ci.yml/badge.svg)](https://github.com/hoducvinhh/devradar/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

---

## 📑 Table of Contents

- [Why DevRadar?](#-why-devradar)
- [Key Features](#-key-features)
- [Installation](#-installation)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Install via Go](#install-via-go)
  - [Build from Source](#build-from-source)
- [Authentication & Setup](#-authentication--setup)
- [Complete Command Reference](#-complete-command-reference)
  - [1. `devradar find` – Discover & Rank Issues](#1-devradar-find--discover--rank-issues)
  - [2. `devradar view` – Inspect Issue Details & Rationale](#2-devradar-view--inspect-issue-details--rationale)
  - [3. `devradar profile` – Manage Developer Preferences](#3-devradar-profile--manage-developer-preferences)
  - [4. `devradar repo` – Repository Health & Activity Audit](#4-devradar-repo--repository-health--activity-audit)
  - [5. `devradar save` / `saved` / `unsave` – Local Bookmarking](#5-devradar-save--saved--unsave--local-bookmarking)
- [Scoring & Classification Algorithm](#-scoring--classification-algorithm)
  - [Difficulty Classification Heuristics](#difficulty-classification-heuristics)
  - [Scoring Formula](#scoring-formula)
- [Local Storage & Privacy](#-local-storage--privacy)
- [Project Architecture](#-project-architecture)
- [Development & Testing](#-development--testing)
- [Contributing](#-contributing)
- [License](#-license)

---

## 💡 Why DevRadar?

Finding the right open-source issue to contribute to is often frustrating:
- **Noise & Volume**: Manually browsing GitHub Trending or Explore returns thousands of issues without context.
- **Inconsistent Difficulty Labels**: The label `good first issue` can range from fixing a typo to rewriting an asynchronous engine.
- **Ghost / Abandoned Projects**: Spending days crafting a Pull Request only to discover the repository was abandoned months ago.
- **First-Timer Friction**: Lack of clarity on issue requirements, scope, or whether an issue is already assigned.

**DevRadar** is an intelligent **Open Source Decision Assistant** running right inside your terminal. It evaluates real issue complexity, filters out unmaintained repositories, and scores opportunities based on your personal skill stack.

---

## ✨ Key Features

- 🎯 **Multi-Factor Ranking Engine**: Computes a normalized score (0–100) combining language match, difficulty alignment, freshness, repository health, clarity, and assignee availability.
- 🧠 **Rule-Based Difficulty Classifier**: Categorizes issues into `BEGINNER`, `INTERMEDIATE`, and `ADVANCED` using title, body heuristics, label semantics, and code snippets.
- 🩺 **Repository Health Diagnostics**: Inspects push frequency, activity recency, and flags archived/abandoned repositories before you write a single line of code.
- 🔍 **Flexible Discovery**: Search across GitHub or target specific organizations (`--org kubernetes`), repositories (`--repo gin-gonic/gin`), tags, or minimum stars.
- 📊 **Multi-Format Export**: Pretty-printed Lip Gloss terminal tables, structured `json` for automation scripts, or standard `csv` for spreadsheet analysis.
- 🔖 **Built-in Bookmark Workflow**: Save promising tasks locally and manage them with simple terminal commands.
- 🔒 **100% Privacy & Offline-First**: Your profile and bookmarks stay locally in `~/.devradar/`. No telemetry, no external accounts, no cloud sync.

---

## 🚀 Installation

### Pre-built Binaries

Download pre-compiled binaries for **Linux**, **macOS**, and **Windows** from the [GitHub Releases](https://github.com/hoducvinhh/devradar/releases) page. Extract the archive and place `devradar` in your system `PATH`.

### Install via Go

If you have Go 1.22 or higher installed:

```bash
go install github.com/hoducvinhh/devradar/cmd/devradar@latest
```

### Build from Source

```bash
# 1. Clone the repository
git clone https://github.com/hoducvinhh/devradar.git
cd devradar

# 2. Build executable
go build -o devradar ./cmd/devradar

# Optional: Install locally using Makefile
make build
make install
```

---

## 🔐 Authentication & Setup

By default, DevRadar uses the unauthenticated GitHub REST API, which provides **60 requests per hour**. To unlock the full quota of **5,000 requests per hour**, configure a GitHub Personal Access Token (PAT).

### 1. Generate a GitHub Token
1. Go to [GitHub Settings → Developer Settings → Personal Access Tokens (Classic or Fine-Grained)](https://github.com/settings/tokens).
2. For Classic tokens, choose the **`public_repo`** scope (or no scopes if you only read public issues).
3. Copy your generated token (`ghp_...`).

### 2. Configure Token in DevRadar

```bash
# Save your token locally (stored securely in ~/.devradar/config.json)
devradar auth set-token ghp_your_personal_access_token_here

# Verify your authentication status and remaining quota
devradar auth status

# If you ever want to remove the saved token:
devradar auth clear
```

Alternatively, you can provide the token via environment variable or CLI flag:
```bash
# Environment variable
export GITHUB_TOKEN="ghp_your_token_here"
devradar find --language go

# Global CLI flag
devradar find --language go --token ghp_your_token_here
```

---

## 🛠️ Complete Command Reference

Global flags available for all commands:
- `--format, -f`: Output format (`table`, `json`, `csv`). Default: `table`.
- `--token`: Override GitHub Personal Access Token for the current command.
- `--verbose, -v`: Enable verbose logging.
- `--help, -h`: Show command help.

---

### 1. `devradar find` – Discover & Rank Issues

Searches open GitHub issues, evaluates repository health and issue difficulty, and ranks them according to your preferences.

```bash
devradar find [flags]
```

#### Flags

| Flag | Shorthand | Type | Default | Description |
| :--- | :---: | :---: | :---: | :--- |
| `--language` | `-l` | `string` | `""` | Programming language filter (e.g. `go`, `python`, `rust`, `typescript`) |
| `--label` | | `string` | `""` | GitHub label filter (e.g. `"good first issue"`, `"help wanted"`, `"bug"`) |
| `--level` | | `string` | `""` | Difficulty level: `beginner`, `intermediate`, `advanced` |
| `--org` | | `string` | `""` | Search strictly within a GitHub organization (e.g. `kubernetes`, `golang`) |
| `--repo` | | `string` | `""` | Search inside a specific repository (e.g. `gin-gonic/gin`, `cli/cli`) |
| `--unassigned` | `-u` | `bool` | `false` | Only return open issues that currently have no assignee |
| `--limit` | | `int` | `15` | Maximum number of ranked issues to return |
| `--min-stars` | | `int` | `0` | Minimum repository star count |
| `--sort` | | `string` | `"updated"` | Sort criteria: `updated`, `created`, `comments` |
| `--order` | | `string` | `"desc"` | Sort direction: `desc`, `asc` |
| `--format` | `-f` | `string` | `"table"` | Output format: `table`, `json`, `csv` |

#### Examples

```bash
# Find beginner-friendly issues in Go with at least 500 stars
devradar find --language go --level beginner --min-stars 500

# Search for unassigned 'good first issue' tasks in Kubernetes org
devradar find --org kubernetes --label "good first issue" --unassigned --limit 10

# Search issues in a specific repository
devradar find --repo gin-gonic/gin --level intermediate

# Export issues to CSV for external spreadsheet review
devradar find --language rust --level beginner --format csv > rust_beginner.csv

# Pipe JSON output to jq for automated scripts
devradar find --language python --format json | jq '.[] | {title: .title, score: .score, url: .html_url}'
```

---

### 2. `devradar view` – Inspect Issue Details & Rationale

Displays the full context of an issue, including its complete description, labels, difficulty rationale, repository health metrics, and a factor-by-factor breakdown of why it received its score.

```bash
# By full URL
devradar view https://github.com/gin-gonic/gin/issues/3937

# Or by shorthand repository#number
devradar view gin-gonic/gin#3937
```

#### Example Output

```text
gin-gonic/gin #3937: Bind recursion struct param fatal exit
🔗 URL: https://github.com/gin-gonic/gin/issues/3937
🏷️ Labels: bug, enhancement
📊 Difficulty: INTERMEDIATE (Standard issue complexity)

Match Breakdown & Score
Score: 84/100 [Good Match]
  • [30/30] Ngôn ngữ Go khớp chính xác với tiêu chí tìm kiếm
  • [13/25] Độ khó intermediate (chênh lệch 1 bậc so với mong muốn beginner)
  • [15/15] Cập nhật 1 ngày trước (Rất mới)
  • [14/15] Active (pushed 20 days ago)
  • [07/10] Mô tả chi tiết (có code block, trình bày có cấu trúc)
  • [05/05] Issue đang để trống assignee (Sẵn sàng tiếp nhận)

💡 Contribution Tip:
Check CONTRIBUTING.md, reproduce the issue with a minimal test case, and discuss your approach in the thread before submitting a PR.
```

---

### 3. `devradar profile` – Manage Developer Preferences

Save your preferred languages, target difficulty level, and topics to your local machine (`~/.devradar/profile.json`). When set, `devradar find` automatically uses your profile to personalize issue scores without needing manual flags every time.

```bash
# View your current profile configuration
devradar profile show

# Update your preferences
devradar profile set --level beginner --languages go,python,rust --interests cli,backend,k8s --min-stars 100

# Reset profile back to defaults
devradar profile reset
```

#### Profile Flags for `devradar profile set`:
- `--level`: Target difficulty (`beginner`, `intermediate`, `advanced`)
- `--languages`: Comma-separated list of preferred languages (e.g. `go,python,typescript`)
- `--interests`: Comma-separated domain topics (e.g. `cli,docker,devtools`)
- `--min-stars`: Default minimum star threshold

---

### 4. `devradar repo` – Repository Health & Activity Audit

Before investing effort into a contribution, verify if the repository is healthy, actively reviewed, and friendly to external contributors.

```bash
devradar repo gin-gonic/gin
devradar repo kubernetes/kubernetes
```

#### What It Checks:
- **Health Score (0–100)**: Evaluates maintenance recency and responsiveness.
- **Commit / Push Recency**: Detects stale projects (e.g. inactive > 6 months).
- **Archival / Read-Only Warning**: Warns you immediately if a repository is archived so you don't waste time preparing a PR.
- **Community Scale**: Reports stars, forks, and open issues.

---

### 5. `devradar save` / `saved` / `unsave` – Local Bookmarking

Manage tasks you plan to work on without leaving your command line.

```bash
# 1. Bookmark an issue
devradar save https://github.com/gin-gonic/gin/issues/3937
# Or with shorthand
devradar save cli/cli#8420

# 2. View all saved bookmarks in a neat table
devradar saved

# 3. Export bookmarks as JSON
devradar saved -f json

# 4. Remove an issue from bookmarks when done
devradar unsave 3937
# Or by full URL
devradar unsave https://github.com/gin-gonic/gin/issues/3937
```

---

## 🧮 Scoring & Classification Algorithm

### Difficulty Classification Heuristics

DevRadar's classification engine inspects issue metadata using rule-based scoring:

| Level | Key Indicators |
| :--- | :--- |
| **BEGINNER** | Labels: `good first issue`, `beginner-friendly`, `easy`, `documentation`, `starter`.<br>Content: Typos, docs fixes, lint cleanups, small configuration adjustments. |
| **INTERMEDIATE** | Labels: `enhancement`, `feature`, `help wanted`, `bug`.<br>Content: Standard bug fixes, adding new CLI flags, unit tests, refactoring module logic. |
| **ADVANCED** | Labels: `architecture`, `performance`, `security`, `breaking-change`, `critical`.<br>Content: Memory leaks, concurrency/deadlocks, distributed systems, protocol changes. |

---

### Scoring Formula

Each issue receives a composite recommendation score from **0 to 100**:

$$\text{Total Score} = (S_{\text{skill}} \times 0.30) + (S_{\text{difficulty}} \times 0.25) + (S_{\text{freshness}} \times 0.15) + (S_{\text{repo\_activity}} \times 0.15) + (S_{\text{clarity}} \times 0.10) + (S_{\text{unassigned}} \times 0.05)$$

1. **Skill Match ($S_{\text{skill}}$, 30%)**:
   - Exact match with your target languages: **100 pts**
   - No language filter specified: **70 pts**
   - Mismatch: **0 pts**
2. **Difficulty Alignment ($S_{\text{difficulty}}$, 25%)**:
   - Exact match with target difficulty: **100 pts**
   - 1 step difference (e.g. beginner searching intermediate): **50 pts**
   - 2 steps difference: **20 pts**
3. **Freshness ($S_{\text{freshness}}$, 15%)**:
   - Updated within 3 days: **100 pts**
   - Updated within 14 days: **80 pts**
   - Updated within 30 days: **50 pts**
   - Older: **20 pts**
4. **Repository Activity ($S_{\text{repo\_activity}}$, 15%)**:
   - Based on repository health score (pushed recently, active commit history, non-archived).
5. **Issue Clarity ($S_{\text{clarity}}$, 10%)**:
   - Rewards issues containing structured descriptions, reproduction steps, checklists, and code snippets.
6. **Unassigned Availability ($S_{\text{unassigned}}$, 5%)**:
   - Open issue with no current assignee: **100 pts**
   - Already assigned: **0 pts**

---

## 💾 Local Storage & Privacy

DevRadar keeps all state strictly on your local machine under `~/.devradar/`:

```text
~/.devradar/
├── config.json       # GitHub API credentials & CLI settings
├── profile.json      # Your skill profile, target level, and preferences
├── bookmarks.json    # Saved issues and tasks
└── cache/            # Local TTL cache to minimize GitHub API rate consumption
```

- **Zero External Telemetry**: DevRadar never collects tracking data or analytics.
- **Direct GitHub Communication**: The CLI only connects to `api.github.com`.

---

## 📁 Project Architecture

```text
devradar/
├── cmd/
│   └── devradar/             # CLI application entrypoint (main.go)
├── internal/
│   ├── app/                  # Application container & dependency injection
│   ├── cli/                  # Cobra commands (find, view, profile, repo, bookmark, auth)
│   ├── github/               # GitHub REST API v3 client, rate limit tracker & HTTP client
│   ├── discovery/            # Issue discovery & search orchestrator
│   ├── classification/       # Difficulty classifier & heuristic rules
│   ├── ranking/              # Multi-factor weighting & match rationale generator
│   ├── repository/           # Repository activity & health analyzer
│   ├── profile/              # Developer profile state manager
│   ├── storage/              # JSON-based disk persistence & TTL cache
│   └── output/               # Lip Gloss terminal tables, CSV & JSON formatters
├── pkg/
│   └── models/               # Domain data structures (Issue, Repository, Profile, Bookmark)
├── docs/                     # Architecture documentation and guides
├── .github/                  # GitHub Actions CI workflow & GoReleaser config
└── Makefile                  # Build, test, lint, and run tasks
```

---

## 🧪 Development & Testing

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- Make (optional, for automation commands)

### Running Unit Tests
```bash
# Run all tests with coverage
go test -v -cover ./...

# Or using Makefile
make test
```

### Build & Run Locally
```bash
make build
./devradar find --language go --level beginner
```

---

## 🤝 Contributing

We welcome contributions of all levels! Whether you want to improve classification heuristics, add new export formats, or improve documentation:

1. Fork the repository.
2. Create your feature branch (`git checkout -b feat/amazing-feature`).
3. Commit your changes (`git commit -m 'feat: add amazing feature'`).
4. Push to the branch (`git push origin feat/amazing-feature`).
5. Open a Pull Request.

Please read [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for detailed guidelines.

---

## 📄 License

This project is licensed under the **MIT License** – see the [LICENSE](LICENSE) file for details.
