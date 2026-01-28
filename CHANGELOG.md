# Changelog

All notable changes to this project will be documented in this file.

## [v2.2.3] Improved Reset Time Display

### Changed
- **Reset Time Display**: Show "1d 2h" instead of "26h 33m" when reset time exceeds 24 hours for better readability

---

## [v2.2.2] Documentation Improvements

### Changed
- **Help Message**: Added `agu` alias information to help text
  - Users can now see that `agu` is available as a shorter alternative
  - Updated usage examples to include both `agusage` and `agu`

---

## [v2.2.1] Fix Alias Deployment

### Fixed
- **Install Script**: Changed `agu` alias from duplicate .exe file to .bat script (Windows)
  - `agu` now always calls the latest `agusage.exe` version
  - Fixes version mismatch issue when updating
  - Automatically removes old `agu.exe` if present
- **Release Install Instructions**: Now includes `-Version` parameter for accurate version installation
  - Ensures users install the exact release version they're viewing
  - Fixes issue where pre-releases would install latest stable instead

---

## [v2.2.0] Account Display & Deployment Improvements

### Added
- **Account Display**: Show active account email in Antigravity
  - Email is fetched from Antigravity Local Server API (most accurate source)
  - Automatically displays the currently active account

### Changed
- **Alias Deployment**: Improved `agu` alias creation
  - Changed from duplicate .exe file to .bat script
  - `agu` now automatically calls the latest `agusage` version
  - No manual copying needed when updating

---

## [v2.1.0] Add Alias Command

### Added
- **Alias Command**: Added `agu` as a shorthand alias for `agusage`
  -Windows: Creates `agu.exe` alongside `agusage.exe`
  - macOS/Linux: Creates symlink `agu` -> `agusage`
- Updated install scripts to automatically create the alias during installation
- Updated documentation to introduce both command names

---

## [v2.0.2] Fix Calculation & Display

### Fixed
- **Total Calculation**: Corrected total usage percentage (now uses weighted average).
- **Display Alignment**: Fixed broken table borders and column alignment.
- **Color Bleeding**: Fixed issue where text color spilled into subsequent rows.
- **Data Accuracy**: Improved detection of shared quotas using reset time.
- **UX**: Expanded "Reset" column to show full timestamp without truncation.

---

## [v2.0.1] Documentation & CI Updates

### Changed
- Add author to LICENSE
- Improved README with bilingual support (Vietnamese/English) and anchor links
- CI: Automatically append quick install instructions to release notes

---

## [v2.0.0] Cleaner Display & Percentage Format

### Changed
- Display usage as percentages (`Used %`) instead of misleading absolute numbers
- Removed redundant `Left %` column (same info as 100 - Used)
- Removed redundant percentage number from progress bar
- Expanded model name column from 22 to 30 characters for full names

### Fixed
- Clarified that values are percentages, not request counts

---

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
