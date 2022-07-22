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

package config

import (
	"github.com/go-logr/logr"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/tidwall/gjson"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strings"
)

const (
	VanceConfigPathDv string = "/vance/config/config.json"
	VanceSecretPathDv string = "/vance/secret/secret.json"
	VanceSink         string = "v_target"
	VanceSinkDv       string = "http://localhost:8080"
	VancePort         string = "v_port"
	VancePortDv       string = "8080"
)

// ConfigAccessor provides an easy way to obtain configs
type ConfigAccessor struct {
	DefaultValues map[string]string
	Logger        logr.Logger
}

var Accessor ConfigAccessor
var userConfig map[string]string
var userSecret map[string]string

func init() {
	log.SetLogger(zap.New())
	Accessor = ConfigAccessor{
		DefaultValues: map[string]string{
			VanceSink: VanceSinkDv,
			VancePort: VancePortDv,
		},
		Logger: log.Log.WithName("ConfigAccessor"),
	}
	configPath := VanceConfigPathDv
	secretPath := VanceSecretPathDv
	userConfig = make(map[string]string)
	userSecret = make(map[string]string)
	content, err := os.ReadFile(configPath)

	if err != nil {
		Accessor.Logger.Info("READ user config failed", "configPath", configPath)
		content, err = os.ReadFile("./config.json")
		if err != nil {
			Accessor.Logger.Info("READ local config failed", "configPath", "./config.json")
		}
	}
	if len(content) != 0 {
		conf := gjson.ParseBytes(content).Map()
		Accessor.Logger.Info("conf length", "len", len(conf))

		for k, v := range conf {
			userConfig[k] = v.Str
		}
	}
	content, err = os.ReadFile(secretPath)
	if err != nil {
		Accessor.Logger.Info("READ user secret failed", "secretPath", secretPath)
		content, err = os.ReadFile("./secret.json")
		if err != nil {
			Accessor.Logger.Info("READ local secret failed", "secretPath", "./secret.json")
		}
	}
	if len(content) != 0 {
		conf := gjson.ParseBytes(content).Map()
		Accessor.Logger.Info("secret length", "len", len(conf))

		for k, v := range conf {
			userSecret[k] = v.Str
		}
	}
}

// GetString method retrieves by following steps:
// 1. Try to get an environment value by the key
// 2. Try to get the value from a user-specific json config file.
// Use config.Accessor.Get(key) to get any config value the user pass to the program
func (a *ConfigAccessor) GetString(key string) string {
	var ret string
	ret, existed := os.LookupEnv(strings.ToUpper(key))
	if !existed {
		a.Logger.Info("userConfig length", "len", len(userConfig))
		ret = userConfig[key]
	}
	return ret
}

func (a *ConfigAccessor) GetSecret(key string) string {
	var ret string
	ret, existed := os.LookupEnv(strings.ToUpper(key))
	if !existed {
		a.Logger.Info("userSecret length", "len", len(userSecret))
		ret = userSecret[key]
	}
	return ret
}

func (a *ConfigAccessor) getOrDefault(key string) string {
	ret := a.GetString(key)
	if ret == "" {
		ret = a.DefaultValues[key]
	}
	return ret
}

func (a *ConfigAccessor) VanceSink() string {
	return a.getOrDefault(VanceSink)
}
func (a *ConfigAccessor) VancePort() string {
	return a.getOrDefault(VancePort)
}
