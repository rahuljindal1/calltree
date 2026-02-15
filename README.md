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
- 🔁 Re-run last analysis configuration easily

Currently supported language:

- **TypeScript** (via Tree-sitter)

---

## 🚀 Getting Started

### Run interactively

```bash
calltree analyze
```

If no arguments are provided, calltree starts in interactive mode and guides you through:

- selecting files or directories
- choosing output format
- configuring depth, focus, recursion, and more

Analyze a single file

```
calltree analyze src/app.js
```

Analyze a directory recursively

```
calltree analyze src -r
```

### 📤 Output Formats

Tree output (default)

```
initApp
├─ loadConfig
│  └─ parseEnv
└─ startServer
   └─ connectDB
```

JSON output

```
calltree analyze src/app.js --json
```

Example:

```
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
