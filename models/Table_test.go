package models

import (
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/stretchr/testify/assert"
)

func TestConvertToProtoTable(t *testing.T) {
	headers := []string{"Header 1", "Header 2"}

	columns := []string{"Column 1", "Column 2"}

	row := &TableRow{Cells: columns}
	rows := []*TableRow{row}

	tbl := &Table{
		Headers: &TableRow{
			Cells: headers,
		},
		Rows: rows,
	}

	p := tbl.ConvertToProtoTable()

	assert.Contains(t, p.Headers.Cells, "Header 1")
	assert.Contains(t, p.Headers.Cells, "Header 2")
	assert.Contains(t, p.Rows[0].Cells, "Column 1")
	assert.Contains(t, p.Rows[0].Cells, "Column 2")
}
func TestCreateFromProtoTable(t *testing.T) {
	headers := []string{"Header 1", "Header 2"}

	columns := []string{"Column 1", "Column 2"}

	row := &m.ProtoTableRow{Cells: columns}
	rows := []*m.ProtoTableRow{row}

	p := &m.ProtoTable{
		Headers: &m.ProtoTableRow{
			Cells: headers,
		},
		Rows: rows,
	}

	tbl := CreateTableFromProtoTable(p)

	assert.Contains(t, tbl.Headers.Cells, "Header 1")
	assert.Contains(t, tbl.Headers.Cells, "Header 2")
	assert.Contains(t, tbl.Rows[0].Cells, "Column 1")
	assert.Contains(t, tbl.Rows[0].Cells, "Column 2")
}
