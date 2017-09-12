package httprouter

import (
	"net/http"
	"testing"
)

func TestNewHTTPRouter(t *testing.T) {
	r := New()
	if r == nil {
		t.Error("Expected new httprouter, got nil")
	}
}
func TestHTTPRouter(t *testing.T) {
	r := new(httpRouter)
	if r.HandleOPTIONS {
		t.Error("Expected of handle OPTIONS is not set")
	}
	r.UseOptionsReplies(true)
	if !r.HandleOPTIONS {
		t.Error("Expected of handle OPTIONS is set")
	}
	if r.MethodNotAllowed != nil {
		t.Error("Expected nil, got", r.MethodNotAllowed)
	}
	r.SetupNotAllowedHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if r.MethodNotAllowed == nil {
		t.Error("Expected handler, got nil")
	}
	if r.NotFound != nil {
		t.Error("Expected nil, got", r.NotFound)
	}
	r.SetupNotFoundHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if r.NotFound == nil {
		t.Error("Expected handler, got nil")
	}
	if r.PanicHandler != nil {
		t.Error("Expected nil, got not nil")
	}
	r.SetupRecoveryHandler(func(http.ResponseWriter, *http.Request, interface{}) {})
	if r.PanicHandler == nil {
		t.Error("Expected handler, got nil")
	}
	err := r.Listen("$")
	if err == nil {
		t.Error("Expected error if used incorrect host and port")
	}
}
