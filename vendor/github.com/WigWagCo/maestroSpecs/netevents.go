package maestroSpecs
//
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
//

const (
	// Net Events:
	// emitted when an interface state has changed
	// including new address, interface up / down
	InterfaceStateChange = "interface-state-change"
)

// NetEventData - All network events has this data struct
type NetEventData struct {
	Type      string              `json:"type"`
	Interface *InterfaceEventData `json:"interface"`
}

// InterfaceEventData is used if the event involves a change to an interface
type InterfaceEventData struct {
	// Id represents the interfaces, such as "eth0"
	ID        string `json:"id"`
	Num       int    `json:"num"`
	Address   string `json:"address"`
	LinkState string `json:"linkstate"`
}
