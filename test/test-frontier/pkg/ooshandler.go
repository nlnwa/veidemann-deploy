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
	"github.com/golang/protobuf/ptypes/empty"
	ooshandlerV1 "github.com/nlnwa/veidemann-api-go/ooshandler/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

/**
 * Server mocks
 */
type OosHandlerMock struct {
	ln            net.Listener
	listenAddr    net.Addr
	contextDialer grpc.DialOption
	Server        *grpc.Server
}

func NewOosHandlerMock(port string) (*OosHandlerMock, error) {
	m := &OosHandlerMock{}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to start Out of Scope Handler Server: %w", err)
	}

	m.ln = ln
	m.listenAddr = ln.Addr()

	opts := []grpc.ServerOption{}
	m.Server = grpc.NewServer(opts...)
	ooshandlerV1.RegisterOosHandlerServer(m.Server, m)

	go func() {
		log.Infof("Out of Scope Handler Server listening on address: %s", m.listenAddr)
		if err := m.Server.Serve(m.ln); err != nil {
			log.Fatalf("Out of Scope Handler Server exited with error: %v", err)
		}
	}()
	return m, nil
}

func (m *OosHandlerMock) Close() {
	m.Server.GracefulStop()
	m.ln.Close()
}

// SubmitUri implements OosHandlerServer
func (m *OosHandlerMock) SubmitUri(ctx context.Context, request *ooshandlerV1.SubmitUriRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
