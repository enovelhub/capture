package rc

type Home struct {
	name     Filter
	author   Filter
	indexURL Filter
}

func (h *Home) Name(f Filter) *Home {
	h.name = f
	return h
}

func (h *Home) Author(f Filter) *Home {
	h.author = f
	return h
}

func (h *Home) IndexURL(f Filter) *Home {
	h.indexURL = f
	return h
}
