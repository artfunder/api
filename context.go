package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ServiceContext ...
type ServiceContext struct {
	body       map[string]string
	id         int
	actionType string
	w          http.ResponseWriter
	object     interface{}
	err        error
}

// NewServiceContext ...
func NewServiceContext(w http.ResponseWriter, r *http.Request, p httprouter.Params, actionType string) *ServiceContext {
	context := new(ServiceContext)

	context.w = w
	context.body = context.parseBodyFromRequest(r)
	context.id = context.parseIDFromParams(p)
	context.actionType = actionType

	return context
}

/** Public Methods */

// GetID ...
func (c ServiceContext) GetID() int {
	return c.id
}

// GetBody ..
func (c ServiceContext) GetBody() map[string]string {
	return c.body
}

// GetBodyInto ...
func (c ServiceContext) GetBodyInto(dest interface{}) {
	jsonBytes := bytes.NewBuffer([]byte(""))
	json.NewEncoder(jsonBytes).Encode(c.body)
	json.NewDecoder(jsonBytes).Decode(dest)
}

// GetActionType ...
func (c ServiceContext) GetActionType() string {
	return c.actionType
}

/** Private Methods */

func (ServiceContext) parseIDFromParams(params httprouter.Params) int {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return 0
	}
	return id
}

func (ServiceContext) parseBodyFromRequest(r *http.Request) map[string]string {
	if r.Body == nil {
		return nil
	}
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	return body
}
