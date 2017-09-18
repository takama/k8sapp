package system

import "testing"

func TestStubHandling(t *testing.T) {
	handling := new(Handling)
	err := handling.Reload()
	if err != ErrNotImplemented {
		t.Error("Expected error", ErrNotImplemented, "got", err)
	}
	err = handling.Maintenance()
	if err != ErrNotImplemented {
		t.Error("Expected error", ErrNotImplemented, "got", err)
	}
	err = handling.Shutdown()
	if err != ErrNotImplemented {
		t.Error("Expected error", ErrNotImplemented, "got", err)
	}
}
