// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"

	"github.com/takama/k8sapp/pkg/version"
)

const (
	// SERVICENAME contains a service name prefix which used in ENV variables
	SERVICENAME = "K8SAPP"
)

// Run starts the service
func Run() (err error) {
	fmt.Println(SERVICENAME, "Version:", version.RELEASE)

	return
}
