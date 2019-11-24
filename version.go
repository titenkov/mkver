// Copyright 2019 Pavel Titenkov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package mkver provides ability to format and enrich application version.
// In this package, semantic version strings must begin with a leading "v",
// as in "v1.0.0".

package main

import (
	"bytes"
	"log"
	"text/template"
)

// Metadata - container of version meta-information such as origin version, git branch, etc.
type Metadata struct {
	Origin    string
	GitBranch string
	GitSha    string
}

// Version is the representation of a processed version
type Version struct {
	metadata Metadata
	template string
}

// New allocates a new, undefined template with the given name.
func New(metadata Metadata) *Version {
	v := &Version{metadata: metadata}
	return v
}

// Execute - process template string and populate it with metadata
func (v *Version) Execute(versionTemplate string) (string, error) {
	v.template = versionTemplate
	return v.execute()
}

func (v *Version) execute() (string, error) {
	var tpl bytes.Buffer
	t := template.New("version")

	t, err := t.Parse(v.template)
	if err != nil {
		log.Fatal("Parsing error: ", err)
		return "", err
	}

	err = t.Execute(&tpl, v.metadata)
	if err != nil {
		log.Fatal("Execution error: ", err)
		return "", err
	}

	return tpl.String(), nil
}
