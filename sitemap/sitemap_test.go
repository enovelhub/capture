package sitemap

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	testYaml = `
domain: novel-website-domain
home:
  name: novel-name-selector
  author: novel-author-selector
  index: index-page-link-selector
index:
  chapter: chapter-page-link-selector
chapter:
  title: chapter-title-selector
  content: chapter-content-selector
`
)

func TestNew(t *testing.T) {
	sitemapExcept := NewSiteMapExcept(t)

	sitemap := New()
	sitemap.Domain = "novel-website-domain"

	sitemap.Home.Name = "novel-name-selector"
	sitemap.Home.Author = "novel-author-selector"
	sitemap.Home.Index = "index-page-link-selector"

	sitemap.Index.Chapter = "chapter-page-link-selector"

	sitemap.Chapter.Title = "chapter-title-selector"
	sitemap.Chapter.Content = "chapter-content-selector"

	if !reflect.DeepEqual(sitemapExcept, sitemap) {
		t.Error("not equal execpt value")
		return
	}
}

func TestNewFromYaml(t *testing.T) {
	{
		sitemapExcept := NewSiteMapExcept(t)
		sitemap, err := NewFromYaml([]byte(testYaml))
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(sitemapExcept, sitemap) {
			t.Error("not equal except value")
			return
		}
	}

	{
		sitemap, err := NewFromYaml(nil)
		if err == nil {
			t.Error("except err != nil")
			return
		}
		if sitemap != nil {
			t.Error("except sitemapFailure != nil")
			return
		}
	}

	{
		sitemap, err := NewFromYaml([]byte("this is not yaml"))
		if err == nil {
			t.Error("except err != nil")
			return
		}
		if sitemap != nil {
			t.Error("except sitemapFailure != nil")
			return
		}
	}
}

func TestNewFromReader(t *testing.T) {
	{
		sitemapExcept := NewSiteMapExcept(t)
		sitemap, err := NewFromReader(strings.NewReader(testYaml))
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(sitemapExcept, sitemap) {
			t.Error("not equal except value")
			return
		}
	}

	{
		sitemap, err := NewFromReader(nil)
		if err == nil {
			t.Error("except err != nil")
			return
		}
		if sitemap != nil {
			t.Error("except sitemapFailure != nil")
			return
		}

	}
}

func TestNewFromFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", ".temp4test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(tmpdir) // clean up

	sitemapFile, err := ioutil.TempFile(tmpdir, "sitemap.yml")
	if err != nil {
		t.Error(err)
		return
	}
	defer sitemapFile.Close()

	if _, err := sitemapFile.Write([]byte(testYaml)); err != nil {
		t.Error(err)
		return
	}

	{
		sitemapExcept := NewSiteMapExcept(t)
		sitemap, err := NewFromFile(sitemapFile.Name())
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(sitemapExcept, sitemap) {
			t.Error("not equal except value")
			return
		}
	}
	{
		sitemap, err := NewFromFile("file-that-not-exist")
		if err == nil {
			t.Error("except err != nil")
			return
		}
		if sitemap != nil {
			t.Error("except sitemap == nil")
			return
		}

	}
}

func TestString(t *testing.T) {
	sitemapExcept := NewSiteMapExcept(t)
	sitemap, err := NewFromYaml([]byte(sitemapExcept.String()))
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(sitemapExcept, sitemap) {
		t.Error("not equal except value")
		return
	}
}
func NewSiteMapExcept(t *testing.T) *SiteMap {
	sitemapExcept := New()
	err := yaml.Unmarshal([]byte(testYaml), sitemapExcept)
	if err != nil {
		t.Error(err)
		return nil
	}

	return sitemapExcept
}
