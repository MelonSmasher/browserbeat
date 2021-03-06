// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package kubernetes

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "kubernetes", asset.ModuleFieldsPri, AssetKubernetes); err != nil {
		panic(err)
	}
}

// AssetKubernetes returns asset data.
// This is the base64 encoded gzipped contents of ../metricbeat/module/kubernetes.
func AssetKubernetes() string {
	return "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"
}
