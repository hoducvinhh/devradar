# DevRadar 📡

> **Find your next Open Source contribution.**  
> Smart CLI tool connecting developers with open-source GitHub issues tailored to their skills, experience level, and interests.

[![CI](https://github.com/devradar/devradar/actions/workflows/ci.yml/badge.svg)](https://github.com/devradar/devradar/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)

---

## ✨ Why DevRadar?

Finding suitable issues in open source can be overwhelming:
- Manually scrolling through thousands of repositories is exhausting.
- The `good first issue` label is often inconsistent across organizations.
- Sending pull requests to abandoned repositories leads to abandoned contributions.

**DevRadar** acts as your **Open Source Decision Assistant**:
1. 🎯 **Skill & Level Matching**: Scores issues (0–100) based on your real programming stack and target difficulty.
2. 🧠 **Difficulty Classification Engine**: Automatically estimates whether an issue is **Beginner**, **Intermediate**, or **Advanced**.
3. 🩺 **Repository Health Analysis**: Flags abandoned/archived repositories before you invest your time.
4. 💡 **Transparent Match Explanation**: Clearly explains *why* an issue was recommended.
5. 🔒 **Privacy-First & Offline Storage**: Zero telemetry, zero cloud lock-in. Profiles and bookmarks are stored safely on your machine (`~/.devradar/`).

---

## 🚀 Quick Start

### Installation

#### Pre-built Binaries
Download the latest binary for Linux, macOS, or Windows from [Releases](https://github.com/devradar/devradar/releases).

#### Build with Go
```bash
go install github.com/devradar/devradar/cmd/devradar@latest
```

Or clone and build locally:
```bash
git clone https://github.com/devradar/devradar.git
cd devradar
go build -o devradar ./cmd/devradar
```

---

## 🛠️ Usage Examples

### 1. Discover Issues
```bash
# Find beginner issues in Go
devradar find --language go --level beginner

# Find unassigned "good first issue" tasks
devradar find --label "good first issue" --unassigned --limit 15

# Search inside a specific organization or repo
devradar find --org kubernetes --level intermediate
devradar find --repo gin-gonic/gin --unassigned

# Output to JSON or CSV
devradar find --language python --format json
devradar find --language rust --format csv > issues.csv
```

### 2. Inspect Details & Match Rationale
```bash
devradar view https://github.com/gin-gonic/gin/issues/3937
```
```text
gin-gonic/gin #3937: Bind recursion struct param fatal exit
🔗 URL: https://github.com/gin-gonic/gin/issues/3937
🏷️ Labels: bug enhancement
📊 Difficulty: INTERMEDIATE (Standard issue complexity)

Match Breakdown & Score
Score: 84/100 [Good Match]
  • [30/30] Ngôn ngữ Go khớp chính xác với tiêu chí tìm kiếm
  • [13/25] Độ khó intermediate (chênh lệch 1 bậc so với mong muốn beginner)
  • [15/15] Cập nhật 1 ngày trước (Rất mới)
  • [14/15] Active (pushed 20 days ago)
  • [07/10] Mô tả chi tiết (có code block, trình bày có cấu trúc)
  • [05/05] Issue đang để trống assignee (Sẵn sàng tiếp nhận)
```

### 3. Personalize Your Profile
```bash
# Set your preferred difficulty and stack
devradar profile set --level beginner --languages go,python --interests backend,cli

# View current profile
devradar profile show
```

### 4. Bookmark & Track Issues
```bash
# Save an issue to work on later
devradar save https://github.com/gin-gonic/gin/issues/3937

# View all saved issues
devradar saved

# Remove from bookmarks
devradar unsave 3937
```

### 5. Repository Health Check
```bash
devradar repo gin-gonic/gin
```

### 6. Authentication (Increase Quota to 5,000 req/hr)
```bash
devradar auth set-token <your_github_token>
devradar auth status
```

---

## 🧮 Matching & Scoring Algorithm

DevRadar calculates a normalized score from **0 to 100**:

$$\text{Total Score} = (S_{\text{skill}} \times 0.30) + (S_{\text{difficulty}} \times 0.25) + (S_{\text{freshness}} \times 0.15) + (S_{\text{repo\_activity}} \times 0.15) + (S_{\text{clarity}} \times 0.10) + (S_{\text{unassigned}} \times 0.05)$$

- **Skill Match (30%)**: Matches issue language with user profile or search criteria.
- **Difficulty Match (25%)**: Aligns Beginner / Intermediate / Advanced scope with user preference.
- **Freshness (15%)**: Evaluates recent activity on the issue thread.
- **Repository Health (15%)**: Ensures the repo is active and maintainers accept contributions.
- **Issue Clarity (10%)**: Checks for structured descriptions, checklists, and code snippets.
- **Unassigned Status (5%)**: Prioritizes open issues not yet assigned to other developers.

---

## 📁 Project Structure

```text
devradar/
├── cmd/devradar/main.go        # Application entrypoint
├── internal/
│   ├── app/                    # Dependency injection & context
│   ├── cli/                    # Cobra CLI commands (find, view, profile, etc.)
│   ├── github/                 # GitHub REST API client & rate limiting
│   ├── discovery/              # Search orchestrator
│   ├── classification/         # Difficulty classifier
│   ├── ranking/                # Multi-factor matching engine
│   ├── repository/             # Repository health analyzer
│   ├── profile/                # User profile manager
│   ├── storage/                # Local storage engine & TTL cache
│   └── output/                 # Lip Gloss table, JSON, and CSV formatters
├── pkg/models/                 # Domain structs
├── docs/                       # Architecture & usage documentation
├── .github/workflows/          # GitHub Actions CI & GoReleaser
└── Makefile                    # Build & test shortcuts
```

---

## 🧪 Testing

Run all unit tests with code coverage:
```bash
go test -v -cover ./...
```

---

## 🤝 Contributing

Contributions are warmly welcome! Please check [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for guidelines.

## 📄 License

DevRadar is licensed under the [MIT License](LICENSE).
