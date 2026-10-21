package digo

import "fmt"

type configs struct {
	order []string
	byKey map[string]IConfig
}

func newConfigs() *configs {
	return &configs{byKey: make(map[string]IConfig)}
}

func (s *configs) register(cfg IConfig) error {
	if cfg == nil {
		return fmt.Errorf("config must not be nil")
	}

	key := cfg.Key()
	if key == "" {
		return fmt.Errorf("config key must not be empty")
	}
	if _, ok := s.byKey[key]; ok {
		return fmt.Errorf("config %s already registered", key)
	}

	s.order = append(s.order, key)
	s.byKey[key] = cfg
	return nil
}

func (s *configs) get(key string) (IConfig, bool) {
	cfg, ok := s.byKey[key]
	return cfg, ok
}

func (s *configs) load() error {
	for _, key := range s.order {
		cfg := s.byKey[key]
		if err := cfg.Load(); err != nil {
			return fmt.Errorf("error loading config %s: %w", key, err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("error validating config %s: %w", key, err)
		}
	}
	return nil
}
