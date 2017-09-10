package router

import "testing"

func TestNewHTTPRouter(t *testing.T) {
	r := NewHTTPRouter()
	if r != nil {
		t.Error("Expected nil, got", r)
	}
}
