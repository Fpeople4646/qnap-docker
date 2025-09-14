# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of qnap-docker CLI tool
- SSH connection management with key and ssh-agent support
- QNAP Container Station integration
- Dynamic CACHEDEV volume detection
- Docker container lifecycle management (run, ps, rm)
- Docker Compose deployment support
- Multi-platform builds (darwin/linux, amd64/arm64)
- Comprehensive documentation and examples
- Integration test framework
- Quality checks and Go Report Card A+ compliance

### Features
- `qnap-docker init` - Setup connection to QNAP NAS
- `qnap-docker run` - Deploy single containers
- `qnap-docker deploy` - Deploy from docker-compose.yml
- `qnap-docker ps` - List containers
- `qnap-docker rm` - Remove containers

### Supported Platforms
- QNAP NAS devices with Container Station
- QTS 4.5.4+ / QuTS hero h5.0.1+ / QuTScloud c5.1.0+
- macOS (darwin/amd64, darwin/arm64)
- Linux (linux/amd64, linux/arm64)

## [0.1.0] - 2025-01-XX (Planned)

### Added
- Initial public release
- Core functionality for QNAP Container Station
- Sister project to syno-docker for QNAP NAS systems