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

package runtime

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/linkall-labs/cdk-go/config"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
)

type SinkConfigConstructor func() config.SinkConfigAccessor
type SinkConstructor func() connector.Sink

func RunSink(cfgCtor SinkConfigConstructor, sinkCtor SinkConstructor) {
	cfg := cfgCtor()
	sink := sinkCtor()
	err := runConnector(cfg, sink)
	if err != nil {
		log.Error("run sink error", map[string]interface{}{
			log.KeyError: err,
		})
		os.Exit(-1)
	}
}

type sinkWorker struct {
	cfg  config.SinkConfigAccessor
	sink connector.Sink
	wg   sync.WaitGroup
}

func newSinkWorker(cfg config.SinkConfigAccessor, sink connector.Sink) Worker {
	return &sinkWorker{
		cfg:  cfg,
		sink: sink,
	}
}

func (w *sinkWorker) Start(ctx context.Context) error {
	port := w.cfg.GetPort()
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "failed to listen port: %d", port)
	}
	ceClient, err := ce.NewClientHTTP(ce.WithListener(ls))
	if err != nil {
		return errors.Wrap(err, "failed to init cloudevents client")
	}
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		err = ceClient.StartReceiver(ctx, w.receive)
		if err != nil {
			panic(fmt.Sprintf("failed to start cloudevnets receiver: %s", err))
		}
	}()
	log.Info("the connector started", map[string]interface{}{
		log.ConnectorName: w.sink.Name(),
		"listening":       port,
	})
	return nil
}

func (w *sinkWorker) receive(ctx context.Context, event ce.Event) ce.Result {
	result := w.sink.Arrived(ctx, &event)
	return result.ConvertToCeResult()
}

func (w *sinkWorker) Stop() error {
	w.wg.Wait()
	return w.sink.Destroy()
}
