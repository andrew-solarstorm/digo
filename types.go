package digo

type IConfig interface {
	Load() error
	Validate() error
}

// IContainer defines the public API for a dependency injection container
type IContainer interface {
	// Getconfig returns the configuration object
	Config() IConfig

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

// Ensure DIContainer implements DIContainerInterface
var (
	_ IContainer = (*Container)(nil)
	_ IInstance  = (*BaseDIInstance)(nil)
)
