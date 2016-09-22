package rc

import (
	"github.com/PuerkitoBio/goquery"
)

type Filter func(doc *goquery.Document) error

type RC struct {
	domain  string
	Home    Home
	Index   Index
	Chapter Chapter
}

func (rc *RC) Domain(d string) *RC {
	rc.domain = d
	return rc
}
