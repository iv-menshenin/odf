package mappers

import (
	"github.com/iv-menshenin/odf/mappers/attr"
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns/table"
	"github.com/kpmy/ypk/assert"
)

//Table structure holds the table structure and identifies it for TableMapper
type Table struct {
	Rows                int
	Columns             int
	Root                model.Leaf
	rowCache, colsCache []model.Leaf
	cellCache           [][]model.Leaf
}

//TableMapper writes and manages tables in document model
type TableMapper struct {
	List map[string]*Table
	fm   *Formatter
}

//Ready or not
func (t *TableMapper) Ready() bool {
	return t.fm != nil && t.fm.ready
}

func (t *TableMapper) newWriter(old ...model.Writer) model.Writer {
	return t.fm.m.NewWriter(old...)
}

//ConnectTo any valid Formatter and it's document model
func (t *TableMapper) ConnectTo(fm *Formatter) {
	t.fm = fm
	t.List = make(map[string]*Table)
}

//Write a Table with given name and dimensions, table object is stored internally and can be accessed by name, latest set TableAttributes, TableRowAttributes, TableColumnAttributes and TableCellAttributes are used
func (t *TableMapper) Write(name string, rows, cols int) {
	assert.For(t.Ready(), 20)
	assert.For(name != "" && t.List[name] == nil, 21)
	t.fm.attr.Flush()
	this := &Table{Rows: rows, Columns: cols}
	t.List[name] = this
	wr := t.newWriter()
	wr.Pos(t.fm.root)
	this.Root = wr.WritePos(New(table.Table))
	wr.Attr(table.Name, name)
	t.fm.attr.Fit(table.Table, func(a attr.Attributes) {
		wr.Attr(table.StyleName, a.Name())
	})
	for i := 0; i < this.Columns; i++ {
		col := New(table.TableColumn)
		this.colsCache = append(this.colsCache, col)
		this.cellCache = append(this.cellCache, make([]model.Leaf, 0))
		wr.Write(col)
		cwr := t.newWriter(wr)
		cwr.Pos(col)
		t.fm.attr.Fit(table.TableColumn, func(a attr.Attributes) {
			cwr.Attr(table.StyleName, a.Name())
		})
	}
	for i := 0; i < this.Rows; i++ {
		rwr := t.newWriter(wr)
		row := rwr.WritePos(New(table.TableRow))
		t.fm.attr.Fit(table.TableRow, func(a attr.Attributes) {
			rwr.Attr(table.StyleName, a.Name())
		})
		this.rowCache = append(this.rowCache, row)
		for j := 0; j < this.Columns; j++ {
			cell := New(table.TableCell)
			this.cellCache[j] = append(this.cellCache[j], cell)
			rwr.Write(cell)
			cwr := t.newWriter(rwr)
			cwr.Pos(cell)
			t.fm.attr.Fit(table.TableCell, func(a attr.Attributes) {
				cwr.Attr(table.StyleName, a.Name())
			})
		}
	}
}

//WriteRows to existing table latest set TableRowAttributes and TableCellAttributes are used
func (t *TableMapper) WriteRows(this *Table, rows int) {
	assert.For(t.Ready(), 20)
	t.fm.attr.Flush()
	wr := t.newWriter()
	for i := 0; i < rows; i++ {
		wr.Pos(this.Root)
		row := wr.WritePos(New(table.TableRow))
		t.fm.attr.Fit(table.TableRow, func(a attr.Attributes) {
			wr.Attr(table.StyleName, a.Name())
		})
		this.rowCache = append(this.rowCache, row)
		for j := 0; j < this.Columns; j++ {
			cell := New(table.TableCell)
			this.cellCache[j] = append(this.cellCache[j], cell)
			wr.Write(cell)
			cwr := t.newWriter(wr)
			cwr.Pos(cell)
			t.fm.attr.Fit(table.TableCell, func(a attr.Attributes) {
				cwr.Attr(table.StyleName, a.Name())
			})
		}
		this.Rows++
	}
}

//WriteColumns to existing table latest set TableColumnAttributes and TableCellAttributes are used
func (t *TableMapper) WriteColumns(this *Table, cols int) {
	assert.For(t.Ready(), 20)
	t.fm.attr.Flush()
	wr := t.newWriter()
	var last model.Leaf
	if this.Columns > 0 {
		last = this.colsCache[this.Columns-1]
	}
	for i := 0; i < cols; i++ {
		wr.Pos(this.Root)
		col := wr.WritePos(New(table.TableColumn), last)
		t.fm.attr.Fit(table.TableColumn, func(a attr.Attributes) {
			wr.Attr(table.StyleName, a.Name())
		})
		this.colsCache = append(this.colsCache, col)
		this.cellCache = append(this.cellCache, make([]model.Leaf, 0))
		this.Columns++
		for j := 0; j < this.Rows; j++ {
			t.WriteCells(this, j, 1)
		}
	}
}

//Write cells to existing table latest set TableCellAttributes are used
func (t *TableMapper) WriteCells(this *Table, _row int, cells int) {
	assert.For(t.Ready(), 20)
	t.fm.attr.Flush()
	wr := t.newWriter()
	row := this.rowCache[_row]
	wr.Pos(row)
	for i := 0; i < cells; i++ {
		cell := New(table.TableCell)
		this.cellCache[i] = append(this.cellCache[i], cell)
		wr.Write(cell)
		cwr := t.newWriter(wr)
		cwr.Pos(cell)
		t.fm.attr.Fit(table.TableCell, func(a attr.Attributes) {
			cwr.Attr(table.StyleName, a.Name())
		})
	}
}

//Span merges visually
func (t *TableMapper) Span(this *Table, row, col int, rowspan, colspan int) {
	assert.For(t.Ready(), 20)
	assert.For(rowspan > 0, 21)
	assert.For(colspan > 0, 22)
	wr := t.newWriter()
	wr.Pos(this.cellCache[col][row])
	wr.Attr(table.NumberRowsSpanned, rowspan)
	wr.Attr(table.NumberColumnsSpanned, colspan)
}

//Pos sets mapper to the cell with given coordinates
func (t *TableMapper) Pos(this *Table, row, col int) *ParaMapper {
	ret := new(ParaMapper)
	ret.ConnectTo(t.fm)
	ret.rider.Pos(this.cellCache[col][row])
	return ret
}
