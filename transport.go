package transport

// Transport ...
type Transport interface {
	Route(Router)
}

// Router ...
type Router func(path string, method string) (Endpoint, error)

// Endpoint ...
type Endpoint interface {
	Receive([]byte) ([]byte, error)
}

// EndpointFunc ...
type EndpointFunc func(body []byte) ([]byte, error)

// Receive ...
func (f EndpointFunc) Receive(body []byte) ([]byte, error) {
	return f(body)
}
