package toc

import "testing"

func TestFindNearestSectionAnchor(t *testing.T) {
	tests := []struct {
		name             string
		sections         []SectionMarker
		linenumber       int
		expectedDistance int
		expectedID       string
	}{
		{
			name: "1. Finds nearest section",
			sections: []SectionMarker{
				{id: "1", linenumber: 10},
				{id: "2", linenumber: 100},
			},
			linenumber:       20,
			expectedDistance: 10,
			expectedID:       "1",
		},
		{
			name: "2. Finds nearest section",
			sections: []SectionMarker{
				{id: "1", linenumber: 99},
				{id: "2", linenumber: 100},
			},
			linenumber:       101,
			expectedDistance: 1,
			expectedID:       "2",
		},
		{
			name: "3. Finds nearest section",
			sections: []SectionMarker{
				{id: "1", linenumber: 99},
				{id: "2", linenumber: 100},
				{id: "3", linenumber: 102},
			},
			linenumber:       101,
			expectedDistance: 1,
			expectedID:       "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, d := findNearestSectionAnchor(tt.linenumber, tt.sections)
			if d != tt.expectedDistance {
				t.Errorf("expected distance [%d] but have [%d]", tt.expectedDistance, d)
			}

			if m.id != tt.expectedID {
				t.Errorf("expected section ID to be [%s] but is [%s]", tt.expectedID, m.id)
			}
		})
	}
}
