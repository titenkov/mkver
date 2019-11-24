package main

import (
	"gotest.tools/assert"
	"testing"
)

var VersionTests = []struct {
	name     string
	metadata Metadata
	template string
	expected string
	err      error
}{
	{"default", Metadata{Origin: "1.0.0"}, "{{.Origin}}", "1.0.0", nil},
}

func TestVersion(t *testing.T) {
	for _, test := range VersionTests {
		got, _ := New(test.metadata).Execute(test.template)
		assert.Equal(t, test.expected, got, "failed while testing "+test.name)
	}
}
