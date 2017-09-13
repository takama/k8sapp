package handlers

import (
	"fmt"
	"net/http"

	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/router"
	"github.com/takama/k8sapp/pkg/version"
)

// Handler defines common part for all handlers
type Handler struct {
	logger logger.Logger
	config *config.Config
}

// New returns new instance of the Handler
func New(logger logger.Logger, config *config.Config) *Handler {
	return &Handler{
		logger: logger,
		config: config,
	}
}

// Base handler implements middleware logic
func (h *Handler) Base(handle func(router.Control)) func(router.Control) {
	return func(c router.Control) {

		// TODO: Add custom logic here

		handle(c)
	}
}

// Root handler shows version
func (h *Handler) Root(c router.Control) {
	c.Code(http.StatusOK)
	c.Write(fmt.Sprintf("%s v%s", config.SERVICENAME, version.RELEASE))
}
