/*
 * Copyright 2020 National Library of Norway.
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
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// FrontierConn holds the client for the Frontier service
type FrontierConn struct {
	addr           string
	connectTimeout time.Duration
	dialOptions    []grpc.DialOption
	clientConn     *grpc.ClientConn
	client         frontierV1.FrontierClient
}

func NewFrontierConn(addr string) (*FrontierConn, error) {
	c := &FrontierConn{
		addr:           addr,
		connectTimeout: 5 * time.Second,
	}
	var err error

	// Set up frontierClient
	c.clientConn, err = c.connectService()
	if err != nil {
		return nil, err
	}
	c.client = frontierV1.NewFrontierClient(c.clientConn)
	log.Infof("Connected to frontier")
	return c, nil
}

//func (c *FrontierConn) Connect() error {
//	var err error
//
//	// Set up frontierClient
//	c.clientConn, err = c.connectService()
//	if err != nil {
//		return err
//	}
//	c.client = frontierV1.NewFrontierClient(c.clientConn)
//	log.Infof("Connected to frontier")
//
//	return nil
//}

func (c *FrontierConn) connectService() (*grpc.ClientConn, error) {
	log.Infof("Connecting Frontier at: %s", c.addr)

	dialOpts := append(c.dialOptions,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	dialCtx, dialCancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer dialCancel()

	return grpc.DialContext(dialCtx, c.addr, dialOpts...)
}

func (c *FrontierConn) Close() {
	_ = c.clientConn.Close()
}

func (c *FrontierConn) Client() frontierV1.FrontierClient {
	return c.client
}
