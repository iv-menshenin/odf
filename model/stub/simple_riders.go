package stub

import (
	"github.com/iv-menshenin/odf/model"
	"github.com/iv-menshenin/odf/xmlns"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"image/color"
	"reflect"
)

type sr struct {
	base *sm
	pos  model.Leaf
	eol  bool
	this model.Leaf
}

func (r *sr) Base() model.Model {
	return r.base
}

func (r *sr) InitFrom(old model.Reader) {
	panic(126)
}

func (r *sr) Pos(p ...model.Leaf) model.Leaf {
	if len(p) == 1 {
		r.pos = p[0]
		if n, ok := r.pos.(model.Node); ok {
			r.eol = n.NofChild() == 0
		} else {
			r.eol = true
		}
	}
	return r.pos
}

func (r *sr) Read() model.Leaf {
	assert.For(r.pos != nil && !r.eol, 20)
	n, ok := r.pos.(model.Node)
	assert.For(ok, 21)
	idx := 0
	if r.this != nil {
		idx = n.IndexOf(r.this)
		idx++
	}
	if idx < n.NofChild() {
		r.this = n.Child(idx)
	} else {
		r.eol = true
	}
	return r.this
}

func (r *sr) Eol() bool {
	return r.eol
}

type sw struct {
	base *sm
	pos  model.Leaf
}

func (w *sw) Base() model.Model {
	return w.base
}

func (w *sw) InitFrom(old model.Writer) {
	if old != nil {
		w.Pos(old.Pos())
	}
}

func (w *sw) Pos(p ...model.Leaf) model.Leaf {
	if len(p) == 1 {
		w.pos = p[0]
	}
	return w.pos
}

func thisNode(l model.Leaf) model.Node {
	if _n, ok := l.(model.Node); ok {
		switch n := _n.(type) {
		case *sn:
			return n
		case *root:
			return n.inner
		default:
			halt.As(100, reflect.TypeOf(n))
		}
	}
	return nil
}

func (w *sw) Write(l model.Leaf, after ...model.Leaf) {
	assert.For(l != nil, 20)
	assert.For(w.pos != nil, 21)
	var splitter model.Leaf
	split := false
	if len(after) == 1 {
		splitter = after[0]
		split = true
	}
	add := func(source []model.Leaf, x model.Leaf) (ret []model.Leaf) {
		var (
			front []model.Leaf
			tail  []model.Leaf
		)
		switch {
		case split && splitter == nil:
			tail = source
		case split && splitter != nil:
			found := false
			for _, i := range source {
				if !found {
					front = append(front, i)
					found = i == splitter
				} else {
					tail = append(tail, i)
				}
			}
		default:
			front = source
		}
		ret = append(ret, front...)
		ret = append(ret, x)
		ret = append(ret, tail...)
		return
	}
	if _n, ok := w.pos.(model.Node); ok {
		switch n := _n.(type) {
		case *sn:
			n.children = add(n.children, l)
			l.Parent(n)
		case *root:
			n.inner.children = add(n.inner.children, l)
			l.Parent(n.inner)
		default:
			halt.As(100, reflect.TypeOf(n))
		}
	}
}

func (w *sw) Delete(l model.Leaf) {
	del := func(l []model.Leaf, x model.Leaf) (ret []model.Leaf) {
		for _, i := range l {
			if i != x {
				ret = append(ret, i)
			}
		}
		return
	}
	assert.For(l != nil, 20)
	assert.For(l.Parent() == thisNode(w.pos), 21, l.Parent(), w.pos.(model.Node))
	switch n := thisNode(w.pos).(type) {
	case *sn:
		n.children = del(n.children, l)
	case *root:
		n.inner.children = del(n.inner.children, l)
	default:
		halt.As(100, reflect.TypeOf(n))
	}

}

func (w *sw) WritePos(l model.Leaf, after ...model.Leaf) model.Leaf {
	w.Write(l, after...)
	return w.Pos(l)
}

func validateAttr(n model.AttrName, val string) {
	values := xmlns.Enums[n]
	found := false
	for _, v := range values {
		if v == val {
			found = true
		}
	}
	assert.For(found, 60, n, val)
}

func castAttr(n model.AttrName, i interface{}) (ret model.Attribute) {
	if i == nil {
		return nil
	}
	typ := xmlns.Typed[n]
	switch typ {
	case xmlns.NONE, xmlns.STRING:
		ret = &StringAttr{Value: i.(string)}
	case xmlns.INT:
		ret = &IntAttr{Value: i.(int)}
	case xmlns.ENUM:
		validateAttr(n, i.(string))
		ret = &StringAttr{Value: i.(string)}
	case xmlns.MEASURE:
		ret = &MeasureAttr{Value: i.(float64)}
	case xmlns.COLOR:
		ret = &ColorAttr{Value: i.(color.Color)}
	case xmlns.BOOL:
		ret = &BoolAttr{Value: i.(bool)}
	default:
		halt.As(100, typ, reflect.TypeOf(i))
	}
	return ret
}

func (w *sw) Attr(n model.AttrName, val interface{}) model.Writer {
	w.pos.Attr(n, castAttr(n, val))
	return w
}
