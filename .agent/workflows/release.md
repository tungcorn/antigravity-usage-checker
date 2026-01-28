---
description: Create a new release with changelog and verification
---

# Release Workflow

Follow these steps to create a new release:

## 1. Confirm version
- Ask user for the new version (e.g., v2.1.0)
- Check the latest version in CHANGELOG.md to avoid duplicates

## 2. Update CHANGELOG.md
- Add new section for the version at the top (after header)
- Use existing format:
```markdown
## [vX.X.X] Release Title

### Added
- New features

### Changed
- Changes

### Fixed
- Bug fixes

---
```
- Ask user for change descriptions if not provided

## 3. Build verification
// turbo
- Run: `go build -o agusage.exe ./cmd/agusage/`
- Ensure build succeeds
- Remove test build file after verification

## 4. Commit and push changes
// turbo
- Run: `git add . && git commit -m "release: vX.X.X - <short summary of changes>" && git tag vX.X.X && git push origin main && git push origin vX.X.X`

## 6. Confirm completion
- Notify user that workflow is complete
- Provide link to GitHub Actions: https://github.com/tungcorn/antigravity-usage-checker/actions