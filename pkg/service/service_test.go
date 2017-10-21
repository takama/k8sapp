package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/handlers"
)

func TestSetup(t *testing.T) {
	cfg := new(config.Config)
	err := cfg.Load(config.SERVICENAME)
	if err != nil {
		t.Error("Expected loading of environment vars, got", err)
	}
	router, logger, err := Setup(cfg)
	if err != nil {
		t.Errorf("Fail, got '%s', want '%v'", err, nil)
	}
	if router == nil {
		t.Error("Expected new router, got nil")
	}
	if logger == nil {
		t.Error("Expected new logger, got nil")
	}

	h := handlers.New(logger, cfg)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(notFound)(bit.NewControl(w, r))
	})

	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != http.StatusNotFound {
		t.Error("Expected status:", http.StatusNotFound, "got", trw.Code)
	}
}
