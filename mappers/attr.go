package mappers

import (
	"github.com/iv-menshenin/odf/mappers/attr"
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns/office"
	"github.com/iv-menshenin/odf/xmlns/style"
	"github.com/iv-menshenin/odf/xmlns/svg"
	"github.com/iv-menshenin/odf/xmlns/text"
	"github.com/kpmy/ypk/halt"
	"strconv"
)

//Attr holds a number of style nodes of document model and writes attributes. Also it holds cache of recently used attributes for reuse
type Attr struct {
	doc     model.Model
	ds      model.Leaf //document styles
	ffd     model.Leaf //font-face decls
	as      model.Leaf //automatic styles
	ms      model.Leaf //master styles
	asc     model.Leaf //automatic styles
	ffdc    model.Leaf //font-face decls
	s       model.Leaf
	current map[string]attr.Attributes
	old     map[string]attr.Attributes
	fonts   map[string]model.Leaf
	stored  bool
	count   int
}

func (a *Attr) nextName() string {
	a.count++
	return "auto" + strconv.Itoa(a.count)
}

func (a *Attr) reset() {
	a.stored = true
	a.current = make(map[string]attr.Attributes)
}

//Init called when empty document is initialized
func (a *Attr) Init(m model.Model) {
	a.doc = m
	wr := a.doc.NewWriter()
	wr.Pos(a.doc.Root())
	a.ds = wr.WritePos(New(office.DocumentStyles))
	wr.Attr(office.Version, "1.0")
	a.ffd = wr.WritePos(New(office.FontFaceDecls))
	wr.Pos(a.ds)
	a.as = wr.WritePos(New(office.AutomaticStyles))
	wr.Pos(a.ds)
	a.ms = wr.WritePos(New(office.MasterStyles))
	wr.Pos(a.ds)
	a.asc = New(office.AutomaticStyles)
	a.ffdc = New(office.FontFaceDecls)
	a.old = make(map[string]attr.Attributes)
	a.fonts = make(map[string]model.Leaf)
	a.reset()
}

//Fit finds appropriate attributes for given name and calls closure for applying attributes
func (a *Attr) Fit(n model.LeafName, callback func(a attr.Attributes)) {
	fit := make(map[model.LeafName]attr.Attributes)
	for _, v := range a.current {
		fit[v.Fit()] = v
	}
	if a := fit[n]; a != nil {
		callback(a)
	}
}

//RegisterFont writes font entry to Font Face Declaration section
func (a *Attr) RegisterFont(name, fontface string) {
	if a.fonts[name] == nil {
		wr := a.doc.NewWriter()
		wr.Pos(a.ffd)
		wr.WritePos(New(style.FontFace))
		wr.Attr(style.Name, name)
		wr.Attr(svg.FontFamily, fontface)
		//TODO deep copy fontface node
		wr.Pos(a.ffdc)
		a.fonts[name] = wr.WritePos(New(style.FontFace))
		wr.Attr(style.Name, name)
		wr.Attr(svg.FontFamily, name)
	}
}

//Flush writes styles to document model
func (a *Attr) Flush() {
	if !a.stored {
		wr := a.doc.NewWriter()
		for _, v := range a.current {
			if n := v.Name(); n == "" && a.old[n] == nil {
				v.Name(a.nextName())
				wr.Pos(a.asc)
				wr.WritePos(New(style.Style))
				wr.Attr(style.Name, v.Name())
				v.Write(a.doc.NewWriter(wr))
				a.old[v.Name()] = v
			} else if n != "" && a.old[n] == nil {
				halt.As(100, v.Name())
			}
		}
		a.stored = true
	}
}

//OldAttr return already written attributes for new attribute if their contents equals
func (a *Attr) OldAttr(n attr.Attributes) attr.Attributes {
	for _, v := range a.old {
		if v.Equal(n) {
			return v
		}
	}
	return nil
}

//SetDefaults sets default attributes of document, only TextAttributes and ParagraphAttributes supported by now
func (a *Attr) SetDefaults(al ...attr.Attributes) {
	wr := a.doc.NewWriter()
	wr.Pos(a.ds)
	if a.s == nil {
		a.s = wr.WritePos(New(office.Styles))
		for _, x := range al {
			switch x.Fit() {
			case text.P, text.Span:
				wr.WritePos(New(style.DefaultStyle))
				x.Write(a.doc.NewWriter(wr))
				wr.Attr(style.Family, style.FamilyParagraph)
			default:
				halt.As(100, x.Fit())
			}
		}
	} else {
		wr.Delete(a.s)
		a.s = nil
		a.SetDefaults(al...)
	}
}
