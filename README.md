# calltree

**calltree** is a fast, language-aware CLI tool that analyzes source code and generates a readable call hierarchy of functions and methods.

It helps you:

- understand execution flow in complex codebases
- explore unfamiliar projects quickly
- debug deeply nested logic
- reason about architecture and dependencies

The tool is designed with a clean separation between a reusable core and language-specific analysis, making it easy to extend to new languages over time.

---

## ✨ Features

- 📊 Interactive CLI mode (guided analysis, reruns supported)
- 🌳 Human-readable call tree output
- 🎯 Focus analysis on a specific function
- 📏 Limit depth to control large outputs
- 📁 Recursive directory scanning
- 🚫 Exclude noisy directories (e.g. `node_modules`)
- 🧠 Language-aware filtering of built-in calls
- 📄 Optional JSON output for tooling and automation
- 📄 Write JSON output to a file
- 📍 Show source file names per function
- 🔁 Re-run last analysis configuration easily (`--rerun`)

**Currently supported language:**

- **TypeScript / TSX** (via Tree-sitter)

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.22+**
- **CGO** enabled (required for Tree-sitter)

### Build & Install

```bash
# Build
go build -o calltree ./cmd/calltree

# Or install to $GOPATH/bin
go install ./cmd/calltree
```

### Run interactively

```bash
calltree analyze
```

If no arguments are provided, calltree starts in interactive mode and guides you through:

- selecting files or directories
- choosing output format
- configuring depth, focus, recursion, extensions, and more

### Analyze a single file

```bash
calltree analyze src/app.ts
```

### Analyze a directory recursively

```bash
calltree analyze src -r --ext .ts --ext .tsx
```

For recursive scanning, use `--ext` to specify file extensions (e.g. `.ts`, `.tsx`). Exclude directories with `--exclude-dir`:

```bash
calltree analyze . -r --ext .ts --ext .tsx --exclude-dir node_modules
```

### Re-run last analysis

```bash
calltree analyze --rerun
```

---

## 📤 Output Formats

### Tree output (default)

```
initApp
├─ loadConfig
│  └─ parseEnv
└─ startServer
   └─ connectDB
```

### JSON output

```bash
calltree analyze src/app.ts --json
```

Example:

```json
[
  {
    "name": "initApp",
    "children": [
      { "name": "loadConfig" },
      { "name": "startServer" }
    ]
  }
]
```

### Write JSON to file

```bash
calltree analyze src/app.ts --json --json-file output.json
```

---

## 🐳 Docker

### Development

```bash
docker compose -f Docker/docker-compose.yml up -d calltree
docker compose -f Docker/docker-compose.yml exec calltree bash
```

### Production build

```bash
docker compose -f Docker/docker-compose.yml build calltree-exec
```

---

## 🛠 Development

The project includes a **Dev Container** (`.devcontainer/devcontainer.json`) for VSCode/Cursor. Open the folder in a dev container to get a preconfigured Go environment with Tree-sitter and Delve debugger.

---

## 📚 CLI Reference

| Flag | Description |
|------|-------------|
| `--depth`, `-d` | Maximum call depth |
| `--roots-only` | Show only entry-point functions |
| `--json` | Output call tree as JSON |
| `--json-file` | Write JSON output to file |
| `--focus` | Focus on a specific function |
| `--recursive`, `-r` | Scan directories recursively |
| `--exclude-dir` | Directories to exclude (repeatable) |
| `--ext` | File extensions to include (repeatable) |
| `--rerun` | Re-run last analysis configuration |
| `--show-file` | Show source file name per function |
| `--include-builtins` | Include built-in calls (map, includes, etc.) |

For a detailed flags reference, see [docs/cli-flags.md](docs/cli-flags.md).
