// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-defaults. DO NOT EDIT.

package test

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
)

var (
	_ = errors.New("")
	_ = io.EOF
	_ = context.Background()
)

var _ TestServer = (*proxyTest)(nil)

func NewTestProxy(c TestClient, opts ...grpc.CallOption) TestServer {
	return &proxyTest{c: c, opts: opts}
}

type proxyTest struct {
	c    TestClient
	opts []grpc.CallOption
}

// Unary proxies call to backend server
func (x *proxyTest) Unary(ctx context.Context, req *UnaryRequest) (*UnaryResponse, error) {
	return x.c.Unary(ctx, req, x.opts...)
}

// ClientStream proxies call to backend server
func (x *proxyTest) ClientStream(s Test_ClientStreamServer) error {
	cs, err := x.c.ClientStream(s.Context(), x.opts...)
	if err != nil {
		return err
	}
	for {
		r, err := s.Recv()
		if errors.Is(err, io.EOF) {
			res, err := cs.CloseAndRecv()
			if err != nil {
				return err
			}
			return s.SendAndClose(res)
		}
		if err != nil {
			return err
		}
		if err := cs.Send(r); err != nil {
			return err
		}
	}
}

// ServerStream proxies call to backend server
func (x *proxyTest) ServerStream(req *ServerStreamRequest, s Test_ServerStreamServer) error {
	cs, err := x.c.ServerStream(s.Context(), req, x.opts...)
	if err != nil {
		return err
	}
	for {
		res, err := cs.Recv()
		if err != nil {
			return err
		}
		if err := s.Send(res); err != nil {
			return err
		}
	}
}

// Stream proxies call to backend server
func (x *proxyTest) Stream(s Test_StreamServer) error {
	cs, err := x.c.Stream(s.Context(), x.opts...)
	if err != nil {
		return err
	}
	errs := make(chan error, 2)
	recv := func() error {
		for {
			req, err := s.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			if err := cs.Send(req); err != nil {
				return err
			}
		}
	}
	send := func() error {
		for {
			res, err := cs.Recv()
			if err != nil {
				return err
			}
			if err := s.Send(res); err != nil {
				return err
			}
		}
	}
	go func() {
		errs <- recv()
	}()
	go func() {
		errs <- send()
	}()
	return <-errs
}

func (x *proxyTest) mustEmbedUnimplementedTestServer() {}
