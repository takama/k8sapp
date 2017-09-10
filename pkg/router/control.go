// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"
)

type control struct {
	req    *http.Request
	w      http.ResponseWriter
	code   int
	params []struct {
		key   string
		value string
	}
}

// newControl returns new control that implement Control interface.
func newControl(w http.ResponseWriter, req *http.Request) Control {
	return &control{
		req: req,
		w:   w,
	}
}

// Request returns *http.Request
func (c *control) Request() *http.Request {
	return c.req
}

// Query searches URL/Post value by key.
// If there are no values associated with the key, an empty string is returned.
func (c *control) Query(key string) string {
	for idx := range c.params {
		if c.params[idx].key == key {
			return c.params[idx].value
		}
	}

	return c.req.URL.Query().Get(key)
}

// Param sets URL/Post key/value params.
func (c *control) Param(key, value string) {
	c.params = append(c.params, struct{ key, value string }{key: key, value: value})
}

// Response writer section
// Header represents http.ResponseWriter header, the key-value pairs in an HTTP header.
func (c *control) Header() http.Header {
	return c.w.Header()
}

// Code sets HTTP status code e.g. http.StatusOk
func (c *control) Code(code int) {
	if code >= 100 && code < 600 {
		c.code = code
	}
}

// Write writes data into http output.
func (c *control) Write(data interface{}) {
	var content []byte

	if str, ok := data.(string); ok {
		content = []byte(str)
	} else {
		var err error
		content, err = json.Marshal(data)
		if err != nil {
			c.w.WriteHeader(http.StatusInternalServerError)
			c.w.Write([]byte(err.Error()))
			return
		}
		if c.w.Header().Get("Content-type") == "" {
			c.w.Header().Add("Content-type", "application/json")
		}
	}
	if strings.Contains(c.req.Header.Get("Accept-Encoding"), "gzip") {
		c.w.Header().Add("Content-Encoding", "gzip")
		if c.code > 0 {
			c.w.WriteHeader(c.code)
		}
		gz := gzip.NewWriter(c.w)
		gz.Write(content)
		gz.Close()
	} else {
		if c.code > 0 {
			c.w.WriteHeader(c.code)
		}
		c.w.Write(content)
	}
}
