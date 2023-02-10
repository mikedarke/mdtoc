package toc

import "testing"

func TestIsExcludedFile(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		excludedList   []string
		expectExcluded bool
	}{
		{
			name:           "if filename is in the excluded list",
			filename:       "index.md",
			excludedList:   []string{"index.md", "test.md"},
			expectExcluded: true,
		},
		{
			name:           "if filename is not in the excluded list",
			filename:       "index.md",
			excludedList:   []string{"index1.md", "test.md"},
			expectExcluded: false,
		},
		{
			name:           "if excluded list is empty",
			filename:       "index.md",
			excludedList:   []string{},
			expectExcluded: false,
		},
		{
			name:           "if filename is part of an excluded file",
			filename:       "index.md",
			excludedList:   []string{"index.mdf"},
			expectExcluded: false,
		},
		{
			name:           "if filename is part of an excluded file",
			filename:       "index.md",
			excludedList:   []string{"findex.md"},
			expectExcluded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := isExcludedFile(tt.filename, tt.excludedList)
			if is != tt.expectExcluded {
				t.Errorf("Expected [%t] but is [%t]", tt.expectExcluded, is)
			}
		})
	}
}

func TestGetParent(t *testing.T) {
	tests := []struct {
		name               string
		active             Heading
		previous           Heading
		expectedParentText string
	}{
		{
			name: "1. Finds the parent when it's the previous heading",
			active: Heading{
				depth: 2,
			},
			previous: Heading{
				depth: 1,
				text:  "A",
			},
			expectedParentText: "A",
		},
		{
			name: "2. Finds the parent when it's the previous heading parent",
			active: Heading{
				depth: 2,
			},
			previous: Heading{
				depth: 2,
				text:  "A",
				parent: &Heading{
					text:  "B",
					depth: 1,
				},
			},
			expectedParentText: "B",
		},
		{
			name: "3. Finds the parent when it's the previous headings parent's parent",
			active: Heading{
				depth: 2,
			},
			previous: Heading{
				depth: 3,
				text:  "A",
				parent: &Heading{
					text:  "B",
					depth: 2,
					parent: &Heading{
						depth: 1,
						text:  "C",
					},
				},
			},
			expectedParentText: "C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getParent(tt.active, &tt.previous)
			if p.text != tt.expectedParentText {
				t.Errorf("expected parent to be %s but have %s", tt.expectedParentText, p.text)
			}
		})
	}
}
