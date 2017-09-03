package service

import "testing"

func TestRun(t *testing.T) {
	err := Run()
	if err != nil {
		t.Errorf("Fail, got '%s', want '%v'", err, nil)
	}
}
