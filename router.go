package api

import (
	"encoding/json"
	"net/http"

	"github.com/artfunder/structs"
	"github.com/julienschmidt/httprouter"
)

var (
	actionGetOne = "get one"
	actionGetAll = "get all"
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

// Router handles the CRUD methods for an API Router
type Router struct {
	service Service
}

// RouteService creates a new Router
func RouteService(route string, service Service) {
	router := new(Router)
	router.service = service

	httpRouter := GetHTTPRouter()

	httpRouter.GET(route, router.getRouteHandler(actionGetAll))
	httpRouter.GET(route+"/:id", router.getRouteHandler(actionGetOne))
	httpRouter.POST(route, router.getRouteHandler(actionCreate))
	httpRouter.PATCH(route+"/:id", router.getRouteHandler(actionUpdate))
	httpRouter.DELETE(route+"/:id", router.getRouteHandler(actionDelete))
}

func (router Router) getRouteHandler(action string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		router.handleRoute(NewServiceContext(w, r, p, action))
	}
}

func (router Router) handleRoute(ctx *ServiceContext) {
	ctx.w.Header().Set("Content-Type", "application/json")

	ctx.object, ctx.err = router.runHandlerForAction(ctx)

	router.handleRouteIfHandlerExists(ctx)
}

func (router Router) handleRouteIfHandlerExists(ctx *ServiceContext) {
	if ctx.object == nil && ctx.err == nil {
		router.returnError(ctx.w, ErrorNoHandler)
		return
	}

	router.returnObjectOrError(ctx)
}

func (router Router) returnObjectOrError(ctx *ServiceContext) {
	if ctx.err != nil {
		router.returnError(ctx.w, ctx.err)
		return
	}

	router.json(ctx.w, ctx.object, router.getSuccessCode(ctx))
}

func (router Router) json(w http.ResponseWriter, object interface{}, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}

func (router Router) runHandlerForAction(ctx *ServiceContext) (interface{}, error) {
	switch ctx.GetActionType() {
	case actionGetAll:
		return router.service.GetAll(ctx)
	case actionGetOne:
		return router.service.GetOne(ctx)
	case actionCreate:
		return router.service.Create(ctx)
	case actionUpdate:
		return router.service.Update(ctx)
	case actionDelete:
		return router.service.Delete(ctx)
	}

	return nil, nil
}

func (router Router) returnError(w http.ResponseWriter, err error) {
	w.WriteHeader(router.getErrorCode(err))
	json.NewEncoder(w).Encode(structs.Error{
		Message: err.Error(),
	})
}

func (router Router) getErrorCode(err error) int {
	if err == ErrorNotFound {
		return 404
	}
	if err == ErrorNoHandler {
		return 405
	}
	return 400
}

func (router Router) getSuccessCode(ctx *ServiceContext) int {
	if ctx.actionType == actionCreate {
		return 201
	}
	return 200
}
