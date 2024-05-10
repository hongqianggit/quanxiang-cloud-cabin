/*
Copyright 2022 QuanxiangCloud Authors
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

package elastic

import (
	"context"
	"github.com/olivere/elastic/v7"

	pkglogger "github.com/quanxiang-cloud/cabin/logger"
)

// Config config
type Config struct {
	Host     []string
	Username string
	Password string
	Log      bool
}

// NewClient new elasticsearch client
func NewClient(conf *Config, log pkglogger.AdaptedLogger, opts ...elastic.ClientOptionFunc) (*elastic.Client, error) {
	for _, host := range conf.Host {
		opts = append(opts, elastic.SetURL(host))
	}
	if conf.Username != "" && conf.Password != "" {
		opts = append(opts, elastic.SetBasicAuth(conf.Username, conf.Password))
	}
	opts = append(opts, elastic.SetErrorLog(newLogger(log)))
	if conf.Log {
		opts = append(opts, elastic.SetInfoLog(newLogger(log)))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(conf.Host[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = client.ElasticsearchVersion(conf.Host[0])
	if err != nil {
		return nil, err
	}

	return client, nil
}
