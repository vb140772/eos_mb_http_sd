# Release Guide

This project uses GitHub Actions to automatically build binaries and create GitHub Releases.

## Automatic Release Process

### 1. Create a Git Tag

To trigger an automatic release, create and push a git tag:

```bash
# Create a new tag
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

The workflow will automatically:
- Build binaries for multiple platforms (Linux, macOS, Windows)
- Run tests to ensure quality
- Create a GitHub Release with the tag
- Upload all built binaries as release assets
- Generate SHA256 checksums for verification

### 2. Manual Release

You can also trigger a release manually through the GitHub Actions tab:

1. Go to the **Actions** tab in your GitHub repository
2. Select the **Build and Release** workflow
3. Click **Run workflow**
4. Enter the version number (e.g., `v1.0.0`)
5. Click **Run workflow**

## Supported Platforms

The workflow builds binaries for the following platforms:

- **Linux**: AMD64, ARM64
- **macOS**: AMD64, ARM64 (Apple Silicon)
- **Windows**: AMD64, ARM64

## Release Assets

Each release includes:

- Binary files for each platform
- SHA256 checksum files for verification
- A combined checksums.txt file

## Binary Names

Binaries are named according to the pattern:
```
eos_mb_http_sd-{platform}-{architecture}
```

Examples:
- `eos_mb_http_sd-linux-amd64`
- `eos_mb_http_sd-darwin-arm64`
- `eos_mb_http_sd-windows-amd64.exe`

## Verification

To verify the integrity of downloaded binaries:

```bash
# On Linux/macOS
shasum -a 256 -c eos_mb_http_sd-linux-amd64.sha256

# On Windows (PowerShell)
Get-FileHash eos_mb_http_sd-windows-amd64.exe -Algorithm SHA256
```

## Workflow Configuration

The workflow is configured in `.github/workflows/release.yml` and includes:

- **Build Job**: Compiles binaries for all platforms
- **Test Job**: Runs the test suite before building
- **Release Job**: Creates GitHub Release and uploads assets
- **Matrix Strategy**: Builds for multiple platforms in parallel

## Requirements

- Go 1.24.6 or later
- GitHub repository with Actions enabled
- Proper permissions for the `GITHUB_TOKEN`

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure the workflow has `contents: write` permission
2. **Build Failures**: Check that all tests pass and dependencies are available
3. **Release Creation Fails**: Verify the tag format follows `v*` pattern
4. **GitHub CLI Authentication Errors**: Ensure the repository has proper permissions

### GitHub Token Permissions

If you encounter authentication errors with the GitHub CLI in the workflow:

1. **Check Repository Settings**: Go to Settings → Actions → General
2. **Workflow Permissions**: Ensure "Read and write permissions" is selected
3. **Allow GitHub Actions to create and approve pull requests**: This may be needed for some operations

### Debugging Release Issues

The workflow includes extensive debugging output:
- Artifact listing and verification
- GitHub CLI authentication status
- Repository access verification
- Detailed command execution logging

Check the workflow logs in the Actions tab for detailed error information.

### Manual Binary Building

If you need to build binaries manually:

```bash
# Set environment variables
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Build binary
go build -ldflags="-s -w" -o eos_mb_http_sd-linux-amd64 main.go

# Create checksum
shasum -a 256 eos_mb_http_sd-linux-amd64 > eos_mb_http_sd-linux-amd64.sha256
```

## Versioning

Follow [Semantic Versioning](https://semver.org/) for version numbers:

- **Major**: Breaking changes
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, backward compatible

Example: `v1.2.3`
