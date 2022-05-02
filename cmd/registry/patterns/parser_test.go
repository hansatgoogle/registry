package patterns

import (
	"testing"

	"github.com/apigee/registry/server/registry/names"
)

func generateSpec(t *testing.T, specName string) names.Spec {
	t.Helper()
	spec, err := names.ParseSpec(specName)
	if err != nil {
		t.Fatalf("Failed generateSpec(%s): %s", specName, err.Error())
	}
	return spec
}

func generateVersion(t *testing.T, versionName string) names.Version {
	t.Helper()
	version, err := names.ParseVersion(versionName)
	if err != nil {
		t.Fatalf("Failed generateVersion(%s): %s", versionName, err.Error())
	}
	return version
}

func generateArtifact(t *testing.T, artifactName string) names.Artifact {
	t.Helper()
	artifact, err := names.ParseArtifact(artifactName)
	if err != nil {
		t.Fatalf("Failed generateArtifact(%s): %s", artifactName, err.Error())
	}
	return artifact
}

func TestSubstituteReferenceEntity(t *testing.T) {
	tests := []struct {
		desc              string
		resourcePattern   string
		dependencyPattern string
		want              string
	}{
		{
			desc:              "artifact reference",
			resourcePattern:   "projects/demo/locations/global/apis/-/versions/-/specs/-/artifacts/lint-gnostic",
			dependencyPattern: "$resource.artifact",
			want:              "projects/demo/locations/global/apis/-/versions/-/specs/-/artifacts/lint-gnostic",
		},
		{
			desc:              "spec reference",
			resourcePattern:   "projects/demo/locations/global/apis/-/versions/-/specs/-/artifacts/-",
			dependencyPattern: "$resource.spec",
			want:              "projects/demo/locations/global/apis/-/versions/-/specs/-",
		},
		{
			desc:              "version reference",
			resourcePattern:   "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/-",
			dependencyPattern: "$resource.version/artifacts/lintstats",
			want:              "projects/demo/locations/global/apis/petstore/versions/1.0.0/artifacts/lintstats",
		},
		{
			desc:              "api reference",
			resourcePattern:   "projects/demo/locations/global/apis/-/versions/-/specs/-",
			dependencyPattern: "$resource.version/artifacts/lintstats",
			want:              "projects/demo/locations/global/apis/-/versions/-/artifacts/lintstats",
		},
		{
			desc:              "no reference",
			resourcePattern:   "projects/demo/locations/global/apis/-/artifacts/lintstats",
			dependencyPattern: "apis/-/versions/-",
			want:              "projects/demo/locations/global/apis/-/versions/-",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			resourceName, err := ParseResourcePattern(test.resourcePattern)
			if err != nil {
				t.Fatalf("Error in parsing, %s", err)
			}
			got, err := SubstituteReferenceEntity(test.dependencyPattern, resourceName)
			if err != nil {
				t.Errorf("SubstituteReferenceEntity returned unexpected error: %s", err)
			}
			if got.String() != test.want {
				t.Errorf("SubstituteReferenceEntity returned unexpected value want: %q got:%q", test.want, got)
			}
		})
	}
}

func TestSubstituteReferenceEntityError(t *testing.T) {
	tests := []struct {
		desc              string
		resourcePattern   string
		dependencyPattern string
	}{
		{
			desc:              "non-existent reference",
			resourcePattern:   "projects/demo/locations/global/apis/-/versions/-/specs/-",
			dependencyPattern: "$resource.artifact",
		},
		{
			desc:              "incorrect reference keyword",
			resourcePattern:   "projects/demo/locations/global/apis/-/versions/-/specs/-",
			dependencyPattern: "$resource.aip",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			resourceName, err := ParseResourcePattern(test.resourcePattern)
			if err != nil {
				t.Fatalf("Error in parsing, %s", err)
			}
			got, err := SubstituteReferenceEntity(test.dependencyPattern, resourceName)
			if err == nil {
				t.Errorf("expected SubstituteReferenceEntity to return error, got: %q", got.String())
			}
		})
	}
}

