# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Flufu

Flufu is a Go CLI wrapper that runs Claude Code through a PTY and overlays an animated ASCII pet (duck) on top of the terminal UI. It is also experimenting with local LLM inference via llama.cpp (through the `yzma` library) to potentially add on-device AI features.

## Commands

```bash
just run          # Run the app (auto-prepares lib/ if missing)
just prepare      # Install yzma and download llama.cpp dylibs into ./lib
just build        # Build binary to bin/flufu
just test         # Run tests and go vet
just lint         # Run golangci-lint
```

When running locally during development, set `YZMA_LIB=./lib` — `just run` does this automatically. The binary resolves `lib/` relative to the executable path at runtime.

## Architecture

The app runs two concurrent goroutines coordinated by `internal/orchestrator`:

1. **Agent** (`internal/agent/terminal.go`) — Spawns `claude` via PTY, wires stdin/stdout, handles terminal resize via `SIGWINCH`. When the agent exits, it cancels the shared context.
2. **Pet** (`internal/pet/render.go`) — Renders an animated ASCII duck overlay using ANSI escape sequences, anchored to the bottom-right corner. Clears on context cancel.

### Local LLM (`internal/ai/`)

Wraps the `yzma`/`llama.cpp` stack for on-device inference:
- `lama.go` — loads the shared llama libraries from `YZMA_LIB` or `<exe>/lib`
- `model.go` — downloads (on first run) and initializes `Qwen2.5-1.5B-Instruct-Q4_K_M` from HuggingFace; stores the model file via XDG config path (`~/.config/flufu/`)
- `generation.go` — token-by-token generation loop
- `prompt.go` — applies the model's chat template via `llama.ChatApplyTemplate`

The `chat` subcommand (`cmd/chat.go`) triggers this flow directly for testing. The design intent (noted in `model.go`) is for `ai.Init` to run in a goroutine keeping the model warm, with `ai.Respond(prompt)` used for inference.

### Config

`internal/config/load.go` resolves config/data paths: XDG first, falling back to `<exe>/`. Used for model file storage.

## Linter

`golangci-lint` is configured with a large strict ruleset (see `.golangci.yml`). Notable enforced rules: `err113` (no bare error comparisons), `wrapcheck` (all external errors must be wrapped), `mnd` (no magic numbers without named constants), `exhaustive` (exhaustive switches on enums).
