// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package system

// Operator defines reload, maintenance and shutdown interface
type Operator interface {
	Reload() error
	Maintenance() error
	Shutdown() error
}

// Handling implements simplest Operator interface
type Handling struct{}

// Reload implementation
func (h Handling) Reload() error {
	return nil
}

// Maintenance implementation
func (h Handling) Maintenance() error {
	return nil
}

// Shutdown implementation
func (h Handling) Shutdown() error {
	return nil
}
