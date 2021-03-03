//Package xmlns contains any constants related to ODF project, and also string names and values of ODF document nodes and attributes
//Also it contains some validation routines and value format descriptions
package xmlns

import (
	"github.com/iv-menshenin/odf/model"
)

const (
	Mimetype = "mimetype"
	Manifest = "META-INF/manifest.xml"
	Content  = "content.xml"
	Styles   = "styles.xml"
	Meta     = "meta.xml"
)

const (
	NSoffice       = "xmlns:office"
	NSmeta         = "xmlns:meta"
	NSconfig       = "xmlns:config"
	NStext         = "xmlns:text"
	NStable        = "xmlns:table"
	NSdraw         = "xmlns:draw"
	NSpresentation = "xmlns:presentation"
	NSdr3d         = "xmlns:dr3d"
	NSchart        = "xmlns:chart"
	NSform         = "xmlns:form"
	NSscript       = "xmlns:script"
	NSstyle        = "xmlns:style"
	NSnumber       = "xmlns:number"
	NSanim         = "xmlns:anim"
	NSdc           = "xmlns:dc"
	NSxlink        = "xmlns:xlink"
	NSmath         = "xmlns:math"
	NSxforms       = "xmlns:xforms"
	NSfo           = "xmlns:fo"
	NSsvg          = "xmlns:svg"
	NSsmil         = "xmlns:smil"
	NSmanifest     = "xmlns:manifest"
)

type AttrType int

const (
	NONE AttrType = iota
	STRING
	INT
	MEASURE
	ENUM
	COLOR
	BOOL
)

type Mime string

const (
	MimeDefault          = "text/xml"
	MimeText        Mime = "application/vnd.oasis.opendocument.text"
	MimeSpreadsheet Mime = "application/vnd.oasis.opendocument.spreadsheet"
)

var Typed map[model.AttrName]AttrType
var Enums map[model.AttrName][]string

func init() {
	Typed = make(map[model.AttrName]AttrType)
	Enums = make(map[model.AttrName][]string)
}
