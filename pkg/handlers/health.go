package handlers

import (
	"net/http"

	"github.com/takama/k8sapp/pkg/router"
)

// Health returns "OK" if service is alive
func (h *Handler) Health(c router.Control) {
	c.Code(http.StatusOK)
	c.Write(http.StatusText(http.StatusOK))
}
