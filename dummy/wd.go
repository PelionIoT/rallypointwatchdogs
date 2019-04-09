package main

// Copyright (c) 2018, Arm Limited and affiliates.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"time"

	"github.com/armPelionEdge/maestroSpecs"
)

const (
	internalInterval = 15 // seconds
	stop             = 1
	keepalive        = 2
)

// This is a Go plugin, meeting the specs of maestroSpecs/watchdog.go

type watchdog struct {
	ready    bool
	log      maestroSpecs.WatchdogLogger
	ctrlChan chan int
}

// Watchdog is the single instance of this watchdog
var _watchdog *watchdog

// Watchdog this is the exported interface of the watchdog
var Watchdog maestroSpecs.Watchdog

//var Watchdog maestroSpecs.Watchdog

// InitMaestroPlugin is always called by Maestro any time it inits a plugin
func InitMaestroPlugin() (err error) {
	//	fmt.Printf("InitMaestroPlugin() dummy\n")
	_watchdog = new(watchdog)
	Watchdog = _watchdog
	return
}

func (wd *watchdog) watchdogRunner() {
	if wd.ready {
		wd.log.Errorf("Watchdog (dummy): ERROR - the watchdog was enabled twice.")
		return
	}
	wd.ready = true
	interval := internalInterval * time.Second
wdLoop:
	for {
		wd.log.Debugf("Watchdog (dummy): top of wdLoop.")
		select {
		case code := <-wd.ctrlChan:
			switch code {
			case stop:
				break wdLoop
			default:
				interval = internalInterval * time.Second
				continue
			}
		case <-time.After(interval):
			wd.log.Errorf("Watchdog (dummy): WATCHDOG TIMED OUT!!! SHOULD NOT HAPPEN.")
			interval = 2 * time.Second
		}
	}
	wd.log.Debugf("Watchdog (dummy): stopped")
	wd.ready = false
}

// Called by Maestro upon load. If the watchdog needs a setup procedure this
// should occur here. An error returned will prevent further calls from Maestro
// including KeepAlive()
func (wd *watchdog) Setup(config *maestroSpecs.WatchdogConfig, logger maestroSpecs.WatchdogLogger) (err error) {
	wd.log = logger
	wd.log.Debugf("Watchdog (dummy) Got Setup() call with config: %+v", config)
	wd.ctrlChan = make(chan int)
	return
}

// func SetupWatchdog(config *maestroSpecs.WatchdogConfig, logger maestroSpecs.WatchdogLogger) (err error, watchdog maestroSpecs.Watchdog) {
// 	wd = new(watchdog)
// 	wd.log = logger
// 	wd.ready = true
// 	wd.log.Debugf("Got Setup() call with config: %+v", config)
// 	wd.ctrlChan = make(chan int)
// 	return
// }

// CriticalInterval returns the time interval the watchdog *must* be called
// in order to keep the system alive. The implementer should build in enough buffer
// to this for the plugin to do its work of keep any hw / sw watchdog alive.
// Maestro will call KeepAlive() at this interval or quicker, but never slower.
func (wd *watchdog) CriticalInterval() time.Duration {
	// 5 seconds buffer please
	return time.Second * (internalInterval - 5)
}

// KeepAlive is called by Maestro to keep the watch dog up. A call of KeepAlive
// means the watchdog plug can safely assume all is normal on the system
func (wd *watchdog) KeepAlive() (err error) {
	select {
	case wd.ctrlChan <- keepalive:
		wd.log.Debugf("Watchdog (dummy): KeepAlive() called")
	default:
		wd.log.Errorf("Watchdog (dummy): KeepAlive() called - but would block. wd not running?")
	}
	return
}

// NotOk is called by Maestro at the same interval as KeepAlive, but in leiu of it
// when the system is not meeting the criteria to keep the watch dog up
func (wd *watchdog) NotOk() (err error) {
	wd.log.Errorf("Watchdog (dummy) Not OK.")
	return
}

// Disabled is called when Maestro desires, typically from a command it recieved, to
// disabled the watchdog. Usually for debugging. Implementors should disable the watchdog
// or if this is impossible, they can return an error.
func (wd *watchdog) Disable() (err error) {
	if wd.ready {
		select {
		case wd.ctrlChan <- stop:
			wd.log.Debugf("Watchdog (dummy): Disable() called")
		default:
			wd.log.Errorf("Watchdog (dummy): Disable() called - but would block. wd not running?")
		}
	}
	return
}

// Enable() is not called normally. But if Disable() is called, the a call of Enable()
// should renable the watchdog if it was disabled.
func (wd *watchdog) Enable() (err error) {
	wd.log.Debugf("Watchdog (dummy): Enable() called")
	go wd.watchdogRunner()
	return
}
