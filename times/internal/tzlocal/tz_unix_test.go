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

//go:build !windows
// +build !windows

package tzlocal

import "testing"

func TestInferFromPathSuccess(t *testing.T) {
	tz, err := inferFromPath("/usr/share/zoneinfo/Asia/Tokyo")
	if err != nil {
		t.Errorf("got err=%d; want: nil", err)
	}
	want := "Asia/Tokyo"
	if tz != want {
		t.Errorf("got tz=%s; want: %s", tz, want)
	}
}
