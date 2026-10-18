package ui

// Add appends children to parent, or to the node parent's Content returns.
func Add(parent Drawable, children ...Drawable) {
	parent = contentOf(parent)
	for _, c := range children {
		Attach(parent, c)
	}
}

// Insert places child at index i among parent's children.
func Insert(parent Drawable, i int, child Drawable) {
	parent = contentOf(parent)
	Attach(parent, child)
	pb := parent.box()
	copy(pb.children[i+1:], pb.children[i:])
	pb.children[i] = child
}

// Attach appends child to parent directly, ignoring ContentHost. Widgets use
// it to build their internal structure.
func Attach(parent, child Drawable) {
	if _, ok := parent.(Leaf); ok {
		panic("ui: node does not accept children")
	}
	cb := child.box()
	if cb.parent != nil {
		panic("ui: node already has a parent")
	}
	cb.parent = parent
	pb := parent.box()
	pb.children = append(pb.children, child)
	Relayout()
}

// Remove detaches d from its parent and tells the parent. A node may remove
// itself. Focus inside d moves to the nearest remaining node.
func Remove(d Drawable) {
	b := d.box()
	p := b.parent
	if p == nil {
		return
	}
	if app != nil {
		app.Focus.leaving(b)
		app.Input.forget(d)
	}
	pb := p.box()
	for i, c := range pb.children {
		if c.box() == b {
			copy(pb.children[i:], pb.children[i+1:])
			pb.children[len(pb.children)-1] = nil
			pb.children = pb.children[:len(pb.children)-1]
			break
		}
	}
	b.parent = nil
	b.flags &^= flagLaidOut
	if r, ok := p.(ChildRemover); ok {
		r.ChildRemoved(d)
	}
	Relayout()
}

// RemoveChildren detaches every child of parent, or of its Content.
func RemoveChildren(parent Drawable) {
	pb := contentOf(parent).box()
	for len(pb.children) > 0 {
		Remove(pb.children[len(pb.children)-1])
	}
}

func contentOf(d Drawable) Drawable {
	if h, ok := d.(ContentHost); ok {
		if c := h.Content(); c != nil {
			return c
		}
	}
	return d
}

// IsAncestor reports whether a is d or contains d.
func IsAncestor(a, d Drawable) bool {
	ab := a.box()
	for n := d; n != nil; n = n.box().parent {
		if n.box() == ab {
			return true
		}
	}
	return false
}

// Attached reports whether d is in the live tree.
func Attached(d Drawable) bool {
	if app == nil {
		return false
	}
	return IsAncestor(app.root, d)
}
