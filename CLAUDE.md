# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Flufu

Flufu is a Go CLI wrapper that runs Claude Code through a PTY and overlays an animated ASCII pet (duck) on top of the terminal UI. The duck periodically displays AI-generated jokes in a speech bubble using a local LLM (llama.cpp via the `yzma` library).

## Commands

```bash
just run          # Run the app (auto-prepares lib/ if missing)
just prepare      # Install yzma and download llama.cpp dylibs into ./lib
just build        # Build binary to bin/flufu
just test         # Run tests and go vet
just lint         # Run golangci-lint
```

Run a single test: `go test -v -run TestName ./internal/package/`

When running locally during development, set `YZMA_LIB=./lib` — `just run` does this automatically. The binary resolves `lib/` relative to the executable path at runtime.

## Architecture

Three concurrent goroutines coordinated by `internal/orchestrator/apps.go`, all sharing a cancellable context:

1. **AI** (`internal/ai/`) — Loads llama.cpp and the model in the background; closes the `ready` channel on `ai.Instance` when warm. Other goroutines gate on `<-aiInstance.Ready()`.
2. **Agent** (`internal/agent/terminal.go`) — Spawns `claude` via PTY, wires stdin/stdout, handles terminal resize via `SIGWINCH`. When the agent exits, it cancels the shared context, which shuts down the other goroutines.
3. **Pet** (`internal/pet/render.go`) — Renders an animated ASCII duck overlay using ANSI escape sequences anchored to the bottom-right corner. Once the AI is ready, fetches jokes via `aiInstance.Joke()` on a timer and draws them in a speech bubble (see `bubble.go`). Clears both duck and bubble on context cancel.

### Pet rendering details (`internal/pet/`)

- `pets.go` — Static ASCII frame definitions, keyed by pet name (currently only `"duck"`).
- `render.go` — Ticker-driven render loop: clears previous frame, advances animation, manages the joke fetch/display lifecycle.
- `bubble.go` — Draws/clears a word-wrapped speech bubble to the left of the pet using ANSI cursor positioning and DEC save/restore (`\0337`/`\0338`).

### Local LLM (`internal/ai/`)

- `lama.go` — Loads the shared llama libraries from `YZMA_LIB` or `<exe>/lib`.
- `model.go` — Downloads (on first run) and initializes `Qwen2.5-1.5B-Instruct-Q4_K_M` from HuggingFace; stores the model at `~/.config/flufu/` (XDG).
- `generation.go` — Token-by-token generation loop; `Joke()` uses a hardcoded duck-persona prompt.
- `prompt.go` — Builds chat-templated prompts via `llama.ChatApplyTemplate`.

The `chat` subcommand (`cmd/chat.go`) runs the full AI init+generate flow interactively for testing without launching Claude Code.

### Config

`internal/config/load.go` resolves config/data paths: XDG first, falling back to `<exe>/`. Used for model file storage.

## Linter

`golangci-lint` is configured with a large strict ruleset (see `.golangci.yml`). Notable enforced rules:
- `err113` — no bare `errors.New` in comparisons; define sentinel errors as package-level vars
- `wrapcheck` — all errors from external packages must be wrapped with `fmt.Errorf("...: %w", err)`
- `mnd` — no magic numbers; use named constants
- `exhaustive` — switches on enum types must be exhaustive
