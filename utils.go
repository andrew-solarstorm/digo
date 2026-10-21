package digo

// Instance Get instance and cast to TInstance
func Instance[T IInstance](id string, c IContainer) T {
	instance := c.Instance(id)
	return instance.(T)
}

// GetConfig Get config by key and cast to TConfig
func GetConfig[T IConfig](key string, c IContainer) T {
	return c.Config(key).(T)
}
