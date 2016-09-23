package rc

import "github.com/mattn/anko/vm"

type Home struct {
	Name     vm.Func
	Author   vm.Func
	IndexURL vm.Func
}

func (h *Home) WithName(f vm.Func) *Home {
	h.Name = f
	return h
}

func (h *Home) WithAuthor(f vm.Func) *Home {
	h.Author = f
	return h
}

func (h *Home) WithIndexURL(f vm.Func) *Home {
	h.IndexURL = f
	return h
}
