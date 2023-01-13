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

package store

import (
	"context"

	"github.com/linkall-labs/cdk-go/config"
	"github.com/pkg/errors"
)

var (
	ErrKeyNotExist = errors.New("key not exist")
)

type KVStore interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

var kvStore KVStore

func InitKvStore(cfg config.StoreConfig) (err error) {
	switch cfg.Type {
	case config.FileStore:
		kvStore, err = NewFileStore(cfg.StoreFile)
		if err != nil {
			return errors.Wrap(err, "new file store error")
		}
	case config.EtcdStore:
		kvStore, err = NewEtcdStore(cfg.Endpoints, cfg.KeyPrefix)
		if err != nil {
			return errors.Wrap(err, "new etcd store error")
		}
	default:
		kvStore = NewMemoryStore()
	}
	return nil
}

func GetKVStore() KVStore {
	return kvStore
}
