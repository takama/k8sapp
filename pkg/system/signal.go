// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package system

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/takama/k8sapp/pkg/logger"
)

// SignalType describe
type SignalType int

const (
	// Shutdown defines signals for shutdown process
	Shutdown SignalType = iota
	// Reload defines signals for reload process
	Reload
	// Maintenance defines signals for maintenance process
	Maintenance
)

func (s SignalType) String() string {
	switch s {
	case Shutdown:
		return "SHUTDOWN"
	case Reload:
		return "RELOAD"
	case Maintenance:
		return "MAINTENANCE"
	default:
		return strconv.Itoa(int(s))
	}
}

// Signals for defined processes
type Signals struct {
	mutex sync.RWMutex

	interrupt chan os.Signal
	quit      chan struct{}

	shutdown    []os.Signal
	reload      []os.Signal
	maintenance []os.Signal
}

// NewSignals creates default signals
func NewSignals() *Signals {
	signals := &Signals{
		// Set up channel on which to send signal notifications.
		// We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		interrupt: make(chan os.Signal, 1),
		quit:      make(chan struct{}, 1),

		shutdown:    []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		reload:      []os.Signal{syscall.SIGHUP},
		maintenance: []os.Signal{syscall.SIGUSR1},
	}
	signal.Notify(signals.interrupt)
	return signals
}

// Get signals by specified type
func (s *Signals) Get(sigType SignalType) (signals []os.Signal) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	switch sigType {
	case Shutdown:
		signals = make([]os.Signal, len(s.shutdown))
		copy(signals, s.shutdown)
	case Reload:
		signals = make([]os.Signal, len(s.reload))
		copy(signals, s.reload)
	case Maintenance:
		signals = make([]os.Signal, len(s.maintenance))
		copy(signals, s.maintenance)
	}

	return
}

// Add appends signal by specified type
func (s *Signals) Add(sig os.Signal, sigType SignalType) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	switch sigType {
	case Shutdown:
		s.shutdown = append(s.shutdown, sig)
	case Reload:
		s.reload = append(s.reload, sig)
	case Maintenance:
		s.maintenance = append(s.maintenance, sig)
	}
}

// Remove deletes signal by specified type
func (s *Signals) Remove(sig os.Signal, sigType SignalType) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	switch sigType {
	case Shutdown:
		s.shutdown = removeSignal(sig, s.shutdown)
	case Reload:
		s.reload = removeSignal(sig, s.reload)
	case Maintenance:
		s.maintenance = removeSignal(sig, s.maintenance)
	}
}

// Wait needs to catch signal and do graceful shutdown
func (s *Signals) Wait(logger logger.Logger, operator Operator) error {
	for {
		select {
		case <-s.quit:
			logger.Info("Gracefully closed")
			return nil
		case sig := <-s.interrupt:
			s.mutex.RLock()
			logger.Infof("Got signal: %s", sig)
			switch {
			case isSignalAvailable(sig, s.maintenance):
				s.mutex.RUnlock()
				logger.Info("Maintenance request")
				err := operator.Maintenance()
				if err != nil {
					logger.Error(err)
				}
			case isSignalAvailable(sig, s.reload):
				s.mutex.RUnlock()
				logger.Info("Reloading configuration...")
				err := operator.Reload()
				if err != nil {
					logger.Error(err)
				}
			case isSignalAvailable(sig, s.shutdown):
				s.mutex.RUnlock()
				logger.Info("Service was terminated by system signal")
				err := operator.Shutdown()
				if err != nil {
					logger.Error(err)
				}
				s.quit <- struct{}{}
			}
		}
	}
}

// Checks if a signal is available in the specified list
func isSignalAvailable(signal os.Signal, list []os.Signal) bool {
	for _, s := range list {
		if s == signal {
			return true
		}
	}
	return false
}

// Removes a signal from the specified list
func removeSignal(signal os.Signal, list []os.Signal) (signals []os.Signal) {
	for ind, sig := range list {
		if sig == signal {
			signals = append(list[:ind], list[ind+1:]...)
			return
		}
	}
	return list
}
