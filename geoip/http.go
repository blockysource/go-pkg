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

package geoip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

// NewHTTPSource creates a HTTP Source object,
// which can be used to request the (external) IP from.
// The Default HTTP Client will be used if no client is given.
func NewHTTPSource(url string) *HTTPSource {
	return &HTTPSource{url: url}
}

// HTTPSource is the default source, to get the external IP from.
// It does so by requesting the IP from a URL, via an HTTP GET Request.
type HTTPSource struct {
	url    string
	parser ContentParser
}

// ContentParser can be used to add a parser to an HTTPSource
// to parse the raw content returned from a website, and return the IP.
// Spacing before and after the IP will be trimmed by the Consensus.
type ContentParser func(string) (string, error)

// WithParser sets the parser value as the value to be used by this HTTPSource,
// and returns the pointer to this source, to allow for chaining.
func (s *HTTPSource) WithParser(parser ContentParser) *HTTPSource {
	s.parser = parser
	return s
}

// IP implements Source.IP
func (s *HTTPSource) IP(ctx context.Context, protocol uint) (net.IP, error) {
	// Define the GET method with the correct url,
	// setting the User-Agent to our library
	req, err := http.NewRequestWithContext(ctx, "GET", s.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-external-ip (github.com/glendc/go-external-ip)")

	// transport to avoid goroutine leak
	tr := &http.Transport{
		MaxIdleConns:      1,
		IdleConnTimeout:   3 * time.Second,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
			Control: func(network, address string, c syscall.RawConn) error {
				if protocol == 4 && network == "tcp6" {
					return errors.New("rejecting ipv6 connection")
				} else if protocol == 6 && network == "tcp4" {
					return errors.New("rejecting ipv4 connection")
				}
				return nil
			},
		}).DialContext,
	}

	client := &http.Client{Transport: tr}

	// Do the request and read the body for non-error results.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// optionally parse the content
	raw := string(bytes)
	if s.parser != nil {
		raw, err = s.parser(raw)
		if err != nil {
			return nil, err
		}
	}

	// validate the IP
	externalIP := net.ParseIP(strings.TrimSpace(raw))
	if externalIP == nil {
		return nil, fmt.Errorf("returned an invalid IP: %s", raw)
	}

	// returned the parsed IP
	return externalIP, nil
}
