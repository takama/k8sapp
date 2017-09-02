// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"

	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/version"
)

// Run starts the service
func Run() (err error) {
	fmt.Println(config.SERVICENAME, "Version:", version.RELEASE)

	return
}
