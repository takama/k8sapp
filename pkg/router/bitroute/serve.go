// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bitroute

import (
	"net/http"
	"strings"

	"github.com/takama/k8sapp/pkg/router"
)

type bitroute struct {
	// List of handlers that associated with known http methods (GET, POST ...)
	handlers map[string]*parser

	// If enabled, the router automatically replies to OPTIONS requests.
	// Nevertheless OPTIONS handlers take priority over automatic replies.
	optionsRepliesEnabled bool

	// Configurable handler which is called when a request cannot be routed.
	notAllowed func(router.Control)

	// Configurable handler which is called when panic happen.
	recoveryHandler func(router.Control)

	// Configurable handler which is allowed to take control
	// before it is called standard methods e.g. GET, PUT.
	middlewareHandler func(func(router.Control)) func(router.Control)

	// Configurable http.Handler which is called when URL path has not defined method.
	// If it is not set, http.NotFound is used.
	notFound func(router.Control)
}

// New returns new router that implement Router interface.
func New() router.BitRoute {
	return &bitroute{
		handlers: make(map[string]*parser),
	}
}

// GET registers a new request handle for HTTP GET method.
func (r *bitroute) GET(path string, f func(router.Control)) {
	r.register("GET", path, f)
}

// PUT registers a new request handle for HTTP PUT method.
func (r *bitroute) PUT(path string, f func(router.Control)) {
	r.register("PUT", path, f)
}

// POST registers a new request handle for HTTP POST method.
func (r *bitroute) POST(path string, f func(router.Control)) {
	r.register("POST", path, f)
}

// DELETE registers a new request handle for HTTP DELETE method.
func (r *bitroute) DELETE(path string, f func(router.Control)) {
	r.register("DELETE", path, f)
}

// HEAD registers a new request handle for HTTP HEAD method.
func (r *bitroute) HEAD(path string, f func(router.Control)) {
	r.register("HEAD", path, f)
}

// OPTIONS registers a new request handle for HTTP OPTIONS method.
func (r *bitroute) OPTIONS(path string, f func(router.Control)) {
	r.register("OPTIONS", path, f)
}

// PATCH registers a new request handle for HTTP PATCH method.
func (r *bitroute) PATCH(path string, f func(router.Control)) {
	r.register("PATCH", path, f)
}

// If enabled, the router automatically replies to OPTIONS requests.
// Nevertheless OPTIONS handlers take priority over automatic replies.
// By default this option is disabled
func (r *bitroute) UseOptionsReplies(enabled bool) {
	r.optionsRepliesEnabled = enabled
}

// SetupNotAllowedHandler defines own handler which is called when a request
// cannot be routed.
func (r *bitroute) SetupNotAllowedHandler(f func(router.Control)) {
	r.notAllowed = f
}

// SetupNotFoundHandler allows to define own handler for undefined URL path.
// If it is not set, http.NotFound is used.
func (r *bitroute) SetupNotFoundHandler(f func(router.Control)) {
	r.notFound = f
}

// SetupRecoveryHandler allows to define handler that called when panic happen.
// The handler prevents your server from crashing and should be used to return
// http status code http.StatusInternalServerError (500)
func (r *bitroute) SetupRecoveryHandler(f func(router.Control)) {
	r.recoveryHandler = f
}

// SetupMiddleware defines handler is allowed to take control
// before it is called standard methods e.g. GET, PUT.
func (r *bitroute) SetupMiddleware(f func(func(router.Control)) func(router.Control)) {
	r.middlewareHandler = f
}

// Listen and serve on requested host and port
func (r *bitroute) Listen(hostPort string) error {
	return http.ListenAndServe(hostPort, r)
}

// registers a new handler with the given path and method.
func (r *bitroute) register(method, path string, f func(router.Control)) {
	if r.handlers[method] == nil {
		r.handlers[method] = newParser()
	}
	r.handlers[method].register(path, f)
}

func (r *bitroute) recovery(w http.ResponseWriter, req *http.Request) {
	if recv := recover(); recv != nil {
		c := newControl(w, req)
		r.recoveryHandler(c)
	}
}

// AllowedMethods returns list of allowed methods
func (r *bitroute) allowedMethods(path string) []string {
	var allowed []string
	for method, parser := range r.handlers {
		if _, _, ok := parser.get(path); ok {
			allowed = append(allowed, method)
		}
	}

	return allowed
}

// ServeHTTP implements http.Handler interface.
func (r *bitroute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.recoveryHandler != nil {
		defer r.recovery(w, req)
	}
	if _, ok := r.handlers[req.Method]; ok {
		if handle, params, ok := r.handlers[req.Method].get(req.URL.Path); ok {
			c := newControl(w, req)
			if len(params) > 0 {
				for _, item := range params {
					c.Param(item.key, item.value)
				}
			}
			if r.middlewareHandler != nil {
				r.middlewareHandler(handle)(c)
			} else {
				handle(c)
			}
			return
		}
	}
	allowed := r.allowedMethods(req.URL.Path)

	if len(allowed) == 0 {
		if r.notFound != nil {
			c := newControl(w, req)
			r.notFound(c)
		} else {
			http.NotFound(w, req)
		}
		return
	}

	w.Header().Set("Allow", strings.Join(allowed, ", "))
	if req.Method == "OPTIONS" && r.optionsRepliesEnabled {
		return
	}
	if r.notAllowed != nil {
		c := newControl(w, req)
		r.notAllowed(c)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
