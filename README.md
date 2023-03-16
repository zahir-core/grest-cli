# GREST CLI

The command line interface for GREST applications

## Installation

1. Make sure you have [Git](https://git-scm.com) and [Go](https://go.dev) installed with [GOPATH](https://pkg.go.dev/cmd/go#hdr-GOPATH_environment_variable) environment variable setted up.
2. Install the CLI globally by running
```bash
go install grest.dev/cmd/grest@latest
```
3. Check the version
```bash
grest version
```

## Usage

```bash
# Initialize new app in the current directory
grest init

# Add a new end point for the current app
grest add

# Format the struct tag
grest fmt

# Print the grest version
grest version

# Help about any command
grest help
```

## License

GREST officialy created, used, and maintained by [Zahir](https://zahiraccounting.com) Core Team. GREST is free and open-source software licensed under the [MIT License](https://github.com/zahir-core/grest-cli/blob/main/LICENSE).