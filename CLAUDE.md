# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

**Local Development:**
- `make build` - Clean and build the binary
- `make run` - Build and run the application locally  
- `make clean` - Remove all build artifacts from bin/

**Docker Development:**
- `docker compose up --build --watch` - Run with hot reload for development
- Use profile `dev` for development, `prd` for production

## Architecture Overview

This is a multiplayer 3D raycaster "game" built with Go, HTMX, and WebSockets for real-time multiplayer functionality.

**Core Architecture:**
- **Frontend**: Server-side rendered HTML templates using Go's standard `html/template` package, styled with CSS transforms for 3D projection
- **Backend**: Go web server using Fiber framework with WebSocket support via `hxsocketsfiber`
- **Real-time Communication**: WebSockets for player synchronization and game state updates
- **Game Loop**: Server-side player sync runs every 250ms, status reporting every 750ms

**Key Components:**
- `main.go`: Entry point, Fiber app setup, WebSocket server, and concurrent game loops
- `gameobjects/`: Core game entities (Player, GameMap, Wall) with physics and collision detection
- `views/`: HTML templates for rendering 3D scenes, HUD, player controls, and sync data
- `network/`: WebSocket message handlers for player input and state synchronization
- `utils/`: Math utilities (Vector3, Matrix4, 2D dimensions)
- `maps/`: PNG files used as level data (colors mapped to wall types)
- `prefabs/`: HTML/CSS prefab definitions for game objects

**Game Map System:**
- Maps are loaded from PNG files in `maps/` directory (currently uses "smile.png")  
- Color values in PNG are mapped to wall types (white=industrial, yellow=brick, red=copper, blue=htmxcon)
- Map parsing includes padding to prevent edge glitches
- Walls are generated based on neighboring cells in the color map

**Player System:**
- Physics-based movement with velocity, acceleration, and friction
- Collision detection using box colliders
- Jumping mechanics with gravity simulation
- Real-time synchronization across all connected clients
- Player state includes position, rotation, controls, and statistics

**Rendering Pipeline:**
- Server-side template rendering generates HTML/CSS for 3D projection
- CSS transforms handle 3D positioning and perspective
- HTMX enables dynamic updates without page refreshes
- WebSockets push real-time player sync data to all clients

**Development Notes:**
- Templates are loaded automatically from `views/*.html` at application startup
- Server runs on port 3000 with static assets served from `public/`
- Game map can be changed by modifying the LoadMap call in `gameobjects/gameMap.go:174`
- Player sync rate can be adjusted in the main game loop goroutine