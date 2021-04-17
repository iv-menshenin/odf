package attr

import (
	"fmt"
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns/fo"
	"github.com/iv-menshenin/odf/xmlns/style"
	"github.com/iv-menshenin/odf/xmlns/text"
)

//ParagraphAttributes is ODF Paragraph Family style fluent builder
type ParagraphAttributes struct {
	named
	easy
}

func (p *ParagraphAttributes) Equal(_a Attributes) (ok bool) {
	a, ok := _a.(*ParagraphAttributes)
	if ok {
		ok = p.equal(&a.easy)
	}
	return
}

func (p *ParagraphAttributes) Fit() model.LeafName { return text.P }

func (p *ParagraphAttributes) Write(wr model.Writer) {
	wr.Attr(style.Family, style.FamilyParagraph)
	wr.WritePos(New(style.ParagraphProperties))
	p.apply(wr)
}

//AlignLeft on page
func (p *ParagraphAttributes) AlignLeft() {
	p.put(fo.TextAlign, fo.Left, nil)
}

//AlignRight on page
func (p *ParagraphAttributes) AlignRight() {
	p.put(fo.TextAlign, fo.Right, nil)
}

//AlignCenter on page
func (p *ParagraphAttributes) AlignCenter() {
	p.put(fo.TextAlign, fo.Center, nil)
}

//AlignJustify on page
func (p *ParagraphAttributes) AlignJustify() {
	p.put(fo.TextAlign, fo.Justify, nil)
}

//AlignCustom allows you to set the text alignment attribute yourself.
func (p *ParagraphAttributes) AlignCustom(alignment string) {
	p.put(fo.TextAlign, alignment, nil)
}

//PageBreak with new paragraph written (it will be first on new page)
func (p *ParagraphAttributes) PageBreak() {
	p.put(fo.BreakBefore, true, func(v value) {
		if x := v.data.(bool); x {
			v.wr.Attr(fo.BreakBefore, fo.Page)
		}
	})
}

//SetIndent sets the indentation of the first line of a paragraph
func (p *ParagraphAttributes) SetIndent(inch float64) {
	p.put(fo.TextIndent, fmt.Sprintf("%0.4fin", inch), nil)
}

//SetMarginLeft adjusts the white space at the right border of the paragraph
func (p *ParagraphAttributes) SetMarginLeft(inch float64) {
	p.put(fo.MarginLeft, fmt.Sprintf("%0.4fin", inch), nil)
}

//SetMarginRight adjusts the white space at the right border of the paragraph
func (p *ParagraphAttributes) SetMarginRight(inch float64) {
	p.put(fo.MarginRight, fmt.Sprintf("%0.4fin", inch), nil)
}

//SetMarginLeft adjusts the white space at the right border of the paragraph
func (p *ParagraphAttributes) SetMarginTop(inch float64) {
	p.put(fo.MarginTop, fmt.Sprintf("%0.4fin", inch), nil)
}

//SetMarginRight adjusts the white space at the right border of the paragraph
func (p *ParagraphAttributes) SetMarginBottom(inch float64) {
	p.put(fo.MarginBottom, fmt.Sprintf("%0.4fin", inch), nil)
}
