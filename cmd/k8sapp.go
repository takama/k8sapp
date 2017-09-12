// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/service"
)

func main() {
	// Load ENV configuration
	cfg := new(config.Config)
	if err := cfg.Load(config.SERVICENAME); err != nil {
		log.Fatal(err)
	}

	// Configure service and get router
	router, err := service.Setup(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Listen and serve handlers
	router.Listen(fmt.Sprintf("%s:%d", cfg.LocalHost, cfg.LocalPort))
}
