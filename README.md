# SCMT - Server Configuration Management Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/jvzantvoort/scmt)](https://goreportcard.com/report/github.com/jvzantvoort/scmt)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0-brightgreen.svg)](https://github.com/jvzantvoort/scmt/releases)

SCMT (Server Configuration Management Tool) is a command-line utility for managing server configuration parameters and roles with full audit logging and change tracking. It provides a simple yet powerful interface for maintaining server metadata, configuration values, and role assignments.

## 🚀 Features

- **Configuration Management**: Set, get, and track server configuration parameters
- **Role Management**: Add, remove, and list server roles  
- **Template Processing**: Generate configuration files from templates using server data
- **Change Auditing**: Full audit trail with engineer attribution and timestamps
- **Multiple Output Formats**: Table and JSON output support
- **Persistent Storage**: JSON-based configuration storage with backup support
- **Logging**: Comprehensive change logging for compliance and debugging
- **CLI Interface**: Intuitive command-line interface with comprehensive help

## 📦 Installation

### Prerequisites

- Go 1.22.6 or later
- Git (for version information)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/jvzantvoort/scmt.git
cd scmt

# Build the application
go build -o scmt cmd/scmt/*.go

# Or use the build script for cross-compilation
./build.sh build
```

### Install Locally

```bash
# Install to $GOPATH/bin
./build.sh install
```

## 🎯 Quick Start

### 1. Initialize Configuration

```bash
# Initialize with default settings
scmt init

# Initialize with custom engineer name
scmt -E "John Doe" init
```

### 2. View Current Configuration

```bash
# Table format (default)
scmt dump

# JSON format
scmt -J dump
```

### 3. Set Configuration Values

```bash
# Set a configuration parameter
scmt set OWNER "DevOps Team"
scmt set ENVIRONMENT "production"
```

### 4. Manage Server Roles

```bash
# Add roles
scmt role add web-server
scmt role add database
scmt role add load-balancer

# List roles
scmt role list

# Remove a role
scmt role remove load-balancer
```

### 5. View Change History

```bash
# View log for a specific parameter
scmt log OWNER
```

## 📚 Command Reference

### Global Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--configdir` | `-C` | Directory for config files | `/etc/scmt` |
| `--engineer` | `-E` | Engineer name for changes | Current user |
| `--json` | `-J` | JSON output format | `false` |
| `--logfile` | `-L` | Specify logfile | `/var/log/scmt.log` |
| `--loglevel` | `-l` | Log level (debug/info/warn/error) | `info` |
| `--message` | `-M` | Message for changes | Empty |

### Commands

#### `scmt init`
Initialize server configuration with default values.

```bash
scmt init [flags]

Flags:
  -d, --description string   Description of the installation
  -t, --type string          Type of installation (default "default")
```

**Default Values:**
- `TYPE`: server
- `OWNER`: Mad House
- `COUNTRY_CODE`: NL
- `REGION_CODE`: EU
- `TIMEZONE`: Europe/Amsterdam
- `COMPUTE_ZONE`: europe-west4-a
- `COMPUTE_REGION`: europe-west3

#### `scmt set <key> <value>`
Set a configuration parameter.

```bash
# Examples
scmt set ENVIRONMENT production
scmt set OWNER "DevOps Team"
scmt -M "Monthly update" set VERSION "2.1.0"
```

#### `scmt dump`
Display all configuration parameters.

```bash
# Table format
scmt dump

# JSON format
scmt -J dump
```

#### `scmt role <command>`
Manage server roles.

```bash
# Add a role
scmt role add <role>

# Remove a role  
scmt role remove <role>

# List all roles
scmt role list

# Examples
scmt role add web-server
scmt role add database
scmt role list
scmt role remove database
```

#### `scmt log <parameter>`
View change history for a parameter.

```bash
scmt log OWNER
scmt log TYPE
```

#### `scmt version`
Display version information.

```bash
scmt version
```

#### `scmt write <template> [output]`
Process template files using server configuration data.

```bash
# Write to stdout
scmt write template.conf

# Write to file
scmt write template.conf /etc/myapp/config.conf

# Process Kubernetes template
scmt write k8s-template.yaml deployment.yaml
```

**Template Features:**
- Access all configuration parameters via `{{.Config.KEY_NAME}}`
- Check roles with `{{.HasRole "role-name"}}`
- Iterate over roles with `{{range .Roles}}`
- Include metadata with `{{.Timestamp}}` and `{{.Engineer}}`
- Built-in string functions (join, upper, lower, etc.)

#### `scmt completion <shell>`
Generate shell completion scripts.

```bash
# Generate bash completion
scmt completion bash > /etc/bash_completion.d/scmt

# Generate zsh completion  
scmt completion zsh > "${fpath[1]}/_scmt"

# Generate fish completion
scmt completion fish > ~/.config/fish/completions/scmt.fish
```

## 🏗️ Configuration Structure

SCMT stores configuration in JSON format with the following structure:

```json
{
  "elements": [
    {
      "option": "OWNER",
      "value": {
        "value": "DevOps Team",
        "engineer": "jvzantvoort", 
        "message": "Updated team ownership",
        "changed": "2026-02-02T13:22:37.644131369Z"
      }
    }
  ],
  "roles": [
    "web-server",
    "database"
  ]
}
```

### Data Elements

Each configuration parameter includes:
- **option**: Parameter name
- **value**: Current value
- **engineer**: Who made the change
- **message**: Change description
- **changed**: Timestamp of change

### Roles

Simple string array containing assigned server roles.

## 🗂️ File Locations

### Default Paths

| File | Default Location | Purpose |
|------|------------------|---------|
| Configuration | `/etc/scmt/data.json` | Main configuration storage |
| Log File | `/var/log/scmt.log` | Change audit log |
| Config File | `~/.scmt.yaml` | User configuration (optional) |

### Custom Paths

Override default paths using flags:

```bash
# Use custom config directory
scmt -C /opt/myapp/config init

# Use custom log file
scmt -L /var/log/myapp-scmt.log set OWNER "MyApp Team"
```

## 🔧 Development

### Project Structure

```
scmt/
├── cmd/scmt/           # CLI commands and main application
├── config/             # Configuration management
├── data/               # Data models and persistence
├── logger/             # Audit logging functionality
├── messages/           # Help text and UI messages
├── utils/              # Utility functions
├── build.sh           # Build and development script
├── go.mod             # Go module definition
└── README.md          # This file
```

### Build Commands

```bash
# Format code
./build.sh fmt

# Run tests with coverage
go test -cover ./...

# Static analysis
./build.sh check

# Build for multiple platforms
./build.sh package

# Clean build artifacts
./build.sh cleanup
```

### Dependencies

Key dependencies:
- **cobra**: CLI framework
- **viper**: Configuration management
- **logrus**: Structured logging
- **tablewriter**: Table output formatting

## 🔍 Examples

### Infrastructure Management

```bash
# Set up a web server
scmt init
scmt set ENVIRONMENT production
scmt set SERVICE_NAME web-api
scmt role add web-server
scmt role add monitoring-target
scmt role add backup-client

# View configuration
scmt dump
```

### JSON Integration

```bash
# Export configuration for automation
scmt -J dump > server-config.json

# Role information for scripts
scmt -J role list | jq '.roles[]'
```

### Change Tracking

```bash
# Make documented changes
scmt -E "Alice Smith" -M "Quarterly security update" set PATCH_LEVEL "2024.Q1"

# Track who changed what
scmt log PATCH_LEVEL
```

### Template Processing

```bash
# Create a server configuration template
cat > server.template << 'EOF'
# Server: {{.Config.TYPE}}
# Owner: {{.Config.OWNER}} 
# Generated: {{.Timestamp}} by {{.Engineer}}

[server]
type = "{{.Config.TYPE}}"
environment = "{{if .Config.ENVIRONMENT}}{{.Config.ENVIRONMENT}}{{else}}development{{end}}"

[roles]
{{range .Roles}}
{{.}} = true
{{end}}

{{if .HasRole "web-server"}}
[web]
port = 8080
enabled = true
{{end}}
EOF

# Process template with current server data
scmt write server.template /etc/myapp/server.conf
```

### Template Processing

```bash
# Create a server configuration template
cat > server.template << 'EOF'
# Server: {{.Config.TYPE}}
# Owner: {{.Config.OWNER}} 
# Generated: {{.Timestamp}} by {{.Engineer}}

[server]
type = "{{.Config.TYPE}}"
environment = "{{if .Config.ENVIRONMENT}}{{.Config.ENVIRONMENT}}{{else}}development{{end}}"

[roles]
{{range .Roles}}
{{.}} = true
{{end}}

{{if .HasRole "web-server"}}
[web]
port = 8080
enabled = true
{{end}}
EOF

# Process template with current server data
scmt write server.template /etc/myapp/server.conf
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Run linting: `./build.sh check`
6. Submit a pull request

## 📋 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**John van Zantvoort** - [john@vanzantvoort.org](mailto:john@vanzantvoort.org)

---

*SCMT helps maintain clear, auditable server configuration with role-based organization and comprehensive change tracking.*
