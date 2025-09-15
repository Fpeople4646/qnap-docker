package qnap

const (
	// DockerBinary is the default path to the Docker binary in Container Station (will be auto-detected)
	DockerBinary = "/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker"
	// DockerCompose is the default path to the Docker Compose binary
	DockerCompose = "/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker-compose"
	// ContainerStationBinPath is the default Container Station binary directory
	ContainerStationBinPath = "/share/CACHEDEV1_DATA/.qpkg/container-station/bin"
	// ContainerStationSbinPath is the default Container Station sbin directory
	ContainerStationSbinPath = "/share/CACHEDEV1_DATA/.qpkg/container-station/sbin"
	// SocketPath is the path to the Docker socket
	SocketPath = "/var/run/docker.sock"
	// DefaultVolume is the default volume path on first CACHEDEV
	DefaultVolume = "/share/CACHEDEV1_DATA/docker"

	// DefaultSSHPort is the default SSH port for QNAP
	DefaultSSHPort = 22
	// DefaultSSHUser is the default SSH username
	DefaultSSHUser = "admin"

	// DefaultRestartPolicy is the default container restart policy
	DefaultRestartPolicy = "unless-stopped"
	// DefaultNetwork is the default Docker network mode
	DefaultNetwork = "bridge"

	// CacheDevPrefix is the prefix for QNAP cache device paths
	CacheDevPrefix = "/share/CACHEDEV"
	// CacheDevSuffix is the suffix for QNAP cache device paths
	CacheDevSuffix = "_DATA"
	// SharePrefix is the prefix for QNAP share paths
	SharePrefix = "/share"
)
