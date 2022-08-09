package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRenderStringTemplate(t *testing.T) {
	template := `
Plainmarkdown: {{ plainmarkdown .Text }}
Split: {{ $arr := split .Text " "}}{{ index $arr 3 }}
Trimspace: {{ trimspace .Text }}
Lower: {{ upper .Text }}
Upper: {{ lower .Text }}
Title: {{ title .Text }}
Prefixlines:
{{ prefixlines "  " .MultiLineTest }}
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
`
	result, err := renderStringTemplate("testTemplate", template, struct {
		Text          string
		MultiLineTest string
	}{
		Text: "my Odly cAsed striNg",
		MultiLineTest: `This text used
multiple lines`,
	})

	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(expectedString, result) {
		t.Errorf("expected: %+v, got: %+v", expectedString, result)
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

func TestExtractMetadata(t *testing.T) {
	for _, tt := range []struct {
		name      string
		delimiter string
		full      string
		expected  map[string]string
	}{
		{
			name:      "nometa",
			delimiter: ":",
			full:      "description",
			expected:  map[string]string{},
		},
		{
			name:      "full",
			delimiter: ":",
			full:      ":meta:subcategory:mysubcategory:This is a regular description.",
			expected: map[string]string{
				"subcategory": "mysubcategory",
			},
		},
		{
			name:      "full with different delimiter",
			delimiter: "!!!",
			full:      "!!!meta!!!subcategory!!!mysubcategory!!!This is a regular description.",
			expected: map[string]string{
				"subcategory": "mysubcategory",
			},
		},
		{
			name:      "full with crazy delimiter",
			delimiter: "!#!33<>",
			full:      "!#!33<>meta!#!33<>subcategory!#!33<>mysubcategory!#!33<>This is a regular description.",
			expected: map[string]string{
				"subcategory": "mysubcategory",
			},
		},
		{
			name:      "only metadata",
			delimiter: ":",
			full:      ":meta:subcategory:mysubcategory:",
			expected: map[string]string{
				"subcategory": "mysubcategory",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractMetadata(tt.full, tt.delimiter)

			if err != nil {
				t.Error(err)
			}

			if !cmp.Equal(result, tt.expected) {
				t.Errorf("expected: %+v, got: %+v", tt.expected, result)
			}
		})
	}
}
