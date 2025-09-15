# qnap-docker Integration Test Suite

This comprehensive integration test suite validates all 40+ qnap-docker commands against a real QNAP NAS running QTS 4.5.4+ with Container Station.

## 🎯 **Complete Test Coverage (v0.2.2+)**

### Comprehensive Command Testing
- **Container Operations**: logs, exec, start/stop/restart, stats with real container scenarios
- **Image Management**: pull, images, rmi, export/import with registry interactions
- **Volume Management**: Complete volume lifecycle, mounting, and data persistence validation
- **Network Management**: Network creation, container connectivity, multi-container communication
- **System Operations**: system df/info/prune with actual resource verification
- **Advanced Features**: inspect, backup/restore workflows with real data

### Legacy Core Functionality Tests
- **Basic Deployment**: Single container deployment with various configurations
- **Compose Deployment**: Multi-container applications using docker-compose.yml
- **Lifecycle Management**: Container start, stop, restart, and removal operations
- **Volume Mapping**: Host path mounting and named volume management (CACHEDEV/ZFS)
- **Network Connectivity**: Inter-container communication and external access
- **Error Handling**: Invalid configurations and failure scenarios

### Advanced Scenarios
- **Performance Tests**: Resource usage and deployment speed benchmarks
- **Stress Tests**: High load and concurrent deployment scenarios
- **Edge Cases**: Unusual configurations and boundary conditions
- **Security Tests**: Permission validation and access control
- **QNAP-Specific**: ZFS volume detection, Container Station path variations, multi-CACHEDEV testing

## 🛠️ **Prerequisites**

### QNAP NAS Requirements
- **QTS Version**: 4.5.4 or later (or QuTS hero h5.0.1+, QuTScloud c5.1.0+)
- **Container Station**: Installed and running
- **SSH Service**: Enabled (Control Panel → Network & File Services → Telnet/SSH)
- **User Account**: Admin privileges with Container Station access
- **Storage**: At least 2GB free space on primary volume
- **Network**: Accessible from test machine

### Test Environment Setup
```bash
# Required tools
- Go 1.21+
- SSH client with key-based authentication
- Network connectivity to QNAP NAS
- Docker client (for comparison/validation)

# Optional tools
- curl/wget (for HTTP endpoint testing)
- PostgreSQL client (for database connectivity tests)
```

## 📋 **Setup Instructions**

### 1. Prepare Your QNAP NAS

**Enable SSH Access:**
```bash
# On your NAS via SSH or QTS Control Panel
# Control Panel → Network & File Services → Telnet/SSH
# Enable SSH service and configure user access
```

**Create Test Directory:**
```bash
# SSH to your NAS
ssh admin@your-qnap-ip

# Create test directory (auto-detects available volumes)
mkdir -p /share/CACHEDEV1_DATA/qnap-docker-test
# Or for ZFS storage:
mkdir -p /share/ZFS530_DATA/qnap-docker-test

chmod 755 /share/*/qnap-docker-test
```

### 2. Configure SSH Keys
```bash
# Generate SSH key if needed
ssh-keygen -t rsa -b 4096 -C "qnap-docker-test"

# Copy public key to NAS
ssh-copy-id admin@your-qnap-ip

# Test SSH connection and Container Station
ssh admin@your-qnap-ip "find /share -name docker -type f | grep container-station | head -1"
```

### 3. Configure Test Environment
```bash
# Set environment variables
export QNAP_HOST="192.168.1.100"
export QNAP_USER="admin"
export NAS_HOST="192.168.1.100"  # Legacy compatibility
export NAS_USER="admin"          # Legacy compatibility
export NAS_SSH_KEY="~/.ssh/id_rsa"
```

## 🚀 **Running Tests**

