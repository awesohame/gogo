# GoGo: Go (Board Game) in Go(lang)

GoGo is a modular Go game engine written in Go, designed for easy integration into your own projects. It provides a session layer for two-player matches, undo/redo, passing, resigning, and scoring, with a clean public API for developers.

## Features
- Pure Go implementation of the game engine
- Session management for two-player games
- Undo/redo functionality
- Pass and resign support
- Accurate scoring according to chinese rules (territory, captures, komi)
- Simple API for integration

## Getting Started

To use GoGo as a module in your Go project:

1. Add GoGo to your dependencies:

```
go get github.com/awesohame/gogo
```

2. Import the engine in your code:

```
import "github.com/awesohame/gogo/pkg/engine"
```

## CLI Demo

Run a simple two-player match in your terminal:

```
go run ./cmd/dev/main.go
```

## Project Structure
- `pkg/engine/game.go`: Public API for game management
- `internal/game/session.go`: Session logic
- `internal/engine/`: Core engine (board, moves, scoring)
- `cmd/app/`: Main application entry
- `cmd/dev/`: CLI demo
- `tests/`: Integration tests

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
