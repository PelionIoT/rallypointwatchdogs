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

import (
    "github.com/WigWagCo/maestroSpecs/templates"    
)

type PlatformReader interface {
    // implemented by a library which provides a way to pull off needed
    // platform specific variables, typically out of NVRAM, secure flash or EEPROM
    // use the Logger to log any errors or information
    GetPlatformVars(dict *templates.TemplateVarDictionary, log Logger) (err error)     
}
