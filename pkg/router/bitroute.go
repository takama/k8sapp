// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import "net/http"

// Control interface contains methods that control
// HTTP header, URL/post query parameters, request/response
// and HTTP output like Code(), Write(), etc.
type Control interface {
	// Request returns *http.Request
	Request() *http.Request

	// Query searches URL/Post query parameters by key.
	// If there are no values associated with the key, an empty string is returned.
	Query(key string) string

	// Param sets URL/Post key/value query parameters.
	Param(key, value string)

	// Response writer section

	// Header represents http.ResponseWriter header, the key-value pairs in an HTTP header.
	Header() http.Header

	// Code sets HTTP status code e.g. http.StatusOk
	Code(code int)

	// GetCode shows HTTP status code that set by Code()
	GetCode() int

	// Write prepared header, status code and body data into http output.
	Write(data interface{})

	// TODO Add more control methods.
}

// BitRoute interface contains base http methods e.g. GET, PUT, POST
// and defines your own handlers that is useful in some use cases
type BitRoute interface {
	// Standard methods

	// GET registers a new request handle for HTTP GET method.
	GET(path string, f func(Control))
	// PUT registers a new request handle for HTTP PUT method.
	PUT(path string, f func(Control))
	// POST registers a new request handle for HTTP POST method.
	POST(path string, f func(Control))
	// DELETE registers a new request handle for HTTP DELETE method.
	DELETE(path string, f func(Control))
	// HEAD registers a new request handle for HTTP HEAD method.
	HEAD(path string, f func(Control))
	// OPTIONS registers a new request handle for HTTP OPTIONS method.
	OPTIONS(path string, f func(Control))
	// PATCH registers a new request handle for HTTP PATCH method.
	PATCH(path string, f func(Control))

	// User defined options and handlers

	// If enabled, the router automatically replies to OPTIONS requests.
	// Nevertheless OPTIONS handlers take priority over automatic replies.
	// By default this option is disabled
	UseOptionsReplies(bool)

	// SetupNotAllowedHandler defines own handler which is called when a request
	// cannot be routed.
	SetupNotAllowedHandler(func(Control))

	// SetupNotFoundHandler allows to define own handler for undefined URL path.
	// If it is not set, http.NotFound is used.
	SetupNotFoundHandler(func(Control))

	// SetupRecoveryHandler allows to define handler that called when panic happen.
	// The handler prevents your server from crashing and should be used to return
	// http status code http.StatusInternalServerError (500)
	SetupRecoveryHandler(func(Control))

	// SetupMiddleware defines handler that is allowed to take control
	// before it is called standard methods above e.g. GET, PUT.
	SetupMiddleware(func(func(Control)) func(Control))

	// Listen and serve on requested host and port e.g "0.0.0.0:8080"
	Listen(hostPort string) error
}
