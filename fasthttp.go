package transport

import (
	"encoding/json"
	"fmt"

	"github.com/artfunder/structs"
	"github.com/valyala/fasthttp"
)

// NewFastHTTPTransport ...
func NewFastHTTPTransport(port int) *FastHTTPTransport {
	transport := new(FastHTTPTransport)

	fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), transport.HandleFastHTTP)

	return transport
}

// FastHTTPTransport ...
type FastHTTPTransport struct {
	getEndpoint Router
}

// Route ...
func (transport *FastHTTPTransport) Route(router Router) {
	transport.getEndpoint = router
}

// HandleFastHTTP ...
func (transport *FastHTTPTransport) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	err := transport.handleRequest(ctx)

	if err != nil {
		transport.handleError(err, ctx)
	}
}

func (transport FastHTTPTransport) handleRequest(ctx *fasthttp.RequestCtx) error {
	useService, err := transport.getEndpoint(transport.getPathAndMethod(ctx))
	if err != nil {
		return err
	}

	return transport.handleRequestWithEndpoint(ctx, useService)
}

func (transport FastHTTPTransport) handleRequestWithEndpoint(ctx *fasthttp.RequestCtx, useService Endpoint) error {
	return transport.useConnectorWithBody(ctx, useService, ctx.Request.Body())
}

func (transport FastHTTPTransport) useConnectorWithBody(
	ctx *fasthttp.RequestCtx,
	endpoint Endpoint,
	body []byte,
) error {
	res, err := endpoint.Receive(body)
	if err != nil {
		return err
	}

	return transport.sendResponse(res, ctx)
}

func (transport FastHTTPTransport) getPathAndMethod(ctx *fasthttp.RequestCtx) (string, string) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	return path, method
}

func (transport FastHTTPTransport) handleError(err error, ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(transport.getErrorCode(err))
	json.NewEncoder(ctx).Encode(structs.Error{Message: err.Error()})
}

func (transport FastHTTPTransport) getErrorCode(err error) int {
	if err == ErrorInternal {
		return 500
	}
	if err == ErrorNotFound {
		return 404
	}
	if err == ErrorBadMethod {
		return 403
	}
	return 400
}

func (transport FastHTTPTransport) sendResponse(res []byte, ctx *fasthttp.RequestCtx) error {
	isValidJSON := json.Valid(res)
	if !isValidJSON {
		return ErrorInternal
	}

	_, err := ctx.Write(res)
	if err != nil {
		return ErrorInternal
	}

	return nil
}
