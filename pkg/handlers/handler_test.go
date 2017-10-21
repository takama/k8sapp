package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/logger/standard"
	"github.com/takama/k8sapp/pkg/version"
)

func TestRoot(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Root)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, fmt.Sprintf("%s v%s", config.SERVICENAME, version.RELEASE))
}

func testHandler(t *testing.T, handler http.HandlerFunc, code int, body string) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != code {
		t.Error("Expected status code:", code, "got", trw.Code)
	}
	if trw.Body.String() != body {
		t.Error("Expected body", body, "got", trw.Body.String())
	}
}

func TestCollectCodes(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(func(c bit.Control) {
			c.Code(http.StatusBadGateway)
			c.Body(http.StatusText(http.StatusBadGateway))
		})(bit.NewControl(w, r))
	})
	testHandler(t, handler, http.StatusBadGateway, http.StatusText(http.StatusBadGateway))

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(func(c bit.Control) {
			c.Code(http.StatusNotFound)
			c.Body(http.StatusText(http.StatusNotFound))
		})(bit.NewControl(w, r))
	})
	testHandler(t, handler, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
