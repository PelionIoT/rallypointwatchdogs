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
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/armPelionEdge/maestroSpecs"
)

const (
	internalInterval = 90 // seconds
)

// This is a Go plugin, meeting the specs of maestroSpecs/watchdog.go

type watchdog struct {
	path        string // path to deviceOSWD
	upString    string
	upStringLen int
	ready       bool
	log         maestroSpecs.WatchdogLogger
	conn        net.Conn
	interval    int64
}

// Watchdog is the single instance of this watchdog
var _watchdog *watchdog

// Watchdog this is the exported interface of the watchdog
var Watchdog maestroSpecs.Watchdog

// InitMaestroPlugin is always called by Maestro any time it inits a plugin
func InitMaestroPlugin() (err error) {
	_watchdog = new(watchdog)
	Watchdog = _watchdog
	return
}

// Called by Maestro upon load. If the watchdog needs a setup procedure this
// should occur here. An error returned will prevent further calls from Maestro
// including KeepAlive()
func (wd *watchdog) Setup(config *maestroSpecs.WatchdogConfig, logger maestroSpecs.WatchdogLogger) (err error) {
	wd.log = logger
	if len(config.Opt1) < 1 {
		logger.Errorf("RP100 watchdog plugin needs 'opt1' as path to watchdog socket. Failing.")
		err = errors.New("need watchdog path")
	} else {
		wd.path = config.Opt1
		// Go has some non-standard names for stuff
		// See here: https://golang.org/src/net/unixsock_posix.go#L16
		// We need a 'unixgram' ... 'unixgram' ... how cute.
		wd.interval = internalInterval
		if len(config.Opt2) > 0 {
			v, err := strconv.ParseInt(config.Opt2, 10, 64)
			if err != nil {
				logger.Errorf("RP100 watchdog plugin: error on opt2 - conversion error - need interval in seconds: %s", err.Error())
			} else {
				if v < 1 {
					v = internalInterval
					logger.Errorf("RP100 watchdog plugin: error on opt2 - out of range. Need number > 1")
				} else {
					wd.interval = v
				}
			}
		}
		wd.conn, err = net.Dial("unixgram", wd.path)
		if err != nil {
			logger.Errorf("RP100 watchdog plugin failed to connect to unix dgram socket %s for watchdog. Failing.", wd.path)
		} else {
			// send an initial 'up' in order to test the socket
			wd.upString = fmt.Sprintf("up %d", wd.interval)
			wd.upStringLen = len(wd.upString)
			wd.ready = true
			cnt, err := wd.conn.Write([]byte(wd.upString))
			if err != nil {
				logger.Errorf("RP100 watchdog plugin - failed to write to socket. %s", err.Error())
			} else {
				if cnt < wd.upStringLen {
					logger.Errorf("RP100 watchdog plugin - failed to write add data to socket.")
				}
			}
		}
	}
	return
}

// CriticalInterval returns the time interval the watchdog *must* be called
// in order to keep the system alive. The implementer should build in enough buffer
// to this for the plugin to do its work of keep any hw / sw watchdog alive.
// Maestro will call KeepAlive() at this interval or quicker, but never slower.
func (wd *watchdog) CriticalInterval() time.Duration {
	// 15 seconds please
	return time.Second * (time.Duration(wd.interval) - 15)
}

// KeepAlive is called by Maestro to keep the watch dog up. A call of KeepAlive
// means the watchdog plug can safely assume all is normal on the system
func (wd *watchdog) KeepAlive() (err error) {
	if wd.ready {
		cnt, err := wd.conn.Write([]byte(wd.upString))
		if err != nil {
			wd.log.Errorf("RP100 watchdog plugin - failed to write to socket. %s", err.Error())
		} else {
			if cnt < wd.upStringLen {
				wd.log.Errorf("RP100 watchdog plugin - failed to write add data to socket.")
			} else {
				wd.log.Debugf("RP100 watchdog plugin: wrote \"%s\"", wd.upString)
			}

		}
	}
	return
}

// NotOk is called by Maestro at the same interval as KeepAlive, but in leiu of it
// when the system is not meeting the criteria to keep the watch dog up
func (wd *watchdog) NotOk() (err error) {
	wd.log.Errorf("RP100 watchdog plugin - Skipping write. Not OK.")
	return
}

// Disabled is called when Maestro desires, typically from a command it recieved, to
// disabled the watchdog. Usually for debugging. Implementors should disable the watchdog
// or if this is impossible, they can return an error.
func (wd *watchdog) Disable() (err error) {
	// TODO need to disable deviceOSWD
	wd.ready = false
	return
}

// Enable() is not called normally. But if Disable() is called, the a call of Enable()
// should renable the watchdog if it was disabled.
func (wd *watchdog) Enable() (err error) {
	// TODO need to enable deviceOSWD
	return
}
