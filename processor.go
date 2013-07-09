// Copyright 2013 Rangel Reale. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sqldimel provides a SQL DML query builder.
package sqldimel

import (
	"fmt"
)

// BuilderProcessor generates parameter names compatible with the database
type BuilderProcessor interface {
	// Initialize parameter generation
	BeginParams()
	// Generates the next field name
	NextParam(fieldname string) string
}

// BuildProcessorDefault generates parameters using the character "?"
func NewBuildProcessorDefault() *BuildProcessorDefault {
	return &BuildProcessorDefault{}
}

type BuildProcessorDefault struct {
}

func (p *BuildProcessorDefault) BeginParams() {

}

func (p *BuildProcessorDefault) NextParam(fieldname string) string {
	return "?"
}

// BuildProcessorNumeric generates parameters using the character "$"
// followed by a sequential number starting with 1
func NewBuildProcessorNumeric() *BuildProcessorNumeric {
	return &BuildProcessorNumeric{}
}

type BuildProcessorNumeric struct {
	ct int
}

func (p *BuildProcessorNumeric) BeginParams() {
	p.ct = 1
}

func (p *BuildProcessorNumeric) NextParam(fieldname string) (ret string) {
	ret = fmt.Sprintf("$%d", p.ct)
	p.ct++
	return
}
