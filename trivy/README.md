# Trivy Container Scanner

This directory contains configuration and scripts for using [Trivy](https://aquasecurity.github.io/trivy/), a comprehensive security scanner for containers and other artifacts. Trivy helps identify vulnerabilities, misconfigurations, and security issues in your container images.

## Prerequisites

- Docker installed on your system
- Docker daemon running
- Access to the container image you want to scan

## Usage

### Scanning Container Images

You can scan container images in two ways:

1. **Using environment variables**:
   Create a `.env` file with the following variables:
   ```env
   IMAGE_NAME=your-image-name
   IMAGE_TAG=your-image-tag
   ```

2. **Using command-line arguments**:
   ```bash
   make scan-with-docker image_name=your-image-name:your-image-tag
   ```

### Example

```bash
# Using environment variables
echo "IMAGE_NAME=nginx" > .env
echo "IMAGE_TAG=latest" >> .env
make scan-with-docker

# Using command-line arguments
make scan-with-docker image_name=nginx:latest
```

## What Trivy Scans

Trivy performs several types of scans:

- OS package vulnerabilities
- Language-specific package vulnerabilities
- Configuration issues
- Secret leaks
- License compliance

## Output

The scan results will show:
- Vulnerability details
- Severity levels
- CVE IDs
- Fix recommendations

## Version

This setup uses Trivy version 0.63.0, which is a stable and well-tested version.

## Additional Resources

- [Official Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [Trivy GitHub Repository](https://github.com/aquasecurity/trivy)
