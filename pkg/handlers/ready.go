package handlers

import (
	"net/http"

	"github.com/takama/k8sapp/pkg/router"
)

// Ready returns "OK" if service is ready to serve traffic
func (h *Handler) Ready(c router.Control) {
	// TODO: possible use cases:
	// load data from a database, a message broker, any external services, etc

	c.Code(http.StatusOK)
	c.Write(http.StatusText(http.StatusOK))
}
