# Byteload Metrics Agent

A lightweight system metrics collection agent written in Go that exposes system metrics via HTTP endpoints.

## Features

- System metrics collection (CPU, Memory, Disk, Network)
- Service status monitoring
- Docker container metrics
- Basic authentication support
- Systemd service integration
- Multi-platform support (Linux, macOS)

## Installation

### Using Package Managers

#### DEB-based systems (Debian/Ubuntu):
```bash
sudo dpkg -i byteload-agent_*.deb
```

#### RPM-based systems (RHEL/CentOS):
```bash
sudo rpm -i byteload-agent_*.rpm
```

### Manual Installation

1. Download the latest release for your platform
2. Extract the archive:
```bash
tar xzf byteload-agent_*.tar.gz
```
3. Copy the binary:
```bash
sudo cp byteload-agent /usr/local/bin/
```
4. Copy configuration:
```bash
sudo mkdir -p /etc/byteload
sudo cp configs/byteload.yaml /etc/byteload/
sudo cp deploy/env.sample /etc/byteload/env
```
5. Install service:
```bash
sudo cp deploy/byteload.service /lib/systemd/system/
sudo systemctl daemon-reload
```

## Configuration

### Main Configuration File (`/etc/byteload/byteload.yaml`)

```yaml
server:
  port: "9001"
  host: "0.0.0.0"

logging:
  level: "info"
  format: "json"

metrics:
  enabled: true
  endpoint: "/metrics"

security:
  basic_auth:
    enabled: true
    username: "metrics"
    password: "change-me-in-production"
```

### Environment Variables

- `BYTELOAD_CONFIG_FILE`: Path to configuration file (default: `/etc/byteload/byteload.yaml`)

## API Endpoints

All endpoints require basic authentication if enabled in configuration.

- `/system` - Complete system information
- `/os` - Operating system information
- `/cpu` - CPU metrics
- `/memory` - Memory usage
- `/storage` - Disk usage and IO statistics
- `/services` - Service status (with optional filtering)
- `/docker` - Docker container metrics

### Example Usage

```bash
# Get CPU metrics
curl -u metrics:your-password http://localhost:9001/cpu

# Get specific services status
curl -u metrics:your-password 'http://localhost:9001/services?services=["nginx","postgresql"]'
```

## Service Management

```bash
# Start the service
sudo systemctl start byteload.service

# Enable on boot
sudo systemctl enable byteload.service

# Check status
sudo systemctl status byteload.service

# View logs
sudo journalctl -u byteload.service
```

## Security

- The service runs as a dedicated `byteload` user
- Basic authentication support
- Systemd service hardening
- No root privileges required for metric collection

## Building from Source

Requirements:
- Go 1.21 or later
- Make (optional)

```bash
# Clone repository
git clone https://github.com/yourusername/metrics-agent-go
cd metrics-agent-go

# Build
go build -o byteload-agent ./cmd/byteload-agent

# Or use goreleaser for all platforms
goreleaser release --snapshot --clean
```

## License

MIT License - see LICENSE file for details
