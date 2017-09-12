package service

import (
	"testing"

	"github.com/takama/k8sapp/pkg/config"
)

func TestRun(t *testing.T) {
	cfg := new(config.Config)
	err := cfg.Load(config.SERVICENAME)
	if err != nil {
		t.Error("Expected loading of environment vars, got", err)
	}
	router, err := Setup(cfg)
	if err != nil {
		t.Errorf("Fail, got '%s', want '%v'", err, nil)
	}
	if router == nil {
		t.Error("Expected new router, got nil")
	}
}
