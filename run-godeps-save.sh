#!/bin/bash

# Copyright (c) 2018, Arm Limited and affiliates.
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a s    ymlink
    SELF="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative s    ymlink, we need to resolve it relative to the path where the symlink file wa    s located
done

SELF="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

cd $SELF

$GOPATH/bin/godep save ./... "$@"
# if [ -d src ]; then
# 	mv src .src
# 	godep save ./... "$@"
# 	mv .src src
# fi

# godeps is a questionably designed tool.
# it - for instance - arbitrarily drops certain folder, for instance a 
# folder with no Go code, but with C dependencies. Hmm.

# rm -rf vendor/github.com/armPelionEdge/greasego/deps
# cp -a ../greasego/deps vendor/github.com/armPelionEdge/greasego
# cp -a ../greasego/src vendor/github.com/armPelionEdge/greasego
# cd -
