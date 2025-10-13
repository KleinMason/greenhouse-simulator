# Greenhouse Simulator - Implementation Plan

## Project Overview
A tick-based greenhouse simulation system with web interface for monitoring and controlling plant growth, soil moisture, and automated/manual watering systems.

## Core Requirements
- **Simulation**: Tick-based (1 tick = 5 seconds real-time, configurable)
- **Plants**: 3-5 plants, different types with varying water needs
- **Plant Properties**: Soil saturation, health, growth stage, can die
- **Sensors**: Real-time readings, shared per plant-type section, exact readings
- **Watering**: Scheduled + manual override, gradual water application
- **Interface**: Web-based UI
- **Learning Focus**: Concurrency (goroutines/channels), Testing
- **Future-Ready**: Architecture supports temperature, light, humidity sensors

## Architecture Overview

### 1. Core Domain Models (`internal/models/`)

#### Plant Model
```go
type PlantType struct {
    Name                string
    OptimalSaturation   float64  // 0.0 to 1.0
    MinSaturation       float64  // below this, health degrades
    MaxSaturation       float64  // above this, health degrades
    GrowthRate          float64  // per tick
    SaturationDepletion float64  // per tick
}

type Plant struct {
    ID             string
    Type           PlantType
    SoilSaturation float64    // 0.0 to 1.0
    Health         float64    // 0.0 (dead) to 1.0 (perfect)
    GrowthStage    float64    // 0.0 (seed) to 1.0 (mature)
    Alive          bool
    CreatedAt      time.Time
}
```

#### Sensor Model
```go
type SensorType string

const (
    SoilMoisture SensorType = "soil_moisture"
    // Future: Temperature, Light, Humidity
)

type Sensor struct {
    ID        string
    Type      SensorType
    SectionID string  // which plant section it monitors
}

type SensorReading struct {
    SensorID  string
    Timestamp time.Time
    Value     float64
}
```

#### Watering System Model
```go
type WateringEvent struct {
    SectionID    string
    Amount       float64       // amount of water to apply
    StartTime    time.Time
    Duration     time.Duration // gradual application
    IsManual     bool
}

type WateringSchedule struct {
    SectionID        string
    TargetSaturation float64
    CheckInterval    int     // in ticks
    WaterAmount      float64
    Enabled          bool
}
```

### 2. Simulation Engine (`internal/engine/`)

#### Simulator
Central coordinator that manages the tick-based simulation.

Key responsibilities:
- Manage simulation clock (tick every 5 seconds by default)
- Coordinate updates across all systems
- Emit events for state changes

Concurrency pattern: Use a ticker and goroutine for the main simulation loop

Questions to explore:
- How will you ensure thread-safe updates to plant state?
- How will you coordinate multiple goroutines (sensor reads, watering, plant updates)?
- What happens if you pause/resume the simulation?

Suggested interface:
```go
type Simulator interface {
    Start() error
    Stop() error
    Pause()
    Resume()
    GetCurrentTick() int
    GetPlants() []Plant
}
```

#### Plant Updater
Updates plant state each tick.

Algorithm per tick:
1. Decrease soil saturation based on plant type depletion rate
2. Check if saturation is in optimal range
3. Update health (increase if optimal, decrease if not)
4. Update growth stage (only if health is above threshold)
5. Check if plant dies (health reaches 0)

Testing focus: Unit tests for plant state transitions

### 3. Sensor System (`internal/sensors/`)

#### Sensor Manager
Manages all sensors and provides real-time readings.

Key features:
- Group sensors by plant section
- Provide current readings on demand
- Calculate average saturation for a section

Concurrency challenge: Sensors read in real-time while simulation updates plant state in ticks. How do you handle concurrent reads/writes?

Hint: Consider using channels or mutexes. Research `sync.RWMutex`.

### 4. Watering System (`internal/watering/`)

#### Watering Controller
Manages both scheduled and manual watering.

Scheduled watering:
- Check sensor readings at intervals
- Trigger watering if below threshold
- Apply water gradually over time

Manual override:
- Accept watering commands from web interface
- Apply water immediately (but still gradually)

Concurrency pattern:
- Use goroutines for scheduled watering checks
- Use channels to receive manual watering commands
- Question: How do you prevent over-watering from simultaneous scheduled + manual events?

Gradual water application:
- Instead of instant saturation increase, spread it over multiple ticks
- Example: 0.3 water over 4 ticks = +0.075 per tick

