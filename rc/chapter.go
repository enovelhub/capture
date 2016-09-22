package rc

type Chapter struct {
	title   Filter
	content Filter
}

func (c *Chapter) Title(f Filter) *Chapter {
	c.title = f
	return c
}

func (c *Chapter) Content(f Filter) *Chapter {
	c.content = f
	return c
}
