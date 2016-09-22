package rc

import "github.com/mattn/anko/vm"

type Chapter struct {
	Title   vm.Func
	Content vm.Func
}

func (c *Chapter) WithTitle(f vm.Func) *Chapter {
	c.Title = f
	return c
}

func (c *Chapter) WithContent(f vm.Func) *Chapter {
	c.Content = f
	return c
}
