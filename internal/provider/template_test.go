// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestRenderStringTemplate(t *testing.T) {
	t.Parallel()

	template := `
Plainmarkdown: {{ plainmarkdown .Text }}
Split: {{ $arr := split .Text " "}}{{ index $arr 3 }}
Trimspace: {{ trimspace .Text }}
Lower: {{ upper .Text }}
Upper: {{ lower .Text }}
Title: {{ title .Text }}
Prefixlines:
{{ prefixlines "  " .MultiLineTest }}
Printf tffile: {{ printf "{{tffile %q}}" .Code }}
tffile: {{ tffile .Code }}
`

	expectedString := `
Plainmarkdown: my Odly cAsed striNg
Split: striNg
Trimspace: my Odly cAsed striNg
Lower: MY ODLY CASED STRING
Upper: my odly cased string
Title: My Odly Cased String
Prefixlines:
  This text used
  multiple lines
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`
	result, err := renderStringTemplate("testdata/test-provider-dir", "testTemplate", template, struct {
		Text          string
		MultiLineTest string
		Code          string
	}{
		Text: "my Odly cAsed striNg",
		MultiLineTest: `This text used
multiple lines`,
		Code: "provider.tf",
	})

	if err != nil {
		t.Error(err)
	}
	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestResourceTemplate_Render(t *testing.T) {
	t.Parallel()

	template := `
Printf tffile: {{ printf "{{tffile %q}}" .ExampleFile }}
tffile: {{ tffile .ExampleFile }}
`
	expectedString := `
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`

	tpl := resourceTemplate(template)

	schema := tfjson.Schema{
		Version: 3,
		Block: &tfjson.SchemaBlock{
			Attributes:      nil,
			NestedBlocks:    nil,
			Description:     "",
			DescriptionKind: "",
			Deprecated:      false,
		},
	}

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "test-provider", "Resource", "provider.tf", []string{"provider.tf"}, "", "", "", &schema, nil, true, ":")
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestProviderTemplate_Render(t *testing.T) {
	t.Parallel()

	template := `
Printf tffile: {{ printf "{{tffile %q}}" .ExampleFile }}
tffile: {{ tffile .ExampleFile }}
`
	expectedString := `
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`

	tpl := providerTemplate(template)

	schema := tfjson.Schema{
		Version: 3,
		Block: &tfjson.SchemaBlock{
			Attributes:      nil,
			NestedBlocks:    nil,
			Description:     "",
			DescriptionKind: "",
			Deprecated:      false,
		},
	}

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "provider.tf", []string{"provider.tf"}, &schema)
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestExtractDescription(t *testing.T) {
	for _, tt := range []struct {
		name      string
		delimiter string
		full      string
		expected  string
	}{
		{
			name:      "nometa",
			delimiter: ":",
			full:      "description",
			expected:  "description",
		},
		{
			name:      "full",
			delimiter: ":",
			full:      ":meta:subcategory:mysubcategory:This is a regular description.",
			expected:  "This is a regular description.",
		},
		{
			name:      "full with different delimiter",
			delimiter: "!!!",
			full:      "!!!meta!!!subcategory!!!mysubcategory!!!This is a regular description.",
			expected:  "This is a regular description.",
		},
		{
			name:      "full with crazy delimiter",
			delimiter: "!#!33<>",
			full:      "!#!33<>meta!#!33<>subcategory!#!33<>mysubcategory!#!33<>This is a regular description.",
			expected:  "This is a regular description.",
		},
		{
			name:      "full with multine description",
			delimiter: ":",
			full: `:meta:subcategory:IP Address Management (IPAM):From the [official documentation](https://docs.netbox.dev/en/stable/core-functionality/ipam/#aggregates):

> NetBox allows us to specify the portions of IP space that are interesting to us by defining aggregates. Typically, an aggregate will correspond to either an allocation of public (globally routable) IP space granted by a regional authority, or a private (internally-routable) designation.`,
			expected: `From the [official documentation](https://docs.netbox.dev/en/stable/core-functionality/ipam/#aggregates):

> NetBox allows us to specify the portions of IP space that are interesting to us by defining aggregates. Typically, an aggregate will correspond to either an allocation of public (globally routable) IP space granted by a regional authority, or a private (internally-routable) designation.`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractDescription(tt.full, tt.delimiter)

			if err != nil {
				t.Error(err)
			}

			if result != tt.expected {
				t.Errorf("expected: %+v, got: %+v", tt.expected, result)
			}
		})
	}
}
