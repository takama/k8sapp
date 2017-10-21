// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/version"
)

// Handler defines common part for all handlers
type Handler struct {
	logger      logger.Logger
	config      *config.Config
	maintenance bool
	stats       *stats
}

type stats struct {
	requests        *Requests
	averageDuration time.Duration
	maxDuration     time.Duration
	totalDuration   time.Duration
	requestsCount   time.Duration
	startTime       time.Time
}

// New returns new instance of the Handler
func New(logger logger.Logger, config *config.Config) *Handler {
	return &Handler{
		logger: logger,
		config: config,
		stats: &stats{
			requests:  new(Requests),
			startTime: time.Now(),
		},
	}
}

// Base handler implements middleware logic
func (h *Handler) Base(handle func(bit.Control)) func(bit.Control) {
	return func(c bit.Control) {
		timer := time.Now()
		handle(c)
		h.countDuration(timer)
		h.collectCodes(c)
	}
}

// Root handler shows version
func (h *Handler) Root(c bit.Control) {
	c.Code(http.StatusOK)
	c.Body(fmt.Sprintf("%s v%s", config.SERVICENAME, version.RELEASE))
}

func (h *Handler) countDuration(timer time.Time) {
	if !timer.IsZero() {
		h.stats.requestsCount++
		took := time.Now()
		duration := took.Sub(timer)
		h.stats.totalDuration += duration
		if duration > h.stats.maxDuration {
			h.stats.maxDuration = duration
		}
		h.stats.averageDuration = h.stats.totalDuration / h.stats.requestsCount
		h.stats.requests.Duration.Max = h.stats.maxDuration.String()
		h.stats.requests.Duration.Average = h.stats.averageDuration.String()
	}
}

func (h *Handler) collectCodes(c bit.Control) {
	if c.GetCode() >= 500 {
		h.stats.requests.Codes.C5xx++
	} else {
		if c.GetCode() >= 400 {
			h.stats.requests.Codes.C4xx++
		} else {
			if c.GetCode() >= 200 && c.GetCode() < 300 {
				h.stats.requests.Codes.C2xx++
			}
		}
	}
}
