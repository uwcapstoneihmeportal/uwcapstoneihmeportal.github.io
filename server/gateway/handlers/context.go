package handlers

import (
	"github.com/nimajalali/go-force/force"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the salesforce rest api object

type SessionContext struct {
	forceApi  *force.ForceApi
}

func NewContext(forceApi *force.ForceApi) *SessionContext {
	return &SessionContext{forceApi}
}
