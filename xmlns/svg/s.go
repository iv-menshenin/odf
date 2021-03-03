package svg

import (
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns"
)

const (
	FontFamily model.AttrName = "svg:font-family"
	Width      model.AttrName = "svg:width"
	Height     model.AttrName = "svg:height"
)

func init() {
	xmlns.Typed[Height] = xmlns.MEASURE
	xmlns.Typed[Width] = xmlns.MEASURE
}
