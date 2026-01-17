# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Terminal-Compatible ASCII Logo**: Replaced Unicode box-drawing characters with standard ASCII characters for maximum terminal compatibility
  - Logo now uses only standard ASCII characters: `+`, `-`, `|`
  - No more Unicode box-drawing characters (`╔═╗╠╣╚╝║`)
  - Works correctly on all terminal emulators, including those with limited Unicode support
  - Maintains the same visual style with "PICK YOUR GO" text in ASCII art
  - Properly centered and aligned in 80-column terminal width

- **Automatic Import Path Updates**: All Go import paths in `.go` files are now automatically updated when generating a new project. This ensures that internal imports match your project's module path instead of the template's module path.
  - When you specify `--module github.com/user/myproject`, all imports like `"go-layered-template/internal/config"` are automatically updated to `"github.com/user/myproject/internal/config"`
  - Supports all import formats:
    - Simple imports: `import "module/path"`
    - Aliased imports: `alias "module/path"`
    - Import blocks: Multi-line import statements
  - Works for all three architecture types: Layered, Modular, and Hexagonal

### Fixed

- Import paths in generated projects now correctly reflect the user's module path instead of the template's module path
- Previously, users had to manually update all import paths after project generation

### Technical Details

The implementation includes several helper functions in `/home/pickez/Workspaces/personal/pick-your-go/internal/generator/layered.go`:

- `extractOriginalModulePath()`: Extracts the original module name from go.mod before updating
- `updateImportPaths()`: Walks through all `.go` files and updates import paths
- `updateImportPathsInFile()`: Updates import paths in a single file
- `replaceImportPaths()`: Handles various import statement formats
- `replaceModulePathInLine()`: Replaces module paths in individual lines

### Safety Features

- Comprehensive validation to ensure old and new module paths are different
- Only processes `.go` files (skips go.mod, vendor, .git)
- Preserves external package imports (only updates internal project imports)
- Error handling that continues processing even if individual files fail

### Testing

Added comprehensive unit tests in `/home/pickez/Workspaces/personal/pick-your-go/internal/generator/import_paths_test.go`:

- Test module path extraction
- Test import path replacement in various formats
- Test file-level import path updates
- Test edge cases (external packages, standard library, aliased imports)

All tests pass successfully.

## [Previous Versions]

See git history for earlier changes.
