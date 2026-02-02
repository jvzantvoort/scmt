# Template Examples for SCMT Write Command

This directory contains example templates demonstrating the capabilities of the `scmt write` command.

## Available Templates

### 1. `server.conf` - Server Configuration File
Generates a configuration file in INI format with server details, roles, and conditional sections.

**Features:**
- Server metadata (type, owner, environment)
- Location and compute configuration
- Role-based conditional sections
- Role listing and counting

**Usage:**
```bash
scmt write examples/templates/server.conf server.conf
```

### 2. `setup.sh` - Server Setup Script
Generates a bash script for server provisioning based on assigned roles.

**Features:**
- Environment setup commands
- Role-based package installation
- Service configuration
- Dynamic script generation

**Usage:**
```bash
scmt write examples/templates/setup.sh setup.sh
chmod +x setup.sh
```

### 3. `kubernetes.yaml` - Kubernetes Deployment
Generates Kubernetes deployment manifests with role-based configuration.

**Features:**
- Dynamic replica count based on environment
- Port mapping based on roles
- Environment variables
- Conditional service creation

**Usage:**
```bash
scmt write examples/templates/kubernetes.yaml k8s-deployment.yaml
kubectl apply -f k8s-deployment.yaml
```

### 4. `config.json` - JSON Configuration
Generates a JSON configuration file with complete server metadata.

**Features:**
- Structured configuration export
- Boolean feature flags based on roles
- Metadata inclusion (timestamp, engineer)
- Role counting and listing

**Usage:**
```bash
scmt write examples/templates/config.json app-config.json
```

## Template Syntax Reference

### Available Data Structure
```go
.Config     - map[string]string  // Configuration parameters
.Roles      - []string           // Server roles
.Timestamp  - string             // Generation timestamp  
.Engineer   - string             // Current engineer name
.HasRole    - func(string) bool  // Check if role exists
```

### Template Functions
```go
join        - strings.Join
upper       - strings.ToUpper
lower       - strings.ToLower
title       - strings.Title
hasPrefix   - strings.HasPrefix
hasSuffix   - strings.HasSuffix
replace     - strings.Replace
split       - strings.Split
trim        - strings.TrimSpace
```

### Common Patterns

#### Conditional Content
```go
{{if .Config.ENVIRONMENT}}
environment = "{{.Config.ENVIRONMENT}}"
{{else}}
environment = "development"
{{end}}
```

#### Role-Based Configuration
```go
{{if .HasRole "web-server"}}
[web]
enabled = true
port = 8080
{{end}}
```

#### Iterating Over Roles
```go
{{range .Roles}}
- {{.}}
{{end}}
```

#### Role Array with Commas
```go
roles = [{{range $i, $role := .Roles}}{{if $i}}, {{end}}"{{$role}}"{{end}}]
```

## Creating Custom Templates

1. **Use Standard Go Template Syntax**: Follow Go's `text/template` syntax
2. **Access Configuration**: Use `{{.Config.KEY_NAME}}` for configuration values
3. **Check Roles**: Use `{{.HasRole "role-name"}}` for conditional content
4. **Add Metadata**: Include `{{.Timestamp}}` and `{{.Engineer}}` for audit trails
5. **Use Functions**: Leverage available template functions for string manipulation

## Testing Templates

Test your templates before deployment:

```bash
# Write to stdout for testing
scmt write my-template.conf

# Write to file
scmt write my-template.conf /tmp/output.conf
cat /tmp/output.conf
```