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

package times

import (
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/blockysource/go-pkg/times/internal/tzlocal"
)

var localTZ struct {
	once sync.Once
	loc  *time.Location
	err  error
}

// Local returns the local timezone with full name.
// In comparison to time.Local, this function returns the full name of the timezone.
// Example:
//
//	time.Local.String() -> "Local"
//	Local().String() -> "Europe/Berlin" (or whatever your local timezone is)
//
// This is useful in cases where a session is timezone oriented.
func Local() (*time.Location, error) {
	localTZ.once.Do(func() {
		tz, err := tzlocal.RuntimeTZ()
		if err != nil {
			localTZ.err = err
			return
		}
		lTZ, err := time.LoadLocation(tz)
		if err != nil {
			localTZ.err = err
		}
		localTZ.loc = lTZ
	})
	if localTZ.err != nil {
		return nil, localTZ.err
	}
	return localTZ.loc, nil
}

// LoadLocation loads the timezone with the given name.
func LoadLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}
