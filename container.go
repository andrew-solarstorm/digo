// Package digo is a dependency injection container for Go
// Author: Andrew Solarstorm
package digo

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const _contanerName = "Container"

// Container is a dependency injection container
type Container struct {
	iOrder    map[int]string
	instances map[string]IInstance
	config    IConfig
	ctx       context.Context
	cancel    context.CancelFunc
}

func (c *Container) Config() IConfig {
	return c.config
}

// New create DIContainer with DIInstances
func New(config IConfig, ins ...IInstance) (*Container, error) {
	c := &Container{
		iOrder:    make(map[int]string),
		instances: make(map[string]IInstance),
		config:    config,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	for order, i := range ins {
		if err := c.register(order, i); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Container) cInfo(action string, instance string) *zerolog.Event {
	return log.Info().Str("DI Container", action).Str("Instance", instance)
}

func (c *Container) cError(action string, instance string) *zerolog.Event {
	return log.Error().Str("DI Container", action).Str("Instance", instance)
}

// register register an instance
func (c *Container) register(idx int, instance IInstance) error {
	// check if instance is already registered
	if _, ok := c.instances[instance.ID()]; !ok {

		// register instance with order
		c.iOrder[idx] = instance.ID()
		c.instances[instance.ID()] = instance

		c.cInfo("register", instance.ID()).Msgf("registered instance: %s order: %d", instance.ID(), idx)
		return nil
	}
	err := fmt.Errorf("instance %s already registered", instance.ID())
	c.cError("register", _contanerName).Err(err)
	return err
}

// Instance returns an instance
// if instance is not found, return error
func (c *Container) Instance(id string) IInstance {
	if i, ok := c.instances[id]; !ok {
		panic(fmt.Errorf("instance %s not found", id))
	} else {
		return i
	}
}

// start the instances in order
func (c *Container) start() error {
	c.cInfo("start", _contanerName).Msg("Application started. Press Ctrl+C to stop.")
	for idx := range len(c.iOrder) {
		iname := c.iOrder[idx]
		if err := c.instances[iname].Start(); err != nil {
			c.cError("start", iname).Err(err).Msg("error staring service")
			return err
		}
		c.cInfo("start", iname).Msgf("starting instance: %s", iname)
	}

	return nil
}

// configure the instances in order
func (c *Container) configure() error {
	if c.config != nil {
		if err := c.config.Load(); err != nil {
			return fmt.Errorf("error loading config: %s", err.Error())
		}

		if err := c.config.Validate(); err != nil {
			return fmt.Errorf("error validating config: %s", err.Error())
		}
	}

	for id := range len(c.iOrder) {
		iname := c.iOrder[id]
		if err := c.instances[iname].Configure(c); err != nil {
			c.cError("configure", iname).Err(err).Msg("error configuring service")
			return err
		}
		c.cInfo("configure", iname).Msgf("configured instance: %d", id)
	}
	return nil
}

// Stop the instances in reverse order
func (c *Container) Stop() error {
	for id := range len(c.iOrder) {
		name := c.iOrder[id]
		if err := c.instances[name].Stop(); err != nil {
			c.cError("stop", name).Err(err).Msg("error stopping service")
			return err
		}
		c.cInfo("stop", name).Msgf("stopped instance: %d", id)
	}
	return nil
}

func (c *Container) run() error {
	// load config
	if c.config != nil {
		if err := c.config.Load(); err != nil {
			return err
		}
	}

	// config instances first
	if err := c.configure(); err != nil {
		return err
	}

	// then run instances
	if err := c.start(); err != nil {
		return err
	}
	return nil
}

// Run the container with context
func (c *Container) Run() error {
	// run with cancel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// stop the container
	go func() {
		select {
		case <-c.ctx.Done():
			c.cInfo("run", _contanerName).Msg("container stopped")
			_ = c.Stop()
			return
		case <-sigChan:
			c.cInfo("run", _contanerName).Msg("container stopped by user cancel")
			c.cancel()
			_ = c.Stop()
			return
		}
	}()

	// run without block
	if err := c.run(); err != nil {
		return err
	}

	return nil
}

// RunBlock implements IContainer.
func (c *Container) RunBlock() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// run block
	if err := c.run(); err != nil {
		return err
	}

	// wait for signal
	<-sigChan

	return nil
}
