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

import (
	"fmt"
	"os"
	"time"
)

// EnvTZ will return the TZ env value if it is set, go will revert any invalid timezone to UTC
func EnvTZ() (string, bool) {
	if name, ok := os.LookupEnv("TZ"); ok {
		// Go treats blank as UTC
		if name == "" {
			return "UTC", true
		}
		_, err := time.LoadLocation(name)
		// Go treats invalid as UTC
		if err != nil {
			return "UTC", true
		}
		return name, true
	}
	return "", false
}

// RuntimeTZ get the full timezone name of the local machine
func RuntimeTZ() (string, error) {
	// Get the timezone from the TZ env variable
	if name, ok := EnvTZ(); ok {
		return name, nil
	}
	// Get the timezone from the system file
	name, err := LocalTZ()
	if err != nil {
		err = fmt.Errorf("failed to get local machine timezone: %w", err)
		return "", err
	}

	return name, err
}
