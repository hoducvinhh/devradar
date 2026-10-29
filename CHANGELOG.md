# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-MVP] - 2026-10-30

### Added
- **Global Issue Discovery**: Search across GitHub issues using `--language`, `--label`, `--level`, `--unassigned`, `--org`, and `--repo`.
- **Difficulty Classification Engine**: Rule-based estimator for Beginner, Intermediate, and Advanced issue scopes.
- **Multi-Factor Matching Engine**: 6-dimension scoring algorithm (0–100) with detailed match explanations.
- **Repository Health Checker**: Analyze maintenance frequency, commit freshness, star metrics, and read-only/archived status (`devradar repo`).
- **Local Profile Management**: Safe offline-first profile storage in `~/.devradar/profile.json` (`devradar profile`).
- **Issue Bookmarking & Tracking**: Local task tracker via `devradar save`, `devradar saved`, and `devradar unsave`.
- **Multiple Output Formats**: Styled Lip Gloss console table, valid indented JSON, and CSV export.
- **API Reliability & Rate Limiter**: Headroom tracking, quota warnings, and transparent support for personal access tokens.
