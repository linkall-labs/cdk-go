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

package runtime

import (
	"context"

	"github.com/pkg/errors"

	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/linkall-labs/cdk-go/util"
)

func runConnector(cfg connector.ConnectorConfigAccessor, c connector.Connector) error {
	err := connector.ParseConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "init source config error")
	}
	ctx := util.SignalContext()
	err = c.Initialize(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "connector initialize failed")
	}
	worker := getWorker(cfg, c)
	err = worker.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "worker start failed")
	}
	select {
	case <-ctx.Done():
		log.Info("received system signal, beginning shutdown", map[string]interface{}{
			"name": c.Name(),
		})
		if err = worker.Stop(); err != nil {
			log.Error("worker stop fail", map[string]interface{}{
				log.KeyError: err,
				"name":       c.Name(),
			})
		} else {
			log.Info("connector shutdown graceful", map[string]interface{}{
				"name": c.Name(),
			})
		}
	}
	return nil
}

type Worker interface {
	Start(ctx context.Context) error
	Stop() error
}

func getWorker(cfg connector.ConnectorConfigAccessor, c connector.Connector) Worker {
	switch cfg.ConnectorType() {
	case connector.SourceConnector:
		return newSourceWorker(cfg.(connector.SourceConfigAccessor), c.(connector.Source))
	case connector.SinkConnector:
		return newSinkWorker(cfg.(connector.SinkConfigAccessor), c.(connector.Sink))
	}
	panic("unknown connector type:" + cfg.ConnectorType())
}