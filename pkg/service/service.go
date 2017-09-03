// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/version"
)

// Run starts the service
func Run() (err error) {
	// Setup logger
	log := logger.New(&logger.Config{
		Level: logger.LevelDebug,
		Time:  true,
		UTC:   true,
	})

	log.Info("Version:", version.RELEASE)
	log.Warnf("%s log level is used", logger.LevelDebug.String())

	return
}
