package text

import (
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns"
)

const (
	P         model.LeafName = "text:p"
	S         model.LeafName = "text:s"
	LineBreak model.LeafName = "text:line-break"
	Tab       model.LeafName = "text:tab"
	Span      model.LeafName = "text:span"
)

const (
	C          model.AttrName = "text:c"
	StyleName  model.AttrName = "text:style-name"
	AnchorType model.AttrName = "text:anchor-type"
)

const (
	Paragraph = "parahraph"
)

func init() {
	xmlns.Typed[C] = xmlns.INT
	xmlns.Typed[AnchorType] = xmlns.ENUM
	xmlns.Enums[AnchorType] = []string{Paragraph}
}
