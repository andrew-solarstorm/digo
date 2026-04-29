package digo

// Instance Get instance and cast to TInstance
func Instance[T IInstance](id string, c IContainer) T {
	instance := c.Instance(id)
	return instance.(T)
}

func GetConfig[T IConfig](c IContainer) T {
	return c.Config().(T)
}