func TestFullResourceNameFromParent(t *testing.T) {
	tests := []struct {
		desc            string
		resourcePattern string
		parent          string
		want            ResourceName
	}{
		{
			desc:            "version pattern",
			resourcePattern: "projects/demo/locations/global/apis/-/versions/1.0.0",
			parent:          "projects/demo/locations/global/apis/petstore",
			want: VersionName{
				Name: generateVersion(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0"),
			},
		},
		{
			desc:            "spec pattern",
			resourcePattern: "projects/demo/locations/global/apis/-/versions/-/specs/openapi.yaml",
			parent:          "projects/demo/locations/global/apis/petstore/versions/1.0.0",
			want: SpecName{
				Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml"),
			},
		},
		{
			desc:            "artifact pattern",
			resourcePattern: "projects/demo/locations/global/apis/-/versions/-/specs/-/artifacts/complexity",
			parent:          "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
			want: ArtifactName{
				Name: generateArtifact(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/complexity"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := FullResourceNameFromParent(test.resourcePattern, test.parent)
			if err != nil {
				t.Errorf("FullResourceNameFromParent returned unexpected error: %s", err)
			}
			if got != test.want {
				t.Errorf("FullResourceNameFromParent returned unexpected value want: %q got:%q", test.want, got)
			}
		})
	}

}

func TestFullResourceNameFromParentError(t *testing.T) {
	tests := []struct {
		desc            string
		resourcePattern string
		parent          string
	}{
		{
			desc:            "incorrect keywords",
			resourcePattern: "projects/demo/locations/global/apis/-/versions/-/apispecs/-",
			parent:          "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
		},
		{
			desc:            "incorrect pattern",
			resourcePattern: "projects/demo/locations/global/apis/-/specs/-",
			parent:          "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := FullResourceNameFromParent(test.resourcePattern, test.parent)
			if err == nil {
				t.Errorf("expected FullResourceNameFromParent to return error, got: %q", got)
			}
		})
	}

}

func TestGetReferenceEntityValue(t *testing.T) {
	tests := []struct {
		desc            string
		resourcePattern string
		referred        ResourceName
		want            string
	}{
		{
			desc:            "api group",
			resourcePattern: "$resource.api/versions/-/specs/-",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
			want:            "projects/demo/locations/global/apis/petstore",
		},
		{
			desc:            "version group",
			resourcePattern: "$resource.version/specs/-",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
			want:            "projects/demo/locations/global/apis/petstore/versions/1.0.0",
		},
		{
			desc:            "spec group",
			resourcePattern: "$resource.spec",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
			want:            "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
		},
		{
			desc:            "artifact group",
			resourcePattern: "$resource.artifact",
			referred:        ArtifactName{Name: generateArtifact(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic")},
			want:            "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
		},
		{
			desc:            "no group",
			resourcePattern: "apis/-/versions/-/specs/-",
			referred:        ArtifactName{Name: generateArtifact(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic")},
			want:            "default",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := GetReferenceEntityValue(test.resourcePattern, test.referred)
			if err != nil {
				t.Errorf("GetReferenceEntityValue returned unexpected error: %s", err)
			}
			if got != test.want {
				t.Errorf("GetReferenceEntityValue returned unexpected value want: %q got:%q", test.want, got)
			}
		})
	}
}

func TestGetReferenceEntityValueError(t *testing.T) {
	tests := []struct {
		desc            string
		resourcePattern string
		referred        ResourceName
	}{
		{
			desc:            "typo",
			resourcePattern: "$resource.apis/versions/-/specs/-",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
		},
		{
			desc:            "incorrect reference",
			resourcePattern: "$resource.name/versions/-/specs/-",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
		},
		{
			desc:            "incorrect resourceKW",
			resourcePattern: "$resources.api/versions/-/specs/-",
			referred:        SpecName{Name: generateSpec(t, "projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml")},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := GetReferenceEntityValue(test.resourcePattern, test.referred)
			if err == nil {
				t.Errorf("expected GetReferenceEntityValue to return error, got: %q", got)
			}
		})
	}
}