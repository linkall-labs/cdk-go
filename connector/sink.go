// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connector

import (
	"context"
	"fmt"

	ce "github.com/cloudevents/sdk-go/v2"
)

type Sink interface {
	Connector
	// Arrived event arrived
	Arrived(ctx context.Context, event ...*ce.Event) Result
}

type Code int
type Result struct {
	c   Code
	msg string
}

func NewResult(c Code, msg string) Result {
	return Result{c: c, msg: msg}
}

func (r Result) ConvertToCeResult() ce.Result {
	if r == Success {
		return nil
	}
	return ce.NewHTTPResult(int(r.c), r.msg)
}

func (r Result) GetCode() Code {
	return r.c
}

func (r Result) GetMsg() string {
	return r.msg
}

func (r Result) Error() error {
	return fmt.Errorf("{\"message\": \"%s\", \"code\": %d}", r.msg, r.c)
}

var (
	Success = NewResult(0, "success")
)
