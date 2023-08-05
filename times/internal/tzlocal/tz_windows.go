// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tzlocal

//go:generate go run ./../tzlocal/cmd/update_tzmapping.go

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const tzKey = `SYSTEM\CurrentControlSet\Control\TimeZoneInformation`
const tzKeyVal = "TimeZoneKeyName"

// LocalTZ obtains the name of the time zone Windows is configured to use. Returns the corresponding IANA standard name
func LocalTZ() (string, error) {
	var winTZname string
	var errTzutil, errReg error

	// try tzutil command first - if that is not available, try to read from registry
	winTZname, errTzutil = localTZfromTzutil()
	if errTzutil != nil {
		winTZname, errReg = localTZfromReg()
		if errReg != nil { // both methods failed, return both errors
			return "", fmt.Errorf("failed to read time zone name with errors\n(1) %s\n(2) %s", errTzutil, errReg)
		}
	}

	if name, ok := WinTZtoIANA[winTZname]; ok {
		return name, nil
	}
	return "", fmt.Errorf("could not find IANA tz name for set time zone \"%s\"", winTZname)
}

// localTZfromTzutil executes command `tzutil /g` to get the name of the time zone Windows is configured to use.
func localTZfromTzutil() (string, error) {
	cmd := exec.Command("tzutil", "/g")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// localTZfromReg obtains the time zone Windows is configured to use from registry.
func localTZfromReg() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, tzKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	winTZname, _, err := k.GetStringValue(tzKeyVal)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(winTZname), nil
}
