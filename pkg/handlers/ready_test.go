package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/logger/standard"
)

func TestReady(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}
