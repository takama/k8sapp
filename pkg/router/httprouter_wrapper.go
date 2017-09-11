// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type httpRouter struct {
	httprouter.Router
}

func newHTTPRouter() HTTPRouter {
	router := new(httpRouter)
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	return router
}

// If enabled, the router automatically replies to OPTIONS requests.
// Nevertheless OPTIONS handlers take priority over automatic replies.
// By default this option is disabled
func (hr *httpRouter) UseOptionsReplies(enabled bool) {
	hr.HandleOPTIONS = enabled
}

// SetupNotAllowedHandler defines own handler which is called when a request
// cannot be routed.
func (hr *httpRouter) SetupNotAllowedHandler(h http.Handler) {
	hr.MethodNotAllowed = h
}

// SetupNotFoundHandler allows to define own handler for undefined URL path.
// If it is not set, http.NotFound is used.
func (hr *httpRouter) SetupNotFoundHandler(h http.Handler) {
	hr.NotFound = h
}

// SetupRecoveryHandler allows to define handler that called when panic happen.
// The handler prevents your server from crashing and should be used to return
// http status code http.StatusInternalServerError (500)
// interface{} will contain value which is transmitted from panic call.
func (hr *httpRouter) SetupRecoveryHandler(f func(http.ResponseWriter, *http.Request, interface{})) {
	hr.PanicHandler = f
}

// Listen and serve on requested host and port e.g "0.0.0.0:8080"
func (hr *httpRouter) Listen(hostPort string) error {
	return http.ListenAndServe(hostPort, hr)
}
