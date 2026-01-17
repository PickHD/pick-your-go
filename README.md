# Pick Your Go

A powerful CLI tool to quickly scaffold production-ready Go projects with various architectural patterns. Choose from Layered, Modular, or Hexagonal architecture templates and get started with your project in seconds.

## Features

- [x] **Multiple Architecture Patterns**
  - Layered Architecture (Traditional 3-tier)
  - Modular Architecture (DDD-based modular monolith)
  - Hexagonal Architecture (Ports and Adapters)

- [x] **Interactive CLI**
  - Beautiful terminal UI with `huh` and `lipgloss`
  - Step-by-step project initialization
  - Form-based configuration

- [x] **Template Management**
  - GitHub integration for template repositories
  - Local caching with 24-hour TTL
  - Easy template updates

- [x] **Production-Ready Projects**
  - Clean code structure
  - Best practices built-in
  - Ready for deployment
  - **Automatic import path updates** - All Go import paths are automatically updated to match your module path

## Installation

### From Source

```bash
git clone https://github.com/pickez/pick-your-go.git
cd pick-your-go
make install
```

### Using Go

```bash
go install github.com/pickez/pick-your-go/cmd/pick-your-go@latest
```

## Prerequisites

- Go 1.21 or higher
- Git
- GitHub token (for accessing private template repositories)

### Setting up GitHub Token

To access private template repositories, set the `PICK_YOUR_GO_GITHUB_TOKEN` environment variable:

```bash
export PICK_YOUR_GO_GITHUB_TOKEN="your_github_token_here"
```

**Note:** This token is required for all three architecture templates as they are hosted in private repositories.

## Usage

### Interactive Mode (Recommended)

The easiest way to create a new project:

```bash
pick-your-go init
```

This will launch an interactive form that guides you through:

1. Choosing your architecture pattern
2. Providing project details (name, module path, author, description)
3. Reviewing and confirming your choices

### Command-Line Mode

You can also provide all options via flags:

```bash
pick-your-go init \
  --architecture layered \
  --name my-awesome-app \
  --module github.com/username/my-awesome-app \
  --author "Your Name" \
  --description "A brief description of your project" \
  --output ./projects
```

### Available Commands

#### `init` - Create a new project

```bash
# Interactive mode
pick-your-go init

# With flags
pick-your-go init --architecture layered --name myapp --module github.com/user/myapp

# Options:
#   -a, --architecture string   Architecture type: layered, modular, or hexagonal
#   -n, --name string           Project name
#   -m, --module string         Go module path (e.g., github.com/user/project)
#   -o, --output string         Output directory (default: current directory)
#   -u, --author string         Author name
#   -d, --description string    Project description
```

#### `templates list` - List available templates

```bash
pick-your-go templates list
```

Shows all available architecture templates with their descriptions and cache status.

#### `templates update` - Update template cache

```bash
pick-your-go templates update
```

Force update the local template cache from remote repositories.

## Architecture Patterns

### Layered Architecture

Traditional three-tier architecture with clear separation:

**Best for:**

- Traditional web applications
- Teams familiar with layered architecture
- Projects with clear separation of concerns

### Modular Architecture

Domain-driven design approach with modular boundaries:

**Best for:**

- Large monolithic applications
- Teams using Domain-Driven Design
- Projects with clear domain boundaries

### Hexagonal Architecture

Ports and adapters pattern for maximum decoupling:

**Best for:**

- Complex business logic domains
- Projects requiring high testability
- Teams focused on clean architecture principles

## Development

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/pickez/pick-your-go.git
cd pick-your-go

# Install dependencies
make deps

# Run in development mode
make run
```

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the binary
make run               # Run the application
make install           # Install to $GOPATH/bin
make test              # Run tests
make lint              # Run linters
make fmt               # Format code
make clean             # Clean build artifacts
make templates-update  # Update template cache
make templates-list    # List available templates
```

### Project Structure

```
pick-your-go/
├── cmd/
│   └── pick-your-go/       # Application entry point
├── internal/
│   ├── cache/              # Template caching logic
│   ├── cli/                # CLI commands (cobra)
│   ├── cmd/                # Command implementations
│   ├── config/             # Configuration management
│   ├── generator/          # Architecture generators
│   └── template/           # Template management
├── pkg/
│   └── ui/                 # Interactive UI components
├── templates/              # Local template cache
├── Makefile                # Development commands
├── go.mod                  # Go module definition
└── README.md               # This file
```

## How It Works

1. **Template Selection**: Choose an architecture pattern (layered/modular/hexagonal)
2. **Configuration**: Provide project details via interactive form or flags
3. **Template Retrieval**: Tool fetches the template from GitHub private repository
4. **Caching**: Template is cached locally for 24 hours to speed up subsequent projects
5. **Generation**: Template is copied to your destination and customized with your project details
6. **Import Path Updates**: All Go import paths in `.go` files are automatically updated from the template's module name to your project's module path
7. **Ready to Use**: Your new project is ready to develop with correct import paths!

## Template Caching

Templates are cached in:

- **Linux**: `~/.cache/pick-your-go/`
- **macOS**: `~/Library/Caches/pick-your-go/`
- **Windows**: `%LocalAppData%\pick-your-go\cache\`

Cache TTL is 24 hours. Force update with:

```bash
pick-your-go templates update
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Huh](https://github.com/charmbracelet/huh) for interactive forms
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling

## Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

Made with by PickHD
