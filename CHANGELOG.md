# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.1] Version Display Fix

### Fixed
- Version now correctly displays from release tag instead of hardcoded value

### Changed
- Auto-inject version via ldflags during release build

---

## [v1.0.0] Enhanced Reset Time Display

### Added
- Display exact reset time alongside countdown (e.g., "1h 12m (01:40)")

### Changed
- Improved reset time display for better user experience

---

## [v0.5.2] Refactoring & Bug Fixes

### Changed
- Cleanup unused code and improve error handling
- Update GitHub username from TungCorn to tungcorn

### Added
- Version parameter to install scripts

### Removed
- TOPICS.md file

---

## [v0.5.1] Code Quality Improvements

### Fixed
- Update mock tests to match new struct-based client implementation

### Changed
- Refactor: improve code quality and cleanup

---

## [v0.5.0] Enhanced CLI Display

### Added
- Enhanced CLI display with smart features
- Individual OS badges in README (Windows, macOS, Linux)

### Changed
- Replace generic platform badge with individual OS badges

---

## [v0.4.0] Multi-platform Support

### Added
- Multi-platform support (Windows, macOS, Linux)
- Unit tests and GitHub Actions CI/CD workflows
- One-line install script for easy installation
- Windows release artifact

### Changed
- Improved installation documentation
- Use dot-sourcing for immediate PATH update

---

## [v0.3.0] Installation Improvements

### Added
- Install script for easier setup
- ZIP download instructions in README

---

## [v0.2.0] Quota Display Enhancement

### Added
- Reset time column to quota display
- Screenshot in documentation
- Star request in README

---

## [v0.1.0] Initial Release

### Added
- CLI tool to analyze Antigravity (Gemini CLI) usage statistics
- Support for Windows, macOS, and Linux
- Display model usage breakdown with token counts
- Show conversation statistics and history
- Cache analysis with size information
- Multiple output formats: table, JSON, plain text
- Color-coded terminal output
- Bilingual README (English/Vietnamese)
- CLI entrypoint with flags
- Terminal UI with progress bars
- Cache fallback for graceful degradation
- API client for GetUserStatus endpoint
- OAuth credentials loader with retry
- Process discovery for Antigravity server
