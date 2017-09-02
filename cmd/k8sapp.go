// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/takama/k8sapp/pkg/service"
)

func main() {
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
