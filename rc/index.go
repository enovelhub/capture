package rc

import "github.com/mattn/anko/vm"

type Index struct {
	ChapterURL vm.Func
}

func (i *Index) WithChapterURL(f vm.Func) *Index {
	i.ChapterURL = f
	return i
}
