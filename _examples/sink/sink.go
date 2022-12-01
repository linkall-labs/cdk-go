/*
Copyright 2022-Present The Vance Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/config"
)

type Config struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	Count            int `json:"count" yaml:"count"`
}

type ExampleSink struct {
	number int
}

func (s *ExampleSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	config := cfg.(*Config)
	s.number = config.Count
	return nil
}

func (s *ExampleSink) Name() string {
	return "ExampleSink"
}

func (s *ExampleSink) Destroy() error {
	return nil
}

func (s *ExampleSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		s.number++
		b, _ := json.Marshal(event)
		fmt.Println(s.number, string(b))
	}
	return cdkgo.SuccessResult
}

func main() {
	os.Setenv(config.EnvConfigFile, "./_examples/sink/config.yaml")
	cdkgo.RunSink(&Config{}, &ExampleSink{})
}
