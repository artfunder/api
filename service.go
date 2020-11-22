package api

// Service ...
type Service interface {
	GetAll(ctx *ServiceContext) (interface{}, error)
	GetOne(ctx *ServiceContext) (interface{}, error)
	Create(ctx *ServiceContext) (interface{}, error)
	Update(ctx *ServiceContext) (interface{}, error)
	Delete(ctx *ServiceContext) (interface{}, error)
}