### Comprehensive Test Run (v0.2.2+)
```bash
# Run complete command suite tests (all 40+ commands)
QNAP_HOST=your-qnap.local go test -v -integration \
    -nas-host=your-qnap.local \
    -nas-user=admin \
    -nas-key=~/.ssh/id_rsa \
    -run TestComprehensiveCommandSuite \
    ./tests/integration/

# Run specific command category tests
QNAP_HOST=your-qnap.local go test -v -integration \
    -nas-host=your-qnap.local \
    -run TestComprehensiveCommandSuite/ContainerOperations \
    ./tests/integration/

QNAP_HOST=your-qnap.local go test -v -integration \
    -nas-host=your-qnap.local \
    -run TestComprehensiveCommandSuite/NetworkManagement \
    ./tests/integration/
```

### Basic Test Run (Legacy)
```bash
# Run all integration tests
QNAP_HOST=your-qnap.local go test -v -integration \
    -nas-host=your-qnap.local \
    -nas-user=admin \
    -nas-key=~/.ssh/id_rsa \
    ./tests/integration/...

# Run specific legacy test suite
QNAP_HOST=your-qnap.local go test -v -integration \
    -nas-host=your-qnap.local \
    -run TestBasicDeployment \
    ./tests/integration/
```

### Advanced Test Options
```bash
# Run with custom timeout
go test -v -integration -timeout=10m \
    -nas-host=your-qnap.local \
    ./tests/integration/...

# Run performance benchmarks
go test -v -integration -bench=. \
    -nas-host=your-qnap.local \
    ./tests/integration/...

# Skip cleanup for debugging
go test -v -integration \
    -nas-host=your-qnap.local \
    -cleanup=false \
    ./tests/integration/...

# Parallel execution (experimental)
go test -v -integration -parallel=4 \
    -nas-host=your-qnap.local \
    ./tests/integration/...
```

### Makefile Targets
```bash
# Run integration tests (requires NAS_HOST environment variable)
NAS_HOST=your-qnap.local make integration-test

# Run with full cleanup and reporting
NAS_HOST=your-qnap.local make integration-test-full

# Generate integration test report
NAS_HOST=your-qnap.local make integration-test-report
```

## 📊 **Test Scenarios**

### 1. Basic Deployment Tests
- Single container deployment (nginx, postgres, redis)
- Port mapping validation
- Volume mounting verification (CACHEDEV/ZFS)
- Environment variable passing
- Restart policy enforcement

### 2. Compose Deployment Tests
- Multi-service application deployment
- Service dependency management
- Named volume creation and mounting
- Custom network configuration
- Environment variable substitution

### 3. Lifecycle Management Tests
- Container start/stop operations
- Graceful shutdown handling
- Container removal with cleanup
- Force removal scenarios
- Container recreation and updates

### 4. Volume and Storage Tests
- Host path mounting (`/share/CACHEDEV*_DATA/`, `/share/ZFS*_DATA/`)
- Named volume creation and management
- Volume permission validation
- Data persistence verification
- Storage cleanup operations

### 5. Network Connectivity Tests
- Container-to-container communication
- External network access validation
- Port binding and forwarding
- Custom bridge network creation
- DNS resolution within containers

### 6. Error Handling Tests
- Invalid image names
- Non-existent volume paths
- Port conflicts
- Resource exhaustion scenarios
- Network connectivity failures

### 7. QNAP-Specific Tests
- Container Station dynamic detection
- Multi-volume storage pool testing
- ZFS volume compatibility
- CACHEDEV path resolution
- Volume migration scenarios

## 🔧 **Test Configuration**

### Environment Variables
```bash
# Required
export QNAP_HOST="192.168.1.100"          # QNAP NAS IP address
export QNAP_USER="admin"                   # SSH username
export NAS_HOST="192.168.1.100"            # Legacy compatibility
export NAS_USER="admin"                    # Legacy compatibility
export NAS_SSH_KEY="~/.ssh/id_rsa"         # SSH private key path

# Optional
export NAS_PORT="22"                       # SSH port (default: 22)
export TEST_VOLUME_PATH="/share/CACHEDEV1_DATA/test"  # Test volume directory
export TEST_TIMEOUT="5m"                   # Test timeout
export TEST_PARALLEL="false"              # Parallel execution
export TEST_CLEANUP="true"                # Cleanup after tests
```

