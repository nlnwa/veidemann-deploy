/*
 * Copyright 2019 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pkg

import (
	"context"
	"fmt"
	dnsresolverV1 "github.com/nlnwa/veidemann-api-go/dnsresolver/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

/**
 * Server mocks
 */
type DnsResolverMock struct {
	ln            net.Listener
	listenAddr    net.Addr
	contextDialer grpc.DialOption
	Server        *grpc.Server
}

func NewDnsResolverMock(port string) (*DnsResolverMock, error) {
	m := &DnsResolverMock{}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to start Dns Resolver Server: %w", err)
	}

	m.ln = ln
	m.listenAddr = ln.Addr()

	opts := []grpc.ServerOption{}
	m.Server = grpc.NewServer(opts...)
	dnsresolverV1.RegisterDnsResolverServer(m.Server, m)

	go func() {
		log.Infof("Dns Resolver Server listening on address: %s", m.listenAddr)
		if err := m.Server.Serve(m.ln); err != nil {
			log.Fatalf("Dns Resolver Server exited with error: %v", err)
		}
	}()
	return m, nil
}

func (m *DnsResolverMock) Close() {
	m.Server.GracefulStop()
	m.ln.Close()
}

// Resolve implements DnsResolverServer
func (m *DnsResolverMock) Resolve(ctx context.Context, request *dnsresolverV1.ResolveRequest) (*dnsresolverV1.ResolveReply, error) {
	var n int
	fmt.Sscanf(request.Host, "foo%d.bar", &n)
	log.Tracef("DNS: %v", request.Host)

	return &dnsresolverV1.ResolveReply{
		Host:      request.Host,
		Port:      request.Port,
		TextualIp: fmt.Sprintf("10.0.0.%02d", n),
		RawIp:     []byte{10, 0, 0, byte(n)},
	}, nil
}
