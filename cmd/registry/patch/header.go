// Copyright 2022 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package patch

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const RegistryV1 = "apigeeregistry/v1"

type Header struct {
	ApiVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata"`
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

func readHeader(bytes []byte) (Header, error) {
	var header Header
	err := yaml.Unmarshal(bytes, &header)
	if err != nil {
		return header, err
	}
	if header.ApiVersion != RegistryV1 {
		return header, fmt.Errorf("unsupported API version: %s", header.ApiVersion)
	}
	return header, nil
}