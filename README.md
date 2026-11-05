# guescape

A data-driven text-based escape room engine.

guescape runs text adventures defined entirely in HCL world files. You point the
engine at a world file and play it from the terminal.

## Playing

New to playing a guescape game? Start with the player guide:

- [Playing a Game](docs/player/playing.md) explains how to launch a game and the basics of typing commands.

To try the bundled example world:

```
go run . -path ./docs/example.hcl
```

## Developer Quick Start

### Prerequisites

- Go 1.27 or newer

### Clone and run

```
git clone https://github.com/en0/guescape.git
cd guescape
go run . -path ./docs/example.hcl
```

The `-path` flag is required and points at the HCL world file to load.

### Build

```
go build -o guescape .
./guescape -path ./docs/example.hcl
```

Cross-compiled release binaries (Linux and Windows, amd64):

```
make build-release
```

### Test

```
make tests
```

Re-run tests automatically on file changes (requires `inotifywait`):

```
make watch-tests
```

## Layout

| Path | Purpose |
|------|---------|
| `main.go` | Entry point; loads a world file and starts the REPL. |
| `internal/game` | Core domain: world state, effects, rules, resolution. |
| `internal/hcl` | Loads and decodes HCL world files into the domain. |
| `internal/lua` | Compiles Lua scripts used by interactive modes. |
| `internal/cli` | The terminal REPL and command parsing. |
| `docs/` | Player guide and design notes. |
| `docs/example.hcl` | A runnable example world. |
