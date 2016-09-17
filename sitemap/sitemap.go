package sitemap

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var (
	ErrYamlIsNil   = errors.New("yaml []byte is nil")
	ErrReaderIsNil = errors.New("r io.Reader is nil")
)

type Home struct {
	Name   string
	Author string
	Index  string
}
type Index struct {
	Chapter string
}
type Chapter struct {
	Title   string
	Content string
}

type SiteMap struct {
	Domain  string
	Home    Home
	Index   Index
	Chapter Chapter
}

func New() *SiteMap {
	return &SiteMap{
		Domain: "",
		Home: Home{
			Name:   "",
			Author: "",
			Index:  "",
		},
		Index: Index{
			Chapter: "",
		},
		Chapter: Chapter{
			Title:   "",
			Content: "",
		},
	}
}

func NewFromYaml(y []byte) (*SiteMap, error) {
	if len(y) == 0 {
		return nil, ErrYamlIsNil
	}

	sitemap := New()
	err := yaml.Unmarshal(y, sitemap)
	if err != nil {
		return nil, err
	}

	return sitemap, nil
}

func NewFromReader(r io.Reader) (*SiteMap, error) {
	if r == nil {
		return nil, ErrReaderIsNil
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewFromYaml(data)
}

func NewFromFile(name string) (*SiteMap, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewFromReader(file)
}

func (sitemap *SiteMap) String() string {
	data, _ := yaml.Marshal(&sitemap)
	return string(data)
}