### Configuration File (test_config.yaml)
```yaml
nas:
  host: "192.168.1.100"
  user: "admin"
  ssh_key_path: "~/.ssh/id_rsa"

test:
  volume_path: "/share/CACHEDEV1_DATA/qnap-docker-test"
  cleanup: true
  timeout: "5m"

scenarios:
  basic_deployment: true
  compose_deployment: true
  performance: false  # Enable for performance testing
```

## 📈 **Test Output and Reporting**

### Console Output
```bash
=== RUN   TestIntegration
=== RUN   TestIntegration/BasicDeployment
    Deploying container: test-nginx-a8b9c2d3
    Container deployed successfully: c4f1e2d9a8b7
    Testing HTTP connectivity to http://192.168.1.100:8080
    ✓ HTTP endpoint accessible and returning expected content
=== RUN   TestIntegration/ComposeDeployment
    Deploying compose project: test-stack-x7y9z1
    Services deployed: web, api, db, cache
    ✓ All services started successfully
    ✓ Inter-service connectivity verified
--- PASS: TestIntegration (45.67s)
```

### Test Reports
```bash
# Generate detailed test report
go test -v -integration -json ./tests/integration/... > test_report.json

# Generate coverage report
go test -integration -coverprofile=integration_coverage.out ./tests/integration/...
go tool cover -html=integration_coverage.out -o integration_coverage.html
```

## 🧹 **Cleanup and Troubleshooting**

### Manual Cleanup
```bash
# Remove all test containers
ssh your-qnap-ip "find /share -name docker | grep container-station | head -1 | xargs -I {} {} ps -a | grep 'test-' | awk '{print \$1}' | xargs {} rm -f"

# Remove test volumes
ssh your-qnap-ip "find /share -name docker | grep container-station | head -1 | xargs -I {} {} volume ls | grep 'test-' | awk '{print \$2}' | xargs {} volume rm"

# Clean test directory
ssh your-qnap-ip "rm -rf /share/*/qnap-docker-test/*"
```

### Common Issues

**SSH Connection Failed:**
```bash
# Verify SSH access
ssh -v admin@your-qnap-ip

# Check SSH key permissions
chmod 600 ~/.ssh/id_rsa
chmod 644 ~/.ssh/id_rsa.pub
```

**Container Station Not Found:**
```bash
# Verify Container Station installation
ssh your-qnap-ip "find /share -name docker -type f | grep container-station"

# Check Container Station service
ssh your-qnap-ip "test -x /share/ZFS530_DATA/.qpkg/container-station/usr/bin/.libs/docker && echo 'Found' || echo 'Not found'"
```

**Volume Mount Errors:**
```bash
# Check volume path permissions
ssh your-qnap-ip "ls -la /share/CACHEDEV1_DATA/qnap-docker-test"
ssh your-qnap-ip "mkdir -p /share/CACHEDEV1_DATA/qnap-docker-test && chmod 755 /share/CACHEDEV1_DATA/qnap-docker-test"
```

## 🎯 **Best Practices**

### Test Development
1. **Idempotent Tests**: Each test should clean up after itself
2. **Unique Names**: Use random suffixes for test resources
3. **Timeout Handling**: Set appropriate timeouts for operations
4. **Error Validation**: Test both success and failure scenarios
5. **Resource Limits**: Be mindful of QNAP NAS resource constraints

### CI/CD Integration
1. **Environment Variables**: Use env vars for CI configuration
2. **Test Reports**: Generate JSON/XML reports for CI systems
3. **Parallel Execution**: Use `-parallel` flag judiciously
4. **Cleanup Verification**: Ensure all resources are cleaned up
5. **Test Isolation**: Each test run should be independent

This integration test suite provides comprehensive validation of qnap-docker functionality against real QNAP hardware, ensuring production readiness and reliability.