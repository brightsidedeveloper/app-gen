# Release Guide

This guide explains how to release new versions of the template-cli tool.

## Automated Release (Recommended)

The easiest way to release is using GitHub Actions:

1. **Update version and commit:**
   ```bash
   git add .
   git commit -m "Release v1.0.0"
   git push origin main
   ```

2. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for all platforms (Linux, macOS, Windows, both amd64 and arm64)
   - Create checksums
   - Create a GitHub release with all binaries attached

## Manual Release

If you prefer to build locally:

1. **Run the release script:**
   ```bash
   ./scripts/release.sh 1.0.0
   ```

2. **This creates binaries in the `dist/` directory:**
   - `template-cli-linux-amd64`
   - `template-cli-linux-arm64`
   - `template-cli-darwin-amd64`
   - `template-cli-darwin-arm64`
   - `template-cli-windows-amd64.exe`
   - `template-cli-windows-arm64.exe`
   - `checksums.txt`

3. **Create a GitHub release:**
   - Go to GitHub → Releases → Draft a new release
   - Tag: `v1.0.0`
   - Title: `v1.0.0`
   - Upload all files from the `dist/` directory
   - Publish release

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR.MINOR.PATCH** (e.g., 1.0.0)
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

## Testing Before Release

Before releasing, test the binary on your platform:

```bash
go build -o template-cli
./template-cli --name test-project
```

Make sure it works correctly before tagging and releasing!