### 5. Web Interface (`internal/web/`)

#### HTTP Server
Use `net/http` or a lightweight framework like `chi` or `echo`.

Endpoints:
- `GET /api/plants` - List all plants with current state
- `GET /api/sensors/readings` - Get current sensor readings
- `POST /api/watering/manual` - Trigger manual watering
- `GET /api/watering/schedules` - Get watering schedules
- `PUT /api/watering/schedules/:id` - Update schedule
- `GET /api/simulation/status` - Get tick count, running status
- `POST /api/simulation/control` - Pause/resume simulation

Frontend (optional for MVP):
- Simple HTML + JavaScript
- Display plant cards with health, saturation, growth
- Buttons for manual watering
- Real-time updates (polling or WebSockets)

### 6. Configuration (`internal/config/`)

Configuration file (YAML or JSON):
```yaml
simulation:
  tick_interval: 5s  # s seconds
  
plant_types:
  - name: Tomato
    optimal_saturation: 0.7
    min_saturation: 0.4
    max_saturation: 0.9
    growth_rate: 0.02
    saturation_depletion: 0.05
    
  - name: Lettuce
    optimal_saturation: 0.8
    min_saturation: 0.6
    max_saturation: 0.95
    growth_rate: 0.03
    saturation_depletion: 0.04

greenhouse:
  sections:
    - id: section-1
      plant_type: Tomato
      plant_count: 3
      
    - id: section-2
      plant_type: Lettuce
      plant_count: 2
```

## Project Structure

```
greenhouse-simulator/
├── main.go                          # Entry point
├── go.mod
├── go.sum
├── config.yaml                      # Configuration file
├── internal/
│   ├── models/                      # Domain models
│   │   ├── plant.go
│   │   ├── sensor.go
│   │   └── watering.go
│   ├── engine/                      # Simulation engine
│   │   ├── simulator.go
│   │   └── plant_updater.go
│   ├── sensors/                     # Sensor system
│   │   └── manager.go
│   ├── watering/                    # Watering system
│   │   └── controller.go
│   ├── config/                      # Configuration
│   │   └── config.go
│   └── web/                         # Web interface
│       ├── server.go
│       ├── handlers.go
│       └── static/                  # HTML, CSS, JS
├── tests/                           # Integration tests
│   └── simulation_test.go
└── README.md
```

## Implementation Phases

### Phase 1: Core Models & Plant Logic (Week 1)
Goal: Get plants updating correctly in a tick-based system

Tasks:
1. Define `PlantType` and `Plant` structs
2. Implement plant update logic (saturation depletion, health, growth)
3. Create basic simulator that ticks and updates plants
4. Write unit tests for plant state transitions

Key Go concepts: Structs, methods, `time.Ticker`, basic testing

Deliverable: A CLI program that prints plant states each tick

### Phase 2: Sensor System (Week 2)
Goal: Read plant sensor data in real-time

Tasks:
1. Define `Sensor` and `SensorReading` models
2. Implement sensor manager that groups sensors by section
3. Integrate with plant state (sensors read current saturation)
4. Handle concurrent reads with mutex or channels
5. Write tests for sensor readings

Key Go concepts: Interfaces, concurrency (`sync.RWMutex` or channels)

Challenge: Ensure sensor reads don't conflict with simulation updates

### Phase 3: Watering System (Week 2-3)
Goal: Implement scheduled and manual watering

Tasks:
1. Define watering models (`WateringEvent`, `WateringSchedule`)
2. Implement scheduled watering (runs in separate goroutine)
3. Implement manual watering via channel communication
4. Add gradual water application logic
5. Write tests for watering logic

Key Go concepts: Goroutines, channels, select statements

Challenge: Coordinate multiple watering events without conflicts

### Phase 4: Configuration & Initialization (Week 3)
Goal: Load greenhouse setup from config file

Tasks:
1. Define config structure
2. Read and parse YAML/JSON config
3. Initialize plants, sensors, schedules from config
4. Make tick interval configurable

Key Go concepts: File I/O, YAML/JSON unmarshaling

Package suggestion: `gopkg.in/yaml.v3` for YAML support

### Phase 5: Web Interface (Week 4)
Goal: Web UI for monitoring and control

Tasks:
1. Set up HTTP server
2. Implement API endpoints (plants, sensors, watering, simulation)
3. Create simple HTML frontend
4. Add real-time updates (polling every 2-5 seconds)
5. Write integration tests for API endpoints

