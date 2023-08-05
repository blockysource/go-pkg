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
	"time"
)

// Clock is an interface that gives access to the current time.
type Clock interface {
	// Now returns the current time.
	Now() time.Time

	// Unix returns the current time in unix format.
	Unix() int64

	// UnixNano returns the current time in unix format.
	UnixNano() int64

	// Location returns the current location.
	Location() *time.Location

	// Since returns the time elapsed since t.
	Since(t time.Time) time.Duration

	// Until returns the time until t.
	Until(t time.Time) time.Duration
}

// GetClockTimezone returns timezone of given clock.
// It is useful for wire dependency injection.
func GetClockTimezone(c Clock) *time.Location {
	return c.Location()
}
