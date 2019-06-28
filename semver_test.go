package main

import (
	"testing"
)

var tests = []struct {
    name string;       branch string;                   version string;                         expected string
} {
  { "",                "develop",                       "1.0.0-SNAPSHOT",                       "1.0.0-SNAPSHOT"                    },
  { "",                "feature/FX-12345",              "1.0.0-SNAPSHOT",                       "1.0.0-feature-fx-12345-SNAPSHOT"   },
  { "",                "defect/FX-12345",               "1.0.0-SNAPSHOT",                       "1.0.0-defect-fx-12345-SNAPSHOT"    },
  { "",                "release/1.0.0",                 "1.0.0",                                "1.0.0-rc0"                         },
  { "",                "hotfix/1.0.1",                  "1.0.1",                                "1.0.1-rc0"                         },
  { "",                "master",                        "1.0.0",                                "1.0.0"                             },
}

func TestCalculateSemanticVersion(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSemanticVersion(tt.branch, tt.version)
			if got != tt.expected {
				t.Errorf("CalculateSemanticVersion(%s, %s) got %v, want %v", tt.branch, tt.version, got, tt.expected)
			}
		})
	}
}