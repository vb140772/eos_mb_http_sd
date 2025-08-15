# GitHub Actions Workflows

This document describes the GitHub Actions workflows available in this repository for building, testing, and publishing Docker images.

## 📋 **Available Workflows**

### 1. **Docker Build and Publish** (`docker.yml`)
**Triggers**: Tags, main/master branches, manual dispatch
**Purpose**: Build and publish Docker images to GitHub Container Registry and GitHub Packages

### 2. **Docker Test** (`docker-test.yml`)
**Triggers**: Pull requests, manual dispatch
**Purpose**: Test Docker images without publishing them

### 3. **Release** (`release.yml`)
**Triggers**: Tags, manual dispatch
**Purpose**: Create GitHub releases with binaries and Docker images

## 🚀 **Docker Build and Publish Workflow**

### **Overview**
The `docker.yml` workflow automatically builds and publishes Docker images when:
- You push a tag (e.g., `v1.0.0`)
- You push to main/master branches
- You manually trigger the workflow

### **Features**
- **Multi-platform support**: Builds for linux/amd64 and linux/arm64
- **Dual registry publishing**: Publishes to both GitHub Container Registry and GitHub Packages
- **Automatic tagging**: Creates semantic version tags and latest tags
- **Security scanning**: Runs Trivy vulnerability scanner on built images
- **Caching**: Uses GitHub Actions cache for faster builds

### **Published Images**

#### **GitHub Container Registry (ghcr.io)**
```
ghcr.io/{username}/{repository}:{version}
ghcr.io/{username}/{repository}:latest
ghcr.io/{username}/{repository}:{major}.{minor}
ghcr.io/{username}/{repository}:{major}
```

#### **GitHub Packages (docker.pkg.github.com)**
```
docker.pkg.github.com/{username}/{repository}/eos-mb-http-sd:{version}
docker.pkg.github.com/{username}/{repository}/eos-mb-http-sd:latest
```

### **Manual Trigger**
You can manually trigger this workflow with custom options:
1. Go to **Actions** → **Docker Build and Publish**
2. Click **Run workflow**
3. Configure options:
   - **Version**: Custom version tag (optional)
   - **Push to Registry**: Enable/disable GitHub Container Registry
   - **Push to Packages**: Enable/disable GitHub Packages

## 🧪 **Docker Test Workflow**

### **Overview**
The `docker-test.yml` workflow runs on pull requests to ensure Docker images build and function correctly without publishing them.

### **Features**
- **Build testing**: Ensures Docker images can be built successfully
- **Runtime testing**: Tests that containers start and respond to health checks
- **Security scanning**: Runs Trivy vulnerability scanner
- **Dockerfile linting**: Uses hadolint to check Dockerfile best practices

### **What It Tests**
1. **Image building**: Verifies Docker images build without errors
2. **Container startup**: Ensures containers start successfully
3. **Health checks**: Tests health endpoint availability
4. **Vulnerability scanning**: Checks for critical and high-severity vulnerabilities
5. **Dockerfile quality**: Validates Dockerfile follows best practices

## 📦 **Release Workflow**

### **Overview**
The `release.yml` workflow creates comprehensive GitHub releases when you push tags, including:
- Binary artifacts for multiple platforms
- Docker images for both registries
- Security scan results
- Detailed release notes

### **Features**
- **Multi-platform binaries**: Builds for linux, darwin, and windows (amd64/arm64)
- **Docker integration**: Automatically builds and publishes Docker images
- **Security scanning**: Includes Trivy vulnerability scan results
- **Rich release notes**: Automatically generates release notes with Docker image information

### **Release Process**
1. **Binary building**: Creates binaries for all supported platforms
2. **Docker building**: Builds multi-platform Docker images
3. **Security scanning**: Scans images for vulnerabilities
4. **Release creation**: Creates GitHub release with all artifacts
5. **Asset upload**: Uploads binaries, checksums, and Docker images

## 🔧 **Configuration**

### **Required Secrets**
No additional secrets are required. The workflows use the built-in `GITHUB_TOKEN` secret.

### **Permissions**
The workflows automatically request the necessary permissions:
- `contents: read/write` - For reading code and creating releases
- `packages: write` - For publishing Docker images
- `id-token: write` - For OIDC token generation
- `security-events: write` - For uploading security scan results

## 📊 **Workflow Status**

### **Success Indicators**
- ✅ **Green checkmark**: All jobs completed successfully
- 🔄 **Yellow circle**: Workflow is running
- ❌ **Red X**: One or more jobs failed

### **Job Dependencies**
```
docker.yml:
├── build-and-push (parallel matrix builds)
├── security-scan (depends on build-and-push)
└── notify (depends on both)

docker-test.yml:
├── test-docker (parallel matrix builds)
└── lint-dockerfile (independent)

release.yml:
├── build (parallel matrix builds)
├── docker-build (depends on build)
└── release (depends on both)
```

## 🚨 **Troubleshooting**

### **Common Issues**

#### **Build Failures**
- Check that your Dockerfile is valid
- Ensure all required files are present
- Verify Go module dependencies are correct

#### **Authentication Errors**
- Ensure the repository has the necessary permissions
- Check that GitHub Actions is enabled for the repository
- Verify the workflow has access to required secrets

#### **Push Failures**
- Check that the target registry is accessible
- Verify image naming conventions
- Ensure the workflow has write permissions

### **Debug Steps**
1. **Check workflow logs**: Review the detailed logs for each step
2. **Verify permissions**: Ensure the repository has the required permissions
3. **Check dependencies**: Verify all required files and configurations are present
4. **Review triggers**: Confirm the workflow is triggered by the expected events

## 📈 **Performance Optimization**

### **Build Caching**
- **Docker layer caching**: Uses GitHub Actions cache for Docker layers
- **Go module caching**: Caches Go dependencies for faster builds
- **Multi-platform builds**: Parallel builds for different architectures

### **Optimization Tips**
1. **Use .dockerignore**: Exclude unnecessary files from Docker context
2. **Multi-stage builds**: Leverage multi-stage Dockerfiles for smaller images
3. **Layer optimization**: Order Dockerfile instructions for better caching
4. **Parallel jobs**: Use matrix strategies for concurrent builds

## 🔒 **Security Features**

### **Vulnerability Scanning**
- **Trivy integration**: Automated vulnerability scanning of built images
- **Security tab integration**: Results uploaded to GitHub Security tab
- **Fail-fast**: Workflows can be configured to fail on critical vulnerabilities

### **Best Practices**
- **Non-root containers**: Images run as non-root users
- **Minimal base images**: Uses Alpine Linux for smaller attack surface
- **Regular updates**: Base images are updated with each build
- **Security scanning**: Automated scanning on every build

## 📚 **Additional Resources**

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitHub Packages](https://docs.github.com/en/packages)
- [Docker Buildx](https://docs.docker.com/buildx/)
- [Trivy Security Scanner](https://aquasecurity.github.io/trivy/)
- [Hadolint](https://github.com/hadolint/hadolint)
