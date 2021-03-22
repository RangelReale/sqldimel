package sqldimel

import (
	"bytes"
	"container/list"
	"database/sql"
	"fmt"
)

type MultiBuilder struct {
	table     string
	fields    []string
	data      *list.List
	processor BuilderProcessor
}

// Creates a new multi builder
func NewMultiBuilder(table string, fields []string) *MultiBuilder {
	b := MultiBuilder{
		table:     table,
		fields:    fields,
		processor: &BuildProcessorDefault{},
	}
	b.data = list.New()
	return &b
}

// Creates a new multi builder for the specified table and processor
func NewMultiBuilderProc(table string, processor BuilderProcessor, fields []string) *MultiBuilder {
	b := MultiBuilder{
		table:     table,
		fields:    fields,
		processor: processor,
	}
	b.data = list.New()
	return &b
}

func (mb *MultiBuilder) HasData() bool {
	return mb.data.Len() > 0
}

func (mb *MultiBuilder) DataLen() int {
	return mb.data.Len()
}

func (mb *MultiBuilder) CreateData() *MultiBuilderData {
	d := NewMultiBuilderData()
	mb.data.PushBack(d)
	return d
}

func (mb *MultiBuilder) ClearData() {
	mb.data = list.New()
}

func (mb *MultiBuilder) Output() (string, []interface{}) {
	var ret bytes.Buffer
	params := make([]interface{}, 0)

	ret.WriteString(fmt.Sprintf("INSERT INTO %s (", mb.table))
	first := true
	for _, f := range mb.fields {
		if !first {
			ret.WriteString(", ")
		} else {
			first = false
		}
		ret.WriteString(f)
	}

	ret.WriteString(") VALUES ")
	first = true
	mb.processor.BeginParams()
	for d := mb.data.Front(); d != nil; d = d.Next() {
		if !first {
			ret.WriteString(", ")
		} else {
			first = false
		}
		ret.WriteString("(")
		firstfield := true
		for _, f := range mb.fields {
			if !firstfield {
				ret.WriteString(", ")
			} else {
				firstfield = false
			}
			ret.WriteString(mb.processor.NextParam(f))
			if v, ok := d.Value.(*MultiBuilderData).fields[f]; ok {
				params = append(params, v)
			} else {
				params = append(params, nil)
			}
		}
		ret.WriteString(")")
	}

	return ret.String(), params
}

// Execute the current SQL on the database
func (mb *MultiBuilder) Exec(db *sql.DB) (res sql.Result, err error) {
	q, qparams := mb.Output()
	res, err = db.Exec(q, qparams...)
	return
}

// Execute the current SQL on the transaction
func (mb *MultiBuilder) ExecTx(tx *sql.Tx) (res sql.Result, err error) {
	q, qparams := mb.Output()
	res, err = tx.Exec(q, qparams...)
	return
}

type MultiBuilderData struct {
	fields map[string]interface{}
}

func NewMultiBuilderData() *MultiBuilderData {
	return &MultiBuilderData{
		fields: make(map[string]interface{}),
	}
}

func (mbd *MultiBuilderData) Add(fieldname string, value interface{}) *MultiBuilderData {
	mbd.fields[fieldname] = value
	return mbd
}
