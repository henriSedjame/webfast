package web

import "net/http"

//ErrorHandler represents a function that handle error
type ErrorHandler func(error, http.ResponseWriter) error

//RequestHandler represents a function that handles a http request
type RequestHandler func(http.ResponseWriter, *http.Request)

//HttpMethod represents a http method name
type HttpMethod = string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	OPTIONS HttpMethod = "OPTIONS"
	PATCH   HttpMethod = "PATCH"
	HEAD    HttpMethod = "HEAD"
	CONNECT HttpMethod = "CONNECT"
	TRACE   HttpMethod = "TRACE"
)
