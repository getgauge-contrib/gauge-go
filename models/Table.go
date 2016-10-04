package models

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
)

type Table struct {
	// / Contains the Headers for the table
	Headers *TableRow
	// / Contains the Rows for the table
	Rows []*TableRow
}

type TableRow struct {
	// / Represents the cells of a given table
	Cells []string
}

func (t *Table) ConvertToProtoTable() *m.ProtoTable {
	//TODO handle error scenarios
	p := &m.ProtoTable{
		Headers: &m.ProtoTableRow{
			Cells: t.Headers.Cells,
		},
	}

	for _, row := range t.Rows {
		p.Rows = append(p.Rows, &m.ProtoTableRow{Cells: row.Cells})
	}
	return p
}

func CreateTableFromProtoTable(p *m.ProtoTable) *Table {
	//TODO handle error scenarios
	t := &Table{
		Headers: &TableRow{
			Cells: p.Headers.Cells,
		},
	}

	for _, row := range p.Rows {
		t.Rows = append(t.Rows, &TableRow{Cells: row.Cells})
	}
	return t
}
