// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package system

import "errors"

// ErrNotImplemented declares error for method that isn't implemented
var ErrNotImplemented = errors.New("This method is not implemented")

// Operator defines reload, maintenance and shutdown interface
type Operator interface {
	Reload() error
	Maintenance() error
	Shutdown() error
}

// Handling implements simplest Operator interface
type Handling struct{}

// Reload operation implementation
func (h Handling) Reload() error {
	return ErrNotImplemented
}

// Maintenance operation implementation
func (h Handling) Maintenance() error {
	return ErrNotImplemented
}

// Shutdown operation implementation
func (h Handling) Shutdown() error {
	return ErrNotImplemented
}
