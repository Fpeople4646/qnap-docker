# qnap-docker

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/qnap-docker)](https://goreportcard.com/report/github.com/scttfrdmn/qnap-docker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/scttfrdmn/qnap-docker)](https://github.com/scttfrdmn/qnap-docker/releases/latest)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Integration Tests](https://img.shields.io/badge/integration%20tests-passing-brightgreen.svg)](#integration-tests)

## 🚀 Comprehensive Docker Management CLI for QNAP NAS with Container Station

**qnap-docker** is the complete Docker management solution for QNAP NAS devices with Container Station. Deploy, manage, and monitor Docker containers on your QNAP NAS with 40+ commands covering the full Docker workflow - container lifecycle, networking, volumes, images, and system operations. Perfect for home labs, self-hosting, and production deployments on QNAP NAS.

**✅ Verified Working** on real QNAP hardware with comprehensive integration testing.

> **Sister Project**: [syno-docker](https://github.com/scttfrdmn/syno-docker) - Comprehensive Docker management for Synology NAS with Container Manager

## Features

### Core Deployment
- 🚀 **One-command deployment** - Deploy containers as easily as `qnap-docker run nginx`
- 📦 **Docker Compose support** - Deploy complex multi-container applications
- 🔧 **Dynamic Docker detection** - Automatically finds Container Station across volumes
- 📂 **Multi-volume support** - Smart handling of CACHEDEV, ZFS, USB, external volumes

### Container Management
- 🔄 **Complete lifecycle** - Start, stop, restart, remove containers
- 📋 **Container inspection** - Detailed container information and logs
- 🖥️ **Interactive execution** - Run commands inside containers (`exec`)
- 📊 **Resource monitoring** - Real-time container statistics

### Image & System Management
- 🏗️ **Image operations** - Pull, list, remove images with advanced filtering
- 📦 **Volume management** - Create, list, inspect, and clean up volumes
- 🌐 **Network management** - Create, list, inspect networks; connect/disconnect containers
- 🧹 **System maintenance** - Disk usage, system info, and cleanup tools
- 📤 **Import/Export** - Backup and restore containers

### Infrastructure
- 🔐 **SSH key & ssh-agent support** - Works with both SSH key files and ssh-agent
- 👤 **Administrator user support** - Compatible with both `admin` and custom admin users
- 🎯 **Container Station optimized** - Built specifically for QNAP Container Station
- ⚡ **Single binary** - No dependencies, just download and use
- 🧪 **Integration tested** - Verified on real QNAP hardware

## Quick Start

### Installation

**Multiple installation methods for macOS, Linux, and direct download:**

```bash
# Install via Homebrew (macOS/Linux) - Recommended
brew tap scttfrdmn/qnap-docker
brew install qnap-docker

# Or in one command:
brew install scttfrdmn/qnap-docker/qnap-docker

# Direct binary download (all platforms)
curl -L https://github.com/scttfrdmn/qnap-docker/releases/latest/download/qnap-docker-$(uname -s)-$(uname -m) -o qnap-docker
chmod +x qnap-docker
sudo mv qnap-docker /usr/local/bin/

# Linux packages (Ubuntu/Debian/CentOS/Alpine)
# Download .deb/.rpm/.apk from releases page
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

# List running containers and monitor resources
qnap-docker ps
qnap-docker stats

# Get container logs and execute commands
qnap-docker logs web-server --follow
qnap-docker exec web-server /bin/bash

# Remove container when done
qnap-docker rm web-server
```

## Commands Overview

qnap-docker provides **21 main commands + 18 subcommands** covering the complete Docker workflow:

### **Container Lifecycle**
- `qnap-docker run` - Deploy single containers with full configuration options
- `qnap-docker ps` - List containers (running/all) with detailed status
- `qnap-docker start/stop/restart` - Control container state
- `qnap-docker rm` - Remove containers (with force option)

### **Container Operations**
- `qnap-docker logs` - View container logs (follow, tail, timestamps)
- `qnap-docker exec` - Execute commands inside containers (interactive/non-interactive)
- `qnap-docker stats` - Real-time resource usage statistics
- `qnap-docker inspect` - Detailed container/image/volume information

### **Image Management**
- `qnap-docker pull` - Pull images from registries (platform-specific, all tags)
- `qnap-docker images` - List images (all, dangling, with digests)
- `qnap-docker rmi` - Remove images (force, preserve parents)
- `qnap-docker import/export` - Backup and restore containers

### **Volume Management**
- `qnap-docker volume ls` - List volumes with driver information
- `qnap-docker volume create` - Create volumes with custom drivers/labels
- `qnap-docker volume rm` - Remove volumes (with force)
- `qnap-docker volume inspect` - Detailed volume information
- `qnap-docker volume prune` - Clean unused volumes

### **Network Management**
- `qnap-docker network ls` - List networks with filtering options
- `qnap-docker network create` - Create custom networks with CIDR, gateways
- `qnap-docker network rm` - Remove networks
- `qnap-docker network inspect` - Detailed network information
- `qnap-docker network connect/disconnect` - Attach/detach containers
- `qnap-docker network prune` - Clean unused networks

### **System Operations**
- `qnap-docker system df` - Show Docker disk usage
- `qnap-docker system info` - Display Docker system information
- `qnap-docker system prune` - Clean unused containers, images, networks

### **Multi-Container Applications**
- `qnap-docker deploy` - Deploy from docker-compose.yml files
- `qnap-docker init` - Setup connection to QNAP NAS

### **Key Command Examples**

```bash
# Container lifecycle
qnap-docker run nginx:latest --name web --port 80:80 --restart unless-stopped
qnap-docker logs web --follow --timestamps
qnap-docker exec -it web /bin/bash
qnap-docker restart web
qnap-docker stop web && qnap-docker rm web

# Image management
qnap-docker pull postgres:13 --platform linux/arm64
qnap-docker images --dangling
qnap-docker rmi old-image --force

# Volume operations
qnap-docker volume create my-data --driver local
qnap-docker volume ls --quiet
qnap-docker volume inspect my-data
qnap-docker volume rm my-data --force

# Network operations
qnap-docker network create my-app-net --driver bridge --subnet 172.20.0.0/16
qnap-docker network ls --filter driver=bridge
qnap-docker network connect my-app-net web-server --alias web
qnap-docker network disconnect my-app-net web-server

# System maintenance
qnap-docker system df --verbose
qnap-docker system prune --all --volumes --force
qnap-docker stats --all --no-stream

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

### `qnap-docker logs <container>`

View container logs with advanced options.

```bash
# Show recent logs
qnap-docker logs web-server

# Follow logs in real-time
qnap-docker logs web-server --follow --timestamps

# Show last 50 lines since 1 hour ago
qnap-docker logs web-server --tail 50 --since 1h
```

### `qnap-docker exec <container> <command>`

Execute commands inside running containers.

```bash
# Interactive shell
qnap-docker exec -it web-server /bin/bash

# Run single command
qnap-docker exec web-server cat /etc/hostname

# Run as specific user
qnap-docker exec --user 1000:1000 web-server whoami
```

### `qnap-docker stats [containers...]`

Display live container resource usage.

```bash
# Show stats for all running containers
qnap-docker stats

# Show stats for specific containers
qnap-docker stats web-server database

# One-time stats (no streaming)
qnap-docker stats --no-stream --all
```

### `qnap-docker volume <subcommand>`

Manage Docker volumes.

```bash
# List volumes
qnap-docker volume ls

# Create volume with custom driver
qnap-docker volume create my-data --driver local --label env=prod

# Inspect volume details
qnap-docker volume inspect my-data

# Clean unused volumes
qnap-docker volume prune --force
```

### `qnap-docker network <subcommand>`

Manage Docker networks.

```bash
# List networks
qnap-docker network ls

# Create custom network
qnap-docker network create app-net --driver bridge --subnet 172.20.0.0/16

# Connect container to network
qnap-docker network connect app-net web-server --alias web

# Remove unused networks
qnap-docker network prune --force
```

### `qnap-docker system <subcommand>`

System-wide Docker management.

```bash
# Show disk usage
qnap-docker system df --verbose

# Display system information
qnap-docker system info

# Clean unused data (containers, images, networks)
qnap-docker system prune --all --volumes --force
```

### `qnap-docker rm <container>`

Remove containers.

```bash
# Remove stopped container
qnap-docker rm web-server

# Force remove running container
qnap-docker rm web-server --force

# Remove multiple containers
qnap-docker rm web-server database cache-server
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

qnap-docker automatically handles QNAP volume paths across different storage types:

```bash
# These are equivalent (auto-detects your primary volume):
qnap-docker run nginx -v /share/ZFS530_DATA/web:/usr/share/nginx/html
qnap-docker run nginx -v ./web:/usr/share/nginx/html  # Expands to detected volume + /docker/web
qnap-docker run nginx -v web:/usr/share/nginx/html    # Expands to detected volume + /docker/web

# Supports multiple volume types:
qnap-docker run app -v /share/CACHEDEV1_DATA/data:/app/data  # Traditional CACHEDEV
qnap-docker run app -v /share/ZFS530_DATA/data:/app/data     # ZFS storage pools
qnap-docker run app -v /share/USB/backup:/app/backup        # USB storage
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
- Container Station binary is not accessible
- Your volumes are not mounted (check `/share/` directory)

qnap-docker uses **dynamic detection** to find Container Station automatically.

### Permission Denied

```bash
# Check if Container Station is accessible
ssh admin@192.168.1.100 'find /share -name docker -type f | grep container-station'

# Ensure your user has access to Container Station
ssh admin@192.168.1.100 'ls -la /share/*/.*qpkg/container-station/*/docker'
```

### Port Already in Use

```bash
# Check what's using the port
ssh admin@192.168.1.100 'netstat -tlnp | grep :8080'
```

### Multiple Volumes

qnap-docker automatically detects all available volumes. You can specify a specific one during init:

```bash
# Use specific CACHEDEV volume
qnap-docker init your-nas.local --volume-path /share/CACHEDEV2_DATA/docker

# Use ZFS volume
qnap-docker init your-nas.local --volume-path /share/ZFS530_DATA/docker

# Use USB storage
qnap-docker init your-nas.local --volume-path /share/USB/docker
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

qnap-docker includes comprehensive integration tests that validate all 40+ commands against real QNAP hardware:

```bash
# Comprehensive test suite for all v0.2.x commands
QNAP_HOST=your-qnap.local go test -v -integration -run TestComprehensiveCommandSuite \
  -nas-host=your-qnap.local -nas-user=admin ./tests/integration/

# Test specific command categories
QNAP_HOST=your-qnap.local go test -v -integration -run TestComprehensiveCommandSuite/ContainerOperations ./tests/integration/
QNAP_HOST=your-qnap.local go test -v -integration -run TestComprehensiveCommandSuite/NetworkManagement ./tests/integration/

# Legacy basic tests
QNAP_HOST=your-qnap.local go test -v -run TestConnectionToQNAP ./tests/integration/
QNAP_HOST=your-qnap.local go test -v -run TestQNAPDockerEndToEnd ./tests/integration/

# All integration tests
QNAP_HOST=your-qnap.local go test -v -integration -nas-host=your-qnap.local ./tests/integration/
```

**Comprehensive integration test coverage (v0.2.2+):**
- ✅ **Container Operations**: logs, exec, start/stop/restart, stats with real scenarios
- ✅ **Image Management**: pull, images, rmi, export/import with registry interactions
- ✅ **Volume Management**: volume lifecycle, mounting, data persistence validation
- ✅ **Network Management**: network creation, container connectivity, isolation testing
- ✅ **System Operations**: system df/info/prune with actual resource cleanup
- ✅ **Advanced Features**: inspect, backup/restore workflows
- ✅ **Error Handling**: Invalid configurations and failure scenarios
- ✅ **Resource Cleanup**: Comprehensive cleanup verification
- ✅ SSH connectivity and authentication (ssh-agent + key file)
- ✅ Container Station dynamic detection and path validation
- ✅ Multi-volume support (CACHEDEV, ZFS) testing

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
| **Docker Binary Path** | `/usr/local/bin/docker` | **Dynamic detection** (CACHEDEV/ZFS volumes) |
| **Service Management** | systemd (`pkg-ContainerManager-dockerd`) | System V / QPKG |
| **Volume Paths** | `/volume1/`, `/volume2/` | `/share/CACHEDEV*_DATA/`, `/share/ZFS*_DATA/` |
| **Multi-Container Support** | Docker only | Docker + LXD + Kata (future) |
| **Volume Detection** | Static volume enumeration | **Dynamic multi-volume detection** |
| **Storage Types** | Synology volumes only | CACHEDEV, ZFS, USB, external |

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
- 🗺️ [Development Roadmap](ROADMAP.md)
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/qnap-docker/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/qnap-docker/discussions)
- ☕ [Donate](https://ko-fi.com/scttfrdmn) - Support development

---

**Made with ❤️ for the QNAP community**