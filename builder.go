// Copyright 2013 Rangel Reale. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sqldimel provides a SQL DML query builder.
package sqldimel

import (
	"bytes"
	"container/list"
	"database/sql"
	"fmt"
	"strings"
)

// DMLType is a custom type to identify the type of DML to generate
// (INSERT, UPDATE or DELETE)
type DMLType int

const (
	INSERT DMLType = iota
	UPDATE
	DELETE
)

// Builder is the generation class
type Builder struct {
	table     string
	fields    *list.List
	where     string
	whereargs []interface{}
	processor BuilderProcessor
}

type field struct {
	name  string
	value interface{}
}

// Creates a new builder for the specified table, with a default processor
func NewBuilder(table string) *Builder {
	b := Builder{
		table:     table,
		processor: &BuildProcessorDefault{},
	}
	b.fields = list.New()
	return &b
}

// Creates a new builder for the specified table and processor
func NewBuilderProc(table string, processor BuilderProcessor) *Builder {
	b := Builder{
		table:     table,
		processor: processor,
	}
	b.fields = list.New()
	return &b
}

// Add a field and value to the builder
func (b *Builder) Add(fieldname string, value interface{}) {
	b.fields.PushBack(&field{name: fieldname, value: value})
}

// Add a query string to be put on the WHERE part if needed by the DML
func (b *Builder) Where(query string, args ...interface{}) {
	b.where = query
	b.whereargs = args
}

// Returns the SQL and parameters at the same time
func (b *Builder) OutputAll(dmltype DMLType) (string, []interface{}) {
	return b.Output(dmltype), b.OutputParams(dmltype)
}

// Returns the generated SQL. It should be used together with OutputParams
// to execute on the database
func (b *Builder) Output(dmltype DMLType) string {
	switch dmltype {
	case INSERT:
		return b.buildInsert()
	case UPDATE:
		return b.buildUpdate()
	case DELETE:
		return b.buildDelete()
	}
	return ""
}

// Returns the parameters to be passed to the database execution
func (b *Builder) OutputParams(dmltype DMLType) []interface{} {
	rlen := 0
	if dmltype != DELETE {
		rlen += b.fields.Len()
	}
	if dmltype != INSERT {
		rlen += len(b.whereargs)
	}

	ret := make([]interface{}, rlen)

	rct := 0

	if dmltype != DELETE {
		for f := b.fields.Front(); f != nil; f = f.Next() {
			ret[rct] = f.Value.(*field).value
			rct++
		}
	}

	if dmltype != INSERT {
		for _, w := range b.whereargs {
			ret[rct] = w
			rct++
		}
	}

	return ret
}

// Generates the WHERE part
func (b *Builder) outputWhere() string {
	if b.where == "" {
		return ""
	}

	var ret bytes.Buffer

	ret.WriteString(" WHERE ")

	inside := false
	var insiderune rune
	for _, c := range b.where {
		if strings.ContainsRune("'\"", c) {
			if !inside || insiderune == c {
				inside = !inside
				insiderune = c
			}
		}

		if !inside && strings.ContainsRune("?", c) {
			ret.WriteString(b.processor.NextParam(""))
		} else {
			ret.WriteRune(c)
		}
	}

	return ret.String()
}

// Build INSERT DML
func (b *Builder) buildInsert() string {
	var ret bytes.Buffer

	ret.WriteString(fmt.Sprintf("INSERT INTO %s (", b.table))
	first := true
	for f := b.fields.Front(); f != nil; f = f.Next() {
		if !first {
			ret.WriteString(", ")
		} else {
			first = false
		}
		ret.WriteString(f.Value.(*field).name)
	}

	ret.WriteString(") VALUES (")
	first = true
	b.processor.BeginParams()
	for f := b.fields.Front(); f != nil; f = f.Next() {
		if !first {
			ret.WriteString(", ")
		} else {
			first = false
		}
		ret.WriteString(b.processor.NextParam(f.Value.(*field).name))
	}
	ret.WriteString(")")

	return ret.String()
}

// Build UPDATE DML
func (b *Builder) buildUpdate() string {
	var ret bytes.Buffer

	ret.WriteString(fmt.Sprintf("UPDATE %s SET ", b.table))
	first := true
	b.processor.BeginParams()
	for f := b.fields.Front(); f != nil; f = f.Next() {
		if !first {
			ret.WriteString(", ")
		} else {
			first = false
		}
		ret.WriteString(f.Value.(*field).name)
		ret.WriteString("=")
		ret.WriteString(b.processor.NextParam(f.Value.(*field).name))
	}

	ret.WriteString(b.outputWhere())

	return ret.String()
}

// Build DELETE DML
func (b *Builder) buildDelete() string {
	var ret bytes.Buffer

	ret.WriteString(fmt.Sprintf("DELETE FROM %s", b.table))
	b.processor.BeginParams()
	ret.WriteString(b.outputWhere())

	return ret.String()
}

// Execute the current SQL on the database
func (b *Builder) Exec(db *sql.DB, dmltype DMLType) (res sql.Result, err error) {
	res, err = db.Exec(b.Output(dmltype), b.OutputParams(dmltype))
	return
}

// Execute the current SQL on the transaction
func (b *Builder) ExecTx(tx *sql.Tx, dmltype DMLType) (res sql.Result, err error) {
	res, err = tx.Exec(b.Output(dmltype), b.OutputParams(dmltype))
	return
}
