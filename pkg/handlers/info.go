// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/takama/k8sapp/pkg/version"
)

// Status contains detailed information about service
type Status struct {
	Host     string   `json:"host"`
	Version  string   `json:"version"`
	Commit   string   `json:"commit"`
	Repo     string   `json:"repo"`
	Compiler string   `json:"compiler"`
	Runtime  Runtime  `json:"runtime"`
	State    State    `json:"state"`
	Requests Requests `json:"requests"`
}

// Runtime defines runtime part of service information
type Runtime struct {
	CPU        int    `json:"cpu"`
	Memory     string `json:"memory"`
	Goroutines int    `json:"goroutines"`
}

// State contains current state of the service
type State struct {
	Maintenance bool   `json:"maintenance"`
	Uptime      string `json:"uptime"`
}

// Requests collects responses statistics
type Requests struct {
	Duration Duration `json:"duration"`
	Codes    Codes    `json:"codes"`
}

// Duration collects responses duration
type Duration struct {
	Average string `json:"average"`
	Max     string `json:"max"`
}

// Codes contains response codes statistics
type Codes struct {
	C2xx int `json:"2xx"`
	C4xx int `json:"4xx"`
	C5xx int `json:"5xx"`
}

// Info returns detailed info about the service
func (h *Handler) Info(c bit.Control) {
	host, _ := os.Hostname()
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	c.Code(http.StatusOK)
	c.Body(Status{
		Host:     host,
		Version:  version.RELEASE,
		Commit:   version.COMMIT,
		Repo:     version.REPO,
		Compiler: runtime.Version(),
		Runtime: Runtime{
			CPU:        runtime.NumCPU(),
			Memory:     fmt.Sprintf("%.2fMB", float64(m.Sys)/(1<<(10*2))),
			Goroutines: runtime.NumGoroutine(),
		},
		State: State{
			Maintenance: h.maintenance,
			Uptime:      time.Now().Sub(h.stats.startTime).String(),
		},
		Requests: Requests{
			Duration: Duration{
				Average: h.stats.requests.Duration.Average,
				Max:     h.stats.requests.Duration.Max,
			},
			Codes: Codes{
				C2xx: h.stats.requests.Codes.C2xx,
				C4xx: h.stats.requests.Codes.C4xx,
				C5xx: h.stats.requests.Codes.C5xx,
			},
		},
	})
}
