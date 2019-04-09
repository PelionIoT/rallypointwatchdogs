package templates
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
	"github.com/WigWagCo/mustache"
	"strconv"
)

const ARCH_PREFIX = "ARCH_"

type TemplateVarDictionary struct {
	Map map[string]string
}

func NewTemplateVarDictionary() (ret *TemplateVarDictionary) {
	ret = new(TemplateVarDictionary)
	ret.Map = make(map[string]string)
	return
}

func (this *TemplateVarDictionary) Add(key string, value string) {
	this.Map[key] = value
}

func (this *TemplateVarDictionary) AddArch(key string, value string) {
	this.Map[ARCH_PREFIX+key] = value
}

func (this *TemplateVarDictionary) Del(key string) {
	delete(this.Map,key)
}

func (this *TemplateVarDictionary) Get(key string) (ret string, ok bool) {
	ret, ok = this.Map[key]
	if ok {
		ret = mustache.Render(ret,this.Map)
	}
	return
}

func (this *TemplateVarDictionary) Render(input string) (output string) {
	output = mustache.Render(input,this.Map)
	// do it again, to handle meta vars which hand meta-vars inside them
	output = mustache.Render(output,this.Map)
	output = mustache.Render(output,this.Map)
	return
}


const (
	TEMPLATEERROR_UNKNOWN = iota
	TEMPLATEERROR_BAD_FILE = iota
	TEMPLATEERROR_PERMISSIONS = iota
	TEMPLATEERROR_NO_TEMPLATE = iota
)

var job_error_map = map[int]string{
  0: "unknown",
  TEMPLATEERROR_BAD_FILE: "TEMPLATEERROR_BAD_FILE",
  TEMPLATEERROR_PERMISSIONS: "TEMPLATEERROR_PERMISSIONS",  
  TEMPLATEERROR_NO_TEMPLATE: "TEMPLATEERROR_NO_TEMPLATE",
}

type TemplateError struct {
	TemplateName string
	Code int
	ErrString string
}

func (this *TemplateError) Error() string {
	s := job_error_map[this.Code]
	return "JobError: " + this.TemplateName + ":"+s+" (" + strconv.Itoa(this.Code) +") " + this.ErrString
}



