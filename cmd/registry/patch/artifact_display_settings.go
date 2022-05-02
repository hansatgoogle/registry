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
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"google.golang.org/protobuf/proto"
)

const DisplaySettingsMimeType = "application/octet-stream;type=google.cloud.apigeeregistry.v1.apihub.DisplaySettings"

type DisplaySettingsData struct {
	Description     string `yaml:"description,omitempty"`
	Organization    string `yaml:"organization"`
	ApiGuideEnabled bool   `yaml:"apiGuideEnabled"`
	ApiScoreEnabled bool   `yaml:"apiScoreEnabled"`
}

func (l *DisplaySettingsData) GetMimeType() string {
	return DisplaySettingsMimeType
}

func (l *DisplaySettingsData) GetMessage() proto.Message {
	return &rpc.DisplaySettings{
		Description:     l.Description,
		Organization:    l.Organization,
		ApiGuideEnabled: l.ApiGuideEnabled,
		ApiScoreEnabled: l.ApiScoreEnabled,
	}
}

func newDisplaySettings(message *rpc.Artifact) (*Artifact, error) {
	artifactName, err := names.ParseArtifact(message.Name)
	if err != nil {
		return nil, err
	}
	value := &rpc.DisplaySettings{}
	err = proto.Unmarshal(message.Contents, value)
	if err != nil {
		return nil, err
	}
	return &Artifact{
		Header: Header{
			ApiVersion: RegistryV1,
			Kind:       "DisplaySettings",
			Metadata: Metadata{
				Name: artifactName.ArtifactID(),
			},
		},
		Data: &DisplaySettingsData{
			Description:     value.Description,
			Organization:    value.Organization,
			ApiGuideEnabled: value.ApiGuideEnabled,
			ApiScoreEnabled: value.ApiScoreEnabled,
		},
	}, nil
}