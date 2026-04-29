package digo

type BaseDIInstance struct {
	IInstance
}

func (i *BaseDIInstance) ID() string {
	return ""
}

func (i *BaseDIInstance) Configure(c IContainer) error {
	return nil
}

func (i *BaseDIInstance) Start() error {
	return nil
}

func (i *BaseDIInstance) Stop() error {
	return nil
}
