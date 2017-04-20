package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeDeps(t *testing.T) {
	tests := []struct {
		name         string
		deps         []Dependency
		expectedDeps []Dependency
	}{
		{
			name: "Simple",
			deps: []Dependency{
				Dependency{ImportPath: "github.com/enrichman/package_1/subpackage_1"},
			},
			expectedDeps: []Dependency{
				Dependency{ImportPath: "github.com/enrichman/package_1"},
			},
		},
		{
			name: "Multiple",
			deps: []Dependency{
				Dependency{ImportPath: "github.com/enrichman/package_1/subpackage_1"},
				Dependency{ImportPath: "github.com/enrichman/package_1/subpackage_2"},
				Dependency{ImportPath: "github.com/enrichman/package_1/subpackage_3"},
			},
			expectedDeps: []Dependency{
				Dependency{ImportPath: "github.com/enrichman/package_1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := mergeDeps(tt.deps)
			assert.Equal(t, tt.expectedDeps, merged)
		})
	}
}