Key Go concepts: HTTP handlers, JSON encoding, middleware

Package suggestions:
- Standard library: `net/http`
- Router: `github.com/go-chi/chi`

### Phase 6: Testing & Polish (Week 5)
Goal: Comprehensive testing and bug fixes

Tasks:
1. Add more unit tests (aim for >70% coverage)
2. Write integration tests for full simulation scenarios
3. Add logging (structured logging with `slog` or `zap`)
4. Handle edge cases (all plants die, negative saturation, etc.)
5. Add graceful shutdown

Key Go concepts: Table-driven tests, test coverage, `context` for cancellation

## Key Go Concepts You'll Learn

### Concurrency
- Goroutines: For simulation loop, scheduled watering, web server
- Channels: For manual watering commands, shutdown signals
- `select`: For handling multiple channels
- Mutexes: For protecting shared state (plant data)
- `context`: For graceful cancellation

### Interfaces
- Define interfaces for simulator, sensor manager, watering controller
- Enables testing via dependency injection and mocking

### Testing
- Unit tests for individual functions/methods
- Table-driven tests for multiple scenarios
- Integration tests for full simulation flows
- Test coverage with `go test -cover`

### Time Management
- `time.Ticker` for simulation ticks
- `time.Duration` for configurable intervals
- `time.Now()` for timestamps

### HTTP & JSON
- Building RESTful APIs
- Encoding/decoding JSON
- Middleware (logging, CORS)

### Configuration
- Reading config files
- Unmarshaling YAML/JSON into structs

## Testing Strategy

### Unit Tests
Test individual components in isolation.

Examples:
- Plant saturation depletion calculation
- Health update based on saturation
- Growth stage progression
- Death condition
- Sensor reading calculation for a section
- Watering amount calculation

### Integration Tests
Test components working together.

Examples:
- Full simulation tick: plants update, sensors read, watering triggers
- Manual watering command flows from API to plant state change
- Config file loading and greenhouse initialization

### Concurrency Tests
Test for race conditions with `go test -race`.

Scenarios:
- Concurrent sensor reads during plant updates
- Multiple manual watering requests
- Pause/resume during active watering

## Future Extensibility

### Architecture Patterns for New Features

#### Strategy pattern for environmental factors
Future sensors (temperature, light, humidity) can be added easily:
```go
type EnvironmentalFactor interface {
    AffectPlant(plant *Plant, reading float64)
}
```

#### Observer pattern for events
Emit events for:
- Plant death
- Low saturation alerts
- Watering completion
- Growth stage changes

#### Plugin-like configuration for plant types
Allow loading custom plant types from config without code changes.

### Planned Future Features
- Temperature: Affects growth rate and water evaporation
- Light: Required for photosynthesis, affects growth
- Humidity: Affects soil moisture evaporation rate
- Pests/Disease: Random events that affect health
- Harvesting: Mature plants can be harvested
- Analytics: Historical data, charts, predictions

## Resources to Learn

### Go Fundamentals
- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)

### Concurrency
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- "Concurrency in Go" by Katherine Cox-Buday

### Testing
- [Add a Test](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)

### Web Development
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)

## Getting Started

### Step 1: Setup Dependencies
```bash
go mod init greenhouse-simulator
go get gopkg.in/yaml.v3  # for YAML config
```

### Step 2: Start with Phase 1
Create `internal/models/plant.go` and define your plant structures.

### Step 3: Write Tests First (TDD)
For each feature, write tests before implementation.

Guiding questions:
1. Models: What data does this component need? What behavior should it have?
2. Concurrency: Does this run concurrently? How do I protect shared data?
3. Interfaces: Can I define an interface to make this testable?
4. Testing: What are the edge cases? How do I test concurrent code?
5. Errors: What can go wrong? How should I handle errors?

## Success Criteria

By the end of this project, you should be able to:

- Run a tick-based simulation with configurable time intervals
- See plants grow, degrade, and die based on water management
- Monitor soil saturation via sensors (grouped by section)
- Schedule automatic watering per section
- Manually trigger watering via web interface
- View real-time plant status in a web browser
- Have >70% test coverage
- Run `go test -race` with no race conditions
- Easily add new plant types via config
- Have an architecture ready for temperature/light/humidity features

## Final Notes

This is a learning project. Iterate, refactor, and explore Go idioms. Experiment with concurrency patterns, refactor when you find better designs, and use tests to guide your implementation.


