/*
 *
 * Copyright 2023 puzzlepassstrengthserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package passstrengthserver

import (
	"context"
	"errors"

	pb "github.com/dvaumoron/puzzlepassstrengthservice"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const PassstrengthKey = "puzzlePassstrength"

var errNotFound = errors.New("locale not found")

// server is used to implement puzzlesaltservice.SaltServer
type server struct {
	pb.UnimplementedPassstrengthServer
	minEntropy     float64
	localizedRules map[string]string
}

func New(defaultPass string, localizedRules map[string]string) pb.PassstrengthServer {
	return server{minEntropy: passwordvalidator.GetEntropy(defaultPass), localizedRules: localizedRules}
}

func (s server) GetRules(ctx context.Context, request *pb.LangRequest) (*pb.PasswordRules, error) {
	description, ok := s.localizedRules[request.Lang]
	if !ok {
		return nil, errNotFound
	}
	return &pb.PasswordRules{Description: description}, nil
}

func (s server) Check(ctx context.Context, request *pb.PasswordRequest) (*pb.Response, error) {
	err := passwordvalidator.Validate(request.Password, s.minEntropy)
	return &pb.Response{Success: err == nil}, nil
}
