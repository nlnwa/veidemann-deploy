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
	robotsV1 "github.com/nlnwa/veidemann-api-go/robotsevaluator/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

/**
 * Server mocks
 */
type RobotsEvaluatorMock struct {
	ln            net.Listener
	listenAddr    net.Addr
	contextDialer grpc.DialOption
	Server        *grpc.Server
}

func NewRobotsEvaluatorMock(port string) (*RobotsEvaluatorMock, error) {
	m := &RobotsEvaluatorMock{}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to start Robots Evaluator Server: %w", err)
	}

	m.ln = ln
	m.listenAddr = ln.Addr()

	opts := []grpc.ServerOption{}
	m.Server = grpc.NewServer(opts...)
	robotsV1.RegisterRobotsEvaluatorServer(m.Server, m)

	go func() {
		log.Infof("Robots Evaluator Server listening on address: %s", m.listenAddr)
		if err := m.Server.Serve(m.ln); err != nil {
			log.Fatalf("Robots Evaluator Server exited with error: %v", err)
		}
	}()
	return m, nil
}

func (m *RobotsEvaluatorMock) Close() {
	m.Server.GracefulStop()
	m.ln.Close()
}

// IsAllowed implements RobotsEvaluatorServer
func (m *RobotsEvaluatorMock) IsAllowed(ctx context.Context, request *robotsV1.IsAllowedRequest) (*robotsV1.IsAllowedReply, error) {
	log.Tracef("ROBOTS: %v", request.Uri)
	return &robotsV1.IsAllowedReply{
		IsAllowed:   true,
		CrawlDelay:  0,
		CacheDelay:  0,
		Sitemap:     nil,
		OtherFields: nil,
	}, nil
}
