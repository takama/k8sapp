// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTPRouter interface contains base http methods e.g. GET, PUT, POST
// and defines your own handlers that is useful in some use cases
type HTTPRouter interface {
	// Standard methods

	// GET registers a new request handle for HTTP GET method.
	GET(path string, h httprouter.Handle)
	// PUT registers a new request handle for HTTP PUT method.
	PUT(path string, h httprouter.Handle)
	// POST registers a new request handle for HTTP POST method.
	POST(path string, h httprouter.Handle)
	// DELETE registers a new request handle for HTTP DELETE method.
	DELETE(path string, h httprouter.Handle)
	// HEAD registers a new request handle for HTTP HEAD method.
	HEAD(path string, h httprouter.Handle)
	// OPTIONS registers a new request handle for HTTP OPTIONS method.
	OPTIONS(path string, h httprouter.Handle)
	// PATCH registers a new request handle for HTTP PATCH method.
	PATCH(path string, h httprouter.Handle)

	// User defined options and handlers

	// If enabled, the router automatically replies to OPTIONS requests.
	// Nevertheless OPTIONS handlers take priority over automatic replies.
	// By default this option is enabled
	UseOptionsReplies(bool)

	// SetupNotAllowedHandler defines own handler which is called when a request
	// cannot be routed.
	SetupNotAllowedHandler(http.Handler)

	// SetupNotFoundHandler allows to define own handler for undefined URL path.
	// If it is not set, http.NotFound is used.
	SetupNotFoundHandler(http.Handler)

	// SetupRecoveryHandler allows to define handler that called when panic happen.
	// The handler prevents your server from crashing and should be used to return
	// http status code http.StatusInternalServerError (500)
	// interface{} will contain value which is transmitted from panic call.
	SetupRecoveryHandler(func(http.ResponseWriter, *http.Request, interface{}))

	// Listen and serve on requested host and port e.g "0.0.0.0:8080"
	Listen(hostPort string) error
}
