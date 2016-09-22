package rc

type RC struct {
	Domain  string
	Home    *Home
	Index   *Index
	Chapter *Chapter
}

func New() *RC {
	return &RC {
		Domain: "",
		Home: &Home {},
		Index: &Index {},
		Chapter: &Chapter {},
	}
}

func (rc *RC) WithDomain(d string) *RC {
	rc.Domain = d
	return rc
}
