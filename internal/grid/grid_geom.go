package grid

// CellOrigin returns the top-left absolute coordinate of a cell.
func (g *Grid) CellOrigin(c Cell) (x, y int) {
	return g.colPref[c.Col], g.rowPref[c.Row]
}

// AreaOrigin returns the top-left absolute coordinate of an area.
func (g *Grid) AreaOrigin(a Area) (x, y int) {
	return g.colPref[a.Col()], g.rowPref[a.Row()]
}

// AreaSize returns the absolute width and height of an area.
func (g *Grid) AreaSize(a Area) (w, h int) {
	w = g.colPref[a.ColEnd()] - g.colPref[a.Col()]
	h = g.rowPref[a.RowEnd()] - g.rowPref[a.Row()]
	return
}

// AreaRect returns x, y, w, h in absolute units for an area.
func (g *Grid) AreaRect(a Area) (x, y, w, h int) {
	x, y = g.AreaOrigin(a)
	w, h = g.AreaSize(a)
	return
}
