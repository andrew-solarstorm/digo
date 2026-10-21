# dicontainer-go

A simple and flexible dependency injection container for Go applications.

## Installation

```bash
go get -u github.com/andrew-solarstorm/digo
```

## Features

- Simple API for managing application components
- Ordered initialization and clean shutdown
- Dependency management between components
- Graceful handling of OS signals
- Interface-based design for better testability
- Generic helper methods for type-safe dependency injection

## Interfaces

### DIContainerInterface

The main interface for dependency injection containers:

```go
// IContainer defines the public API for a dependency injection container
type IContainer interface {
 // Config returns a configuration by its key
 // If the config is not found, it will panic
 Config(key string) IConfig

 // Instance returns an instance by its ID
 // If the instance is not found, it will panic
 Instance(id string) IInstance

 // Run configures and starts all instances in the container
 // It also sets up handlers for graceful shutdown on OS signals
 Run() error

 // Stop all instances in the container
 // It also sets up handlers for graceful shutdown on OS signals
 Stop() error

 RunBlock() error
}
```

### IConfig

The interface that all configurations must implement. Multiple configs can be registered in the container, each identified by a unique, non-empty `Key`:

```go
type IConfig interface {
 // Key returns the name of the config which is unique
 Key() string

 // Load reads the configuration (env, file, ...)
 Load() error

 // Validate checks the loaded configuration
 Validate() error
}
```

### DIInstance

The interface that all injectable components must implement:

```go
type IInstance interface {
 // Id returns the name of the instance which is unique
 ID() string

 // all the configuration && initialization should be done here
 // dependency injection should be done here
 Configure(c IContainer) error

 // start the instance
 Start() error

 // stop the instance
 Stop() error
}
```

## Example Usage

```go
package main

import (
 "fmt"

 "github.com/andrew-solarstorm/digo"
)

// Configuration
type AppConfig struct {
 DBConnectionString string
 APIPort            int
}

const APP_CONFIG = "app-config"

func (c *AppConfig) Key() string { return APP_CONFIG }

func (c *AppConfig) Load() error {
 // read from env, file, ...
 return nil
}

func (c *AppConfig) Validate() error {
 return nil
}

const (
 DATABASE_SERVICE = "database-service"
 API_SERVICE = "api-service"
)

// Database service example
type DatabaseService struct {
 digo.BaseDIInstance
 db     *MockDB
 config *AppConfig
}

func (d *DatabaseService) Id() string {
 return DATABASE_SERVICE
}

func (d *DatabaseService) Configure(c digo.IContainer) error {
 fmt.Println("Configuring database service")
 d.config = c.Config(APP_CONFIG).(*AppConfig)
 return nil
}

func (d *DatabaseService) Start() error {
 fmt.Printf("Starting database service with connection: %s\n", d.config.DBConnectionString)
 d.db = &MockDB{}
 return nil
}

func (d *DatabaseService) Stop() error {
 fmt.Println("Stopping database service")
 return nil
}

// API service that depends on database
type APIService struct {
 digo.BaseDIInstance
 db     *DatabaseService
 config *AppConfig
}

func (a *APIService) Id() string {
 return API_SERVICE
}

func (a *APIService) Configure(c digo.IContainer) error {
 fmt.Println("Configuring API service")
 
 // Get config (standard way)
 a.config = c.Config(APP_CONFIG).(*AppConfig)
 
 // or Get config (type-safe generic helper)
 a.config = digo.GetConfig[*AppConfig](APP_CONFIG, c)
 
 // Get dependency (standard way)
 a.db = c.Instance(DATABASE_SERVICE).(*DatabaseService)
 
 // or Get dependency (type-safe generic helper)
 a.db = digo.Instance[*DatabaseService](DATABASE_SERVICE, c)
 return nil
}

func (a *APIService) Start() error {
 fmt.Printf("Starting API service on port: %d\n", a.config.APIPort)
 return nil
}

func (a *APIService) Stop() error {
 fmt.Println("Stopping API service")
 return nil
}

// Mock database for example
type MockDB struct{}

func main() {
 // Application configuration
 appConfig := &AppConfig{
  DBConnectionString: "postgres://localhost:5432/myapp",
  APIPort:            8080,
 }

 // Create container with configs (unique keys) and ordered instances
 diContainer, err := digo.New(
  []digo.IConfig{appConfig},
  &DatabaseService{},
  &APIService{},
 )
 if err != nil {
  panic(err)
 }

 // Run the container (this will configure, start all instances)
 go func() {
  if err := diContainer.Run(); err != nil {
   panic(err)
  }
 }()

 // Simulate running application
 fmt.Println("Application is running. Press Ctrl+C to stop.")

 // Keep the program running for a while
 time.Sleep(time.Second * 10)
 fmt.Println("Example completed.")
}
```

## Generic Helpers

The package provides type-safe generic helper methods:

```go
// Get config by key with type safety
config := digo.GetConfig[*AppConfig](APP_CONFIG, c)

// Get instance with type safety
dbService := digo.Instance[*DatabaseService](DATABASE_SERVICE, c)
```

## How It Works

1. Create one or more configuration structs that implement `IConfig`, each with a unique key
2. Create service instances that implement the `DIInstance` interface
3. Register the configs and the instances in the container (instances in the desired initialization order)
4. Run the container to start the application

The container handles:

- Loading and validating all configs in registration order (the first failure stops the run and reports the config key)
- Configuration of all instances
- Dependency injection between instances
- Starting instances in the registered order
- Gracefully stopping instances in reverse order when receiving termination signals

## License

MIT
