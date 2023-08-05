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
	"net"
	"sync"
	"time"
)

// ErrNoIP is returned when no IP could be found.
var ErrNoIP = errors.New("no IP could be found")

// DefaultConsensus returns a consensus filled
// with default and recommended HTTPSources.
// TLS-Protected providers get more power,
// compared to plain-text providers.
func DefaultConsensus() (*Consensus, error) {
	c, err := NewConsensus()
	if err != nil {
		return nil, err
	}

	// TLS-protected providers
	if err = c.AddVoter(NewHTTPSource("https://icanhazip.com/"), 3); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("https://myexternalip.com/raw"), 3); err != nil {
		return nil, err
	}

	// Plain-text providers
	if err = c.AddVoter(NewHTTPSource("http://ifconfig.io/ip"), 1); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("http://checkip.amazonaws.com/"), 1); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("http://ident.me/"), 1); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("http://whatismyip.akamai.com/"), 1); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("http://myip.dnsomatic.com/"), 1); err != nil {
		return nil, err
	}
	if err = c.AddVoter(NewHTTPSource("http://diagnostic.opendns.com/myip"), 1); err != nil {
		return nil, err
	}

	return c, nil
}

// NewConsensus creates a new Consensus, with no sources.
func NewConsensus() (*Consensus, error) {
	return &Consensus{timeout: time.Second * 30}, nil
}

// Consensus the type at the center of this library,
// and is the main entry point for users.
// Its `ExternalIP` method allows you to ask for your ExternalIP,
// influenced by all its added voters.
type Consensus struct {
	voters   []voter
	timeout  time.Duration
	protocol uint
}

// AddVoter adds a voter to this consensus.
// The source cannot be <nil> and
// the weight has to be of a value of 1 or above.
func (c *Consensus) AddVoter(source Source, weight uint) error {
	if source == nil {
		return errors.New("no sources provided")
	}
	if weight == 0 {
		return errors.New("weight cannot be 0")
	}

	c.voters = append(c.voters, voter{
		source: source,
		weight: weight,
	})
	return nil
}

// ResolveExternalIP resolves the externalIP from all added voters,
// returning the IP which received the most votes.
// The returned IP will always be valid, in case the returned error is <nil>.
func (c *Consensus) ResolveExternalIP(ctx context.Context) (net.IP, error) {
	voteCollection := make(map[string]uint)
	var vlock sync.Mutex
	var wg sync.WaitGroup

	// start all source Requests on a seperate goroutine
	for _, v := range c.voters {
		wg.Add(1)
		go func(v voter) {
			defer wg.Done()
			ip, err := v.source.IP(ctx, c.protocol)
			if err == nil && ip != nil {
				vlock.Lock()
				defer vlock.Unlock()
				voteCollection[ip.String()] += v.weight
			}
		}(v)
	}

	// wait for all votes to come in,
	// or until their process times out
	done := make(chan struct{}, 1)

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// if no votes were casted succesfully,
	// return early with an error
	if len(voteCollection) == 0 {
		return nil, ErrNoIP
	}

	var max uint
	var externalIP string

	// find the IP which has received the most votes,
	// influinced by the voter's weight.
	vlock.Lock()
	defer vlock.Unlock()
	for ip, votes := range voteCollection {
		if votes > max {
			max, externalIP = votes, ip
		}
	}

	// as the found IP was parsed previously,
	// we know it cannot be nil and is valid
	return net.ParseIP(externalIP), nil
}

// UseIPProtocol will set the IP Protocol to use for http requests
// to the sources. If zero, it will not discriminate. This is useful
// when you want to get the external IP in a specific protocol.
// Protocol only supports 0, 4 or 6.
func (c *Consensus) UseIPProtocol(protocol uint) error {
	if protocol != 0 && protocol != 4 && protocol != 6 {
		return errors.New("only ipv4 and ipv6 protocol is supported")
	}
	c.protocol = protocol
	return nil
}

// Source defines the part of a voter which gives the actual voting value (IP).
type Source interface {
	// IP returns IPv4/IPv6 address in a non-error case
	// net.IP should never be <nil> when error is <nil>
	// It is recommended that the IP function times out,
	// if no result could be found, after the given timeout duration.
	IP(ctx context.Context, protocol uint) (net.IP, error)
}

// voter adds weight to the IP given by a source.
// The weight has to be at least 1, and the more it is, the more power the voter has.
type voter struct {
	source Source // provides the IP (see: vote)
	weight uint   // provides the weight of its vote (acts as a multiplier)
}
