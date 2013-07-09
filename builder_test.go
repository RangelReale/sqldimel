package sqldimel

import (
	"testing"
	"time"
)

func TestBuildGenerate(t *testing.T) {
	b := NewBuilderProc("user", NewBuildProcessorDefault())
	b.Add("id", 1)
	b.Add("name", "Monte Marto")
	b.Add("dob", time.Now())
	b.Add("optional", nil)
	b.Add("weight", 80.2)

	b.Where("id = ? and weight > ?", 1, 70.2)

	if b.Output(UPDATE) != "UPDATE user SET id=?, name=?, dob=?, optional=?, weight=? WHERE id = ? and weight > ?" {
		t.Error("Invalid generated SQL")
	}

	if b.OutputParams(UPDATE)[1] != "Monte Marto" {
		t.Error("Invalid parameter order")
	}
}
