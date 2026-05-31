# mdt - MarkDown Task Manager

![GitHub release](https://img.shields.io/github/v/release/soe1hom-arch/mdt)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

A CLI task manager that stores tasks in a simple markdown file (`TASKS.md`).
Minimal, portable, and works with any project.

## Features

- **Simple** — plain markdown, human-readable
- **Zero dependencies** — single Go binary
- **Colorful terminal** — automatic color output
- **Categories** — organize tasks with `-c <category>`
- **Stats** — see progress, breakdown by category

## Installation

### Download binary

Grab the latest binary from the [Releases page](https://github.com/soe1hom-arch/mdt/releases):

```bash
# Example for Linux ARM64
wget https://github.com/soe1hom-arch/mdt/releases/download/v1.0.0/mdt-linux-arm64.tar.gz
tar xzf mdt-linux-arm64.tar.gz
sudo mv mdt /usr/local/bin/
```

### Or build from source

```bash
git clone https://github.com/soe1hom-arch/mdt.git
cd mdt
go build -o mdt .
sudo mv mdt /usr/local/bin/
```

## Usage

```
mdt              List all tasks
mdt add <desc>   Add a new task
mdt done <n>     Mark task as done
mdt undo <n>     Mark task as not done
mdt rm <n>       Remove a task
mdt stats        Show statistics
mdt init         Create TASKS.md template
```

## Example

```bash
mdt init
mdt add "Write README"
mdt add -c Go "Refactor parser"
mdt done 1
mdt stats
```

### Sample Output

```
 General
────────────────────────────────────────
   1. ✓ Write README
   2. ✗ Refactor parser

 📊 Task Statistics
──────────────────────────────
   Total:       2
   Done:        1
   Progress:    ██████░░░░░░░░░░░░░░  50%
```

## Why TASKS.md?

- Tracked in git alongside your code
- Editable in any text editor
- Renderable on GitHub
- No database, no config files

## Credits

Created and maintained by **[soe1hom-arch](https://github.com/soe1hom-arch)**.

## License

MIT &mdash; &copy; 2026 soe1hom-arch
