// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
)

// Health returns "OK" if service is alive
func (h *Handler) Health(c bit.Control) {
	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}
