package rc

type Index struct {
	chapterURL Filter
}

func (i *Index) ChapterURL(f Filter) *Index {
	i.chapterURL = f
	return i
}
