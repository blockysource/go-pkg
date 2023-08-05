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
var _ Clock = (*ZonedClock)(nil)

// ZonedClock is the default implementation of the Clock interface.
type ZonedClock struct {
	loc *time.Location
}

// NewLocalClock creates a new ZonedClock.
func NewLocalClock() (*ZonedClock, error) {
	loc, err := Local()
	if err != nil {
		return nil, err
	}
	return &ZonedClock{loc: loc}, nil
}

// NewZonedClock creates a new ZonedClock.
func NewZonedClock(loc *time.Location) *ZonedClock {
	return &ZonedClock{loc: loc}
}

// Now returns the current time.
func (c *ZonedClock) Now() time.Time {
	return time.Now().In(c.loc)
}

// Until returns the time until t.
func (c *ZonedClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

// Unix returns the current time in unix format.
func (c *ZonedClock) Unix() int64 {
	return time.Now().In(c.loc).Unix()
}

// UnixNano returns the current time in unix format.
func (c *ZonedClock) UnixNano() int64 {
	return time.Now().In(c.loc).UnixNano()
}

// Location returns the current location.
func (c *ZonedClock) Location() *time.Location {
	return c.loc
}

// Since returns the time elapsed since t.
func (c *ZonedClock) Since(t time.Time) time.Duration {
	return time.Now().In(c.loc).Sub(t)
}
