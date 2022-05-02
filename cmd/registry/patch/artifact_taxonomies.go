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

const TaxonomyListMimeType = "application/octet-stream;type=google.cloud.apigeeregistry.v1.apihub.TaxonomyList"

type TaxonomyListData struct {
	DisplayName string     `yaml:"displayName,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Taxonomies  []Taxonomy `yaml:"taxonomies"`
}

func (a *TaxonomyListData) GetMimeType() string {
	return TaxonomyListMimeType
}

type Taxonomy struct {
	ID              string            `yaml:"id"`
	DisplayName     string            `yaml:"displayName,omitempty"`
	Description     string            `yaml:"description,omitempty"`
	AdminApplied    bool              `yaml:"adminApplied,omitempty"`
	SingleSelection bool              `yaml:"singleSelection,omitempty"`
	SearchExcluded  bool              `yaml:"searchExcluded,omitempty"`
	SystemManaged   bool              `yaml:"systemManaged,omitempty"`
	DisplayOrder    int               `yaml:"displayOrder"`
	Elements        []TaxonomyElement `yaml:"elements"`
}

type TaxonomyElement struct {
	ID          string `yaml:"id"`
	DisplayName string `yaml:"displayName,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func (l *TaxonomyListData) GetMessage() proto.Message {
	return &rpc.TaxonomyList{
		DisplayName: l.DisplayName,
		Description: l.Description,
		Taxonomies:  l.taxonomies(),
	}
}

func (l *TaxonomyListData) taxonomies() []*rpc.TaxonomyList_Taxonomy {
	taxonomies := make([]*rpc.TaxonomyList_Taxonomy, 0)
	for _, t := range l.Taxonomies {
		taxonomies = append(taxonomies,
			&rpc.TaxonomyList_Taxonomy{
				Id:              t.ID,
				DisplayName:     t.DisplayName,
				Description:     t.Description,
				AdminApplied:    t.AdminApplied,
				SingleSelection: t.SingleSelection,
				SearchExcluded:  t.SearchExcluded,
				SystemManaged:   t.SystemManaged,
				DisplayOrder:    int32(t.DisplayOrder),
				Elements:        t.elements(),
			},
		)
	}
	return taxonomies
}

func (t *Taxonomy) elements() []*rpc.TaxonomyList_Taxonomy_Element {
	elements := make([]*rpc.TaxonomyList_Taxonomy_Element, 0)
	for _, e := range t.Elements {
		elements = append(elements, &rpc.TaxonomyList_Taxonomy_Element{
			Id:          e.ID,
			DisplayName: e.DisplayName,
			Description: e.Description,
		})
	}
	return elements
}

func newTaxonomyList(message *rpc.Artifact) (*Artifact, error) {
	artifactName, err := names.ParseArtifact(message.Name)
	if err != nil {
		return nil, err
	}
	value := &rpc.TaxonomyList{}
	err = proto.Unmarshal(message.Contents, value)
	if err != nil {
		return nil, err
	}
	taxonomies := make([]Taxonomy, len(value.Taxonomies))
	for i, t := range value.Taxonomies {
		elements := make([]TaxonomyElement, len(t.Elements))
		for j, e := range t.Elements {
			elements[j] = TaxonomyElement{
				ID:          e.Id,
				DisplayName: e.DisplayName,
				Description: e.Description,
			}
		}
		taxonomies[i] = Taxonomy{
			ID:              t.Id,
			DisplayName:     t.DisplayName,
			Description:     t.Description,
			AdminApplied:    t.AdminApplied,
			SingleSelection: t.SingleSelection,
			SearchExcluded:  t.SearchExcluded,
			SystemManaged:   t.SystemManaged,
			DisplayOrder:    int(t.DisplayOrder),
			Elements:        elements,
		}
	}
	return &Artifact{
		Header: Header{
			ApiVersion: RegistryV1,
			Kind:       "TaxonomyList",
			Metadata: Metadata{
				Name: artifactName.ArtifactID(),
			},
		},
		Data: &TaxonomyListData{
			DisplayName: value.DisplayName,
			Description: value.Description,
			Taxonomies:  taxonomies,
		},
	}, nil
}