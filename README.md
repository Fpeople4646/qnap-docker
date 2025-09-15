# qnap-docker

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/qnap-docker)](https://goreportcard.com/report/github.com/scttfrdmn/qnap-docker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/badge/release-v0.1.0-blue.svg)](https://github.com/scttfrdmn/qnap-docker/releases/tag/v0.1.0)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

**qnap-docker** is a CLI tool that simplifies Docker container deployment to QNAP NAS devices with Container Station. It handles SSH connection management, Docker client setup, and path resolution issues specific to QNAP Container Station.

Sister project to [syno-docker](https://github.com/scttfrdmn/syno-docker) for Synology NAS systems.

## Features

- 🚀 **One-command deployment** - Deploy containers as easily as `qnap-docker run nginx`
- 🔐 **SSH key & ssh-agent support** - Works with both SSH key files and ssh-agent
- 👤 **Administrator user support** - Compatible with both `admin` and custom admin users
- 📦 **docker-compose support** - Deploy complex multi-container applications
- 🎯 **Container Station optimized** - Built specifically for QNAP Container Station
- 🔧 **PATH resolution** - Automatically handles Docker binary location issues
- 📂 **CACHEDEV path helpers** - Smart handling of QNAP CACHEDEV volume paths
- 🔄 **Container lifecycle** - Deploy, list, and remove containers easily
- ⚡ **Single binary** - No dependencies, just download and use
- 🧪 **Integration tested** - Verified on real QNAP hardware

## Quick Start

### Installation

```bash
# Install via Homebrew (recommended)
brew install scttfrdmn/tap/qnap-docker

# Or download binary from releases
curl -L https://github.com/scttfrdmn/qnap-docker/releases/latest/download/qnap-docker-$(uname -s)-$(uname -m) -o qnap-docker
chmod +x qnap-docker
sudo mv qnap-docker /usr/local/bin/
```

### Setup

```bash
# One-time setup - connect to your QNAP NAS
qnap-docker init 192.168.1.100

# Or with custom admin username (if not using 'admin')
qnap-docker init your-nas.local --user your-username

# For ssh-agent users (automatically detected)
qnap-docker init your-nas.local --user your-username
```

This will:
- Test SSH connection to your NAS (supports both SSH keys and ssh-agent)
- Verify Container Station is running
- Test Docker command execution
- Detect available CACHEDEV volumes
- Save connection details to `~/.qnap-docker/config.yaml`

### Deploy Your First Container

```bash
# Deploy Nginx web server
qnap-docker run nginx:latest \
  --name web-server \
  --port 8080:80 \
  --volume /share/CACHEDEV1_DATA/web:/usr/share/nginx/html

# Deploy from docker-compose.yml
qnap-docker deploy ./docker-compose.yml

# List running containers
qnap-docker ps

# Remove container
qnap-docker rm web-server
```

## Commands

### `qnap-docker init <host>`

Setup connection to your QNAP NAS.

```bash
qnap-docker init 192.168.1.100 \
  --user admin \
  --port 22 \
  --key ~/.ssh/id_rsa \
  --volume-path /share/CACHEDEV1_DATA/docker
```

### `qnap-docker run <image>`

Deploy a single container.

```bash
qnap-docker run postgres:13 \
  --name database \
  --port 5432:5432 \
  --volume /share/CACHEDEV1_DATA/postgres:/var/lib/postgresql/data \
  --env POSTGRES_PASSWORD=secretpassword \
  --restart unless-stopped
```

**Options:**
- `--name` - Container name (auto-generated if not specified)
- `--port` - Port mappings (format: `host:container`)
- `--volume` - Volume mappings (format: `host:container`)
- `--env` - Environment variables (format: `KEY=value`)
- `--restart` - Restart policy (`no`, `always`, `unless-stopped`, `on-failure`)
- `--network` - Network mode (default: `bridge`)
- `--user` - User to run container as (format: `uid:gid`)
- `--workdir` - Working directory inside container
- `--command` - Command to run in container

### `qnap-docker deploy <compose-file>`

Deploy from docker-compose.yml file.

```bash
qnap-docker deploy ./docker-compose.yml \
  --project my-app \
  --env-file .env
```

**Supported compose features:**
- Multi-service deployments
- Port mappings
- Volume mounts
- Environment variables
- Environment variable substitution
- Restart policies
- Networks (basic support)
- Dependencies (deployment order only)

### `qnap-docker ps`

List containers.

```bash
# Show running containers
qnap-docker ps

# Show all containers (including stopped)
qnap-docker ps --all
```

### `qnap-docker rm <container>`

Remove container.

```bash
# Remove stopped container
qnap-docker rm web-server

# Force remove running container
qnap-docker rm web-server --force
```

## Configuration

qnap-docker stores configuration in `~/.qnap-docker/config.yaml`:

```yaml
host: 192.168.1.100
port: 22
user: admin
ssh_key_path: /home/user/.ssh/id_rsa
defaults:
  volume_path: /share/CACHEDEV1_DATA/docker
  primary_volume: CACHEDEV1_DATA
  network: bridge
```

## Volume Path Handling

qnap-docker automatically handles QNAP CACHEDEV volume paths:

```bash
# These are equivalent:
qnap-docker run nginx -v /share/CACHEDEV1_DATA/web:/usr/share/nginx/html
qnap-docker run nginx -v ./web:/usr/share/nginx/html  # Expands to /share/CACHEDEV1_DATA/docker/web
qnap-docker run nginx -v web:/usr/share/nginx/html    # Expands to /share/CACHEDEV1_DATA/docker/web
```

## Requirements

### QNAP NAS
- QTS 4.5.4 or later (or QuTS hero h5.0.1+, QuTScloud c5.1.0+)
- Container Station installed and running
- SSH access enabled (Control Panel → Network & File Services → Telnet/SSH)
- User with administrator privileges and docker access

### Local Machine
- SSH key pair configured OR ssh-agent running
- Network access to your NAS
- Go 1.21+ (for building from source)

## Troubleshooting

### Connection Issues

```bash
# Test SSH connection manually
ssh admin@192.168.1.100

# Check if Container Station is running
ssh admin@192.168.1.100 'test -x /share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker && echo "Docker found" || echo "Docker not found"'
```

### Container Station Not Found

This means:
- Container Station is not installed
- Container Station binary is not in the expected location
- Your CACHEDEV path differs (try different CACHEDEV numbers)

### Permission Denied

```bash
# Ensure your user has access to Container Station
ssh admin@192.168.1.100 'ls -la /share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker'
```

### Port Already in Use

```bash
# Check what's using the port
ssh admin@192.168.1.100 'netstat -tlnp | grep :8080'
```

### Multiple CACHEDEV Volumes

If you have multiple volumes, specify the correct one during init:

```bash
qnap-docker init your-nas.local --volume-path /share/CACHEDEV2_DATA/docker
```

## Development

### Building from Source

```bash
git clone https://github.com/scttfrdmn/qnap-docker.git
cd qnap-docker
make build
```

### Running Tests

```bash
make test              # Run unit tests
make quality-check     # Run all quality checks
make coverage         # Generate coverage report
```

### Integration Tests

qnap-docker includes comprehensive integration tests that validate functionality against real QNAP hardware:

```bash
# Test connection to your NAS
NAS_HOST=your-qnap.local make integration-test

# With custom user
NAS_HOST=your-qnap.local NAS_USER=your-username make integration-test

# Full end-to-end testing with coverage
NAS_HOST=your-qnap.local NAS_USER=admin make integration-test-full

# Test with environment variables (alternative)
QNAP_HOST=your-qnap.local QNAP_USER=admin go test ./tests/integration/ -v -short=false
```

**Integration test coverage:**
- ✅ SSH connectivity and authentication (ssh-agent + key file)
- ✅ Dynamic Container Station Docker binary detection
- ✅ Container deployment, lifecycle, and removal
- ✅ HTTP endpoint validation for deployed services
- ✅ Volume mounting and file system access (CACHEDEV + ZFS)
- ✅ Error handling for invalid configurations
- ✅ Multi-volume detection (CACHEDEV, ZFS) and validation
- ✅ Docker version compatibility testing
- ✅ CLI command interface validation

### Quality Checks

qnap-docker maintains Go Report Card A+ grade with:

- `gofmt` - Code formatting
- `go vet` - Static analysis
- `golangci-lint` - Comprehensive linting
- `staticcheck` - Advanced static analysis
- `ineffassign` - Ineffectual assignment detection
- `misspell` - Spelling mistakes
- `gocyclo` - Cyclomatic complexity

## Differences from syno-docker

| Feature | syno-docker | qnap-docker |
|---------|-------------|-------------|
| **Target Platform** | Synology DSM 7.2+ | QNAP QTS 4.5.4+ |
| **Container Platform** | Container Manager | Container Station |
| **Docker Binary Path** | `/usr/local/bin/docker` | `/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker` |
| **Service Management** | systemd (`pkg-ContainerManager-dockerd`) | System V / QPKG |
| **Volume Paths** | `/volume1/`, `/volume2/` | `/share/CACHEDEV1_DATA/`, `/share/CACHEDEV2_DATA/` |
| **Multi-Container Support** | Docker only | Docker + LXD + Kata (future) |
| **Volume Detection** | Static volume enumeration | Dynamic CACHEDEV detection |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks (`make quality-check`)
5. Run tests (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Docker SDK](https://github.com/docker/docker) - Docker client library
- [SSH package](https://golang.org/x/crypto/ssh) - SSH client implementation
- QNAP Community - For documenting Container Station architecture
- [syno-docker](https://github.com/scttfrdmn/syno-docker) - Sister project for Synology NAS

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/qnap-docker/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/qnap-docker/discussions)
- 📧 Email: support@qnap-docker.com

---

**Made with ❤️ for the QNAP community**