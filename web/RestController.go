package web

import (
	"net/http"
)

//RestController represents a web api controller for a given type T
type RestController interface {

	//Path gives the base path for the controller
	Path() string

	//Endpoints lists all endpoints of the controller
	Endpoints() []Endpoint

	//MiddleWare lists all middlewares to apply for all endpoints of the controller
	MiddleWare(handler http.Handler) http.Handler

	//ErrorHandler handles error occurring during request handling
	ErrorHandler() ErrorHandler
}

//Endpoint represents an Rest api endpoint
type Endpoint interface {

	//Path is the subpath for the endpoint
	Path() string

	//Handler represents the handler for the entering request
	Handler() http.Handler

	//Method http verb for the endpoint
	Method() HttpMethod

	//EmptyRequestBody default request body
	// Only needed for POST and PUT Method
	EmptyRequestBody() Request

	//ModelKey the key for storing the requestBody into the context
	// Only needed for POST and PUT Method
	ModelKey() interface{}
}
