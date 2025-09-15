# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-01-15

### Added - Complete Docker Management Suite
- **Container Operations**: `logs`, `exec`, `start`, `stop`, `restart`, `stats` commands
- **Image Management**: `images`, `pull`, `rmi`, `import`, `export` commands
- **Volume Management**: `volume ls/create/rm/inspect/prune` subcommands
- **Network Management**: `network ls/create/rm/inspect/connect/disconnect/prune` subcommands
- **System Management**: `system df/info/prune` subcommands
- **Advanced Inspection**: `inspect` command for detailed object information

### Enhanced Features
- **Interactive Container Execution**: Full TTY and interactive mode support
- **Real-time Log Following**: Stream container logs with timestamps and filtering
- **Resource Monitoring**: Live container statistics and system resource usage
- **Advanced Image Operations**: Platform-specific pulls, digest support, dangling image cleanup
- **Network Isolation**: Custom bridge networks with CIDR configuration
- **Volume Persistence**: Named volume creation with driver options
- **System Maintenance**: Comprehensive cleanup and disk usage monitoring

### Total Command Count
- **23 main commands + 18 subcommands** = Complete Docker workflow coverage
- Full feature parity with syno-docker v0.2.0
- All commands adapted for QNAP Container Station architecture

### Architecture Improvements
- Dynamic Docker binary detection for Container Station
- ZFS volume support (ZFS*_DATA) in addition to CACHEDEV volumes
- Docker binary path caching for improved performance
- Configurable integration tests via environment variables
- Enhanced volume detection for multiple storage pool types

### Fixed
- Container Station detection now works with ZFS storage pools
- Integration tests properly configurable (QNAP_HOST, QNAP_USER)
- Path validation supports both CACHEDEV and ZFS volume patterns

## [0.1.0] - 2025-01-15

### Added
- Initial release of qnap-docker CLI tool
- SSH connection management with key and ssh-agent support
- QNAP Container Station integration with dynamic binary detection
- Multi-volume detection (CACHEDEV, ZFS, USB, external)
- Docker container lifecycle management (run, ps, rm)
- Docker Compose deployment support
- Multi-platform builds (darwin/linux, amd64/arm64)
- Comprehensive documentation and examples
- Integration test framework with real hardware validation
- Quality checks and Go Report Card A+ compliance

### Features
- `qnap-docker init` - Setup connection to QNAP NAS with volume detection
- `qnap-docker run` - Deploy single containers with QNAP path helpers
- `qnap-docker deploy` - Deploy from docker-compose.yml files
- `qnap-docker ps` - List containers with status information
- `qnap-docker rm` - Remove containers (with force option)

### Architecture Support
- **QNAP Storage**: CACHEDEV*_DATA, ZFS*_DATA, USB, external volumes
- **Container Station**: Dynamic Docker binary detection across volumes
- **Docker Version**: Tested with Docker 27.1.2-qnap4
- **Volume Detection**: Automatic discovery of available storage pools

### Supported Platforms
- QNAP NAS devices with Container Station installed
- QTS 4.5.4+ / QuTS hero h5.0.1+ / QuTScloud c5.1.0+
- macOS (darwin/amd64, darwin/arm64)
- Linux (linux/amd64, linux/arm64)

### Verified Hardware
- ✅ **Real QNAP Testing**: Validated on live QNAP NAS with Container Station
- ✅ **Multi-Volume**: CACHEDEV and ZFS storage pool support
- ✅ **Integration Tests**: SSH, Docker commands, container operations
- ✅ **Path Resolution**: Automatic QNAP volume path handling