package api

// Handler ...
type Handler func(c *ServiceContext) (body interface{}, err error)
