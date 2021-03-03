package odf

import (
	"github.com/iv-menshenin/odf/generators"
	"github.com/iv-menshenin/odf/mappers"
	"github.com/iv-menshenin/odf/mappers/attr"
	"github.com/iv-menshenin/odf/model"
	_ "github.com/iv-menshenin/odf/model/stub"
	"github.com/iv-menshenin/odf/xmlns"
	"github.com/iv-menshenin/odf/xmlns/fo"
	"github.com/iv-menshenin/odf/xmlns/table"
	"github.com/kpmy/ypk/assert"
	"image/color"
	"os"
	"testing"
)

func TestModel(t *testing.T) {
	m := model.ModelFactory()
	if m == nil {
		t.Error("model is nil")
	}
	w := m.NewWriter()
	if w == nil {
		t.Error("writer is nil")
	}
}

func TestMappers(t *testing.T) {
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()
}

func TestGenerators(t *testing.T) {
	output, _ := os.OpenFile("test-basics.odf", os.O_CREATE|os.O_WRONLY, 0666)
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()
	generators.GeneratePackage(m, nil, output, fm.MimeType)
	assert.For(output.Close() == nil, 20)
}

func TestStructure(t *testing.T) {
	output, _ := os.OpenFile("test-text.odf", os.O_CREATE|os.O_WRONLY, 0666)
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()
	fm.WriteString("Hello, World!   \t   \n   \r	фыва 	фыва		\n фыва")
	generators.GeneratePackage(m, nil, output, fm.MimeType)
	assert.For(output.Close() == nil, 20)
}

func TestStylesMechanism(t *testing.T) {
	output, _ := os.OpenFile("test-styles.odf", os.O_CREATE|os.O_WRONLY, 0666)
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()
	fm.RegisterFont("Arial", "Arial")
	fm.RegisterFont("Courier New", "Courier New")
	fm.SetDefaults(new(attr.TextAttributes).Size(18).FontFace("Courier New"))
	fm.SetDefaults(new(attr.TextAttributes).Size(16).FontFace("Courier New"))
	fm.WriteString("Hello, World!\n")
	fm.SetAttr(new(attr.TextAttributes).Size(32).FontFace("Arial"))
	fm.WriteString(`Hello, Go!`)
	fm.SetAttr(new(attr.TextAttributes).Size(36).FontFace("Courier New").Bold().Italic())
	fm.WriteString(`	Hello, Again!`)
	fm.SetAttr(new(attr.TextAttributes).Size(32).FontFace("Arial")) //test attribute cache
	fm.SetAttr(new(attr.TextAttributes).Size(32).FontFace("Arial").Color(color.RGBA{0x00, 0xff, 0xff, 0xff}))
	fm.WriteString("\nNo, not you again!")
	fm.SetAttr(new(attr.ParagraphAttributes).AlignRight().PageBreak())
	fm.WritePara("Page break!\r")
	fm.SetAttr(nil)
	fm.WriteString(`Hello, Пщ!`)
	generators.GeneratePackage(m, nil, output, fm.MimeType)
	assert.For(output.Close() == nil, 20)
}

func TestTables(t *testing.T) {
	table := func(fm *mappers.Formatter) {
		tm := &mappers.TableMapper{}
		tm.ConnectTo(fm)
		tm.Write("test", 5, 10)
		tt := tm.List["test"]
		tm.WriteColumns(tt, 4)
		tm.WriteRows(tt, 3)
		tm.Span(tt, 1, 2, 1, 3)
		tm.Pos(tt, 0, 0).WritePara("Hello, table world!")
		tm.Pos(tt, 1, 2).WritePara("Hello, table world!")
	}
	{
		output, _ := os.OpenFile("test-odt-tables.odf", os.O_CREATE|os.O_WRONLY, 0666)
		m := model.ModelFactory()
		fm := &mappers.Formatter{}
		fm.ConnectTo(m)
		fm.MimeType = xmlns.MimeText
		fm.Init()
		table(fm)
		generators.GeneratePackage(m, nil, output, fm.MimeType)
		assert.For(output.Close() == nil, 20)
	}
	{
		output, _ := os.OpenFile("test-ods-tables.odf", os.O_CREATE|os.O_WRONLY, 0666)
		m := model.ModelFactory()
		fm := &mappers.Formatter{}
		fm.ConnectTo(m)
		fm.MimeType = xmlns.MimeSpreadsheet
		fm.Init()
		table(fm)
		generators.GeneratePackage(m, nil, output, fm.MimeType)
		assert.For(output.Close() == nil, 20)
	}
}

func TestDraw(t *testing.T) {
	const ImagePng xmlns.Mime = "image/png"
	output, _ := os.OpenFile("test-draw.odf", os.O_CREATE|os.O_WRONLY, 0666)
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()
	embed := make(map[string]generators.Embeddable)
	{
		img, _ := os.Open("2go.png")
		d := mappers.NewDraw(img, ImagePng)
		url := d.WriteTo(fm, "Two Gophers", 6.07, 3.53) //magic? real size of `project.png`
		embed[url] = d
	}
	generators.GeneratePackage(m, embed, output, fm.MimeType)
	assert.For(output.Close() == nil, 20)
}

func TestTableStyles(t *testing.T) {
	output, _ := os.OpenFile("test-table-styles.odf", os.O_CREATE|os.O_WRONLY, 0666)
	m := model.ModelFactory()
	fm := &mappers.Formatter{}
	fm.ConnectTo(m)
	fm.MimeType = xmlns.MimeText
	fm.Init()

	fm.SetAttr(new(attr.TableAttributes).BorderModel(table.BorderModelCollapsing).AlignCenter().Width(10.0))
	fm.SetAttr(new(attr.TableRowAttributes).UseOptimalRowHeight()).SetAttr(new(attr.TableColumnAttributes).UseOptimalColumnWidth())
	fm.SetAttr(new(attr.TableCellAttributes).Border(attr.Border{Width: 0.01, Color: color.Black, Style: fo.Solid}))
	tm := &mappers.TableMapper{}
	tm.ConnectTo(fm)
	tm.Write("test", 5, 10)
	tt := tm.List["test"]
	tm.Pos(tt, 0, 0).WriteString("Hello!")

	generators.GeneratePackage(m, nil, output, fm.MimeType)
	assert.For(output.Close() == nil, 20)
}
