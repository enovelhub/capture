package engine

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/enovelhub/capture/book"
	"github.com/enovelhub/capture/charset"
	"github.com/enovelhub/capture/httpool"
	"github.com/enovelhub/capture/rc"
	enovelhub_core "github.com/enovelhub/capture/rc/builtins"
	"github.com/mattn/anko/vm"
)

type Engine struct {
	sync.Mutex
	ctx     context.Context
	name    string
	src     []byte
	log     *logrus.Logger
	ankovm  *vm.Env
	rc      *rc.RC
	returns []string
	envs    map[string]string
	err     error
}

func New(ctx context.Context, log *logrus.Logger, name string, src []byte) *Engine {
	ankovm := vm.NewEnv()
	ret := &Engine{
		ctx:     ctx,
		name:    name,
		src:     src,
		log:     log,
		ankovm:  ankovm,
		rc:      rc.New(),
		returns: nil,
		envs:    make(map[string]string),
		err:     nil,
	}

	enovelhub_core.LoadAllBuiltins(ankovm)
	ankovm.Define("rc", ret.rc)
	ankovm.Define("Get", ret.Get)
	ankovm.Define("Set", ret.Set)
	ankovm.Define("Return", ret.Return)

	ankovm.Define("ToGoqueryDoc", func(v reflect.Value) *goquery.Document {
		if doc, ok := v.Interface().(*goquery.Document); ok {
			return doc
		}
		return nil
	})
	_, err := ankovm.Execute(string(ret.src))
	ret.err = err

	return ret
}

func (e *Engine) Err() error {
	e.Lock()
	e.Unlock()
	return e.err
}

func (e *Engine) Get(k string) (string, bool) {
	e.Lock()
	e.Unlock()
	v, ok := e.envs[k]
	return v, ok
}
func (e *Engine) Set(k, v string) {
	e.Lock()
	e.Unlock()
	e.envs[k] = v
}

func (e *Engine) Return(v string) {
	e.Lock()
	e.Unlock()
	e.returns = append(e.returns, v)
}

func (e *Engine) Returns() []string {
	e.Lock()
	e.Unlock()
	return e.returns
}

func (e *Engine) ResetReturns() {
	e.Lock()
	e.Unlock()
	e.returns = nil
}

var (
	ErrSiteMapNotBelongHomeURL = errors.New("sitemap is not belong HomeURL" +
		"(HomeURL Host is not equal sitemap.Domain)")
	ErrIndexURLNotExist    = errors.New("indexURL is not exist in home page")
	ErrVaildHomeURLDomain  = errors.New("valid homeURL domain is rc Domain")
	ErrNotFoundName        = errors.New("cannot found book name by rc")
	ErrNotFoundAuthor      = errors.New("cannot found book author by rc")
	ErrNotFoundIndexURL    = errors.New("connot found book index page url by rc")
	ErrNotFoundChapterURLs = errors.New("connot found book chapter urls by rc")
)

func (e *Engine) Exec(httpool *httpool.Httpool, homeURL string) (*book.Book, error) {
	log := e.log

	retbook := &book.Book{}
	var indexURL string
	var chapterURLs []string
	{
		log.WithFields(logrus.Fields{
			"domain":  e.rc.Domain,
			"homeURL": homeURL,
		}).Info("valid rc Domain is homeURL domain")

		url, err := url.Parse(homeURL)
		if err != nil {
			return nil, err
		}

		if url.Host != e.rc.Domain {
			return nil, ErrVaildHomeURLDomain
		}
		e.Set("homeURL", url.String())
	}

	{
		log.Info("Get book home page")
		homePage, err := NewGoqueryDocument(e.ctx, httpool, homeURL)
		if err != nil {
			return nil, err
		}

		{
			log.Info("get book name")
			if e.rc.Home.Name != nil {
				_, err := e.rc.Home.Name(reflect.ValueOf(homePage))
				if err != nil {
					return nil, err
				}

				retbook.Name = strings.Join(e.Returns(), " ")
				e.ResetReturns()
			}

			if retbook.Name == "" {
				return nil, ErrNotFoundName
			}
			e.Set("name", retbook.Name)
		}

		{
			log.Info("get book author")
			if e.rc.Home.Author != nil {
				_, err := e.rc.Home.Author(reflect.ValueOf(homePage))
				if err != nil {
					return nil, err
				}

				retbook.Author = strings.Join(e.Returns(), ",")
				e.ResetReturns()
			}

			if retbook.Author == "" {
				log.WithError(ErrNotFoundAuthor).Warn("get book author failure")
			}
			e.Set("author", retbook.Author)
		}

		log.WithFields(logrus.Fields{
			"name":   retbook.Name,
			"author": retbook.Author,
		}).Info("novel info")

		{
			log.Info("get book indexURL")
			if e.rc.Home.IndexURL != nil {
				_, err := e.rc.Home.IndexURL(reflect.ValueOf(homePage))
				if err != nil {
					return nil, err
				}

				returns := e.Returns()
				e.ResetReturns()

				if len(returns) == 0 || len(returns[0]) == 0 {
					return nil, ErrNotFoundIndexURL
				}
				rawIndexURL := returns[0]
				url, err := homePage.Url.Parse(rawIndexURL)
				if err != nil {
					return nil, err
				}

				indexURL = url.String()
				e.Set("indexURL", indexURL)
			}

			if len(indexURL) == 0 {
				return nil, ErrNotFoundIndexURL
			}
		}

	}

	{
		log.Info("get book chapterURLs")
		indexPage, err := NewGoqueryDocument(e.ctx, httpool, indexURL)
		if err != nil {
			return nil, err
		}

		if e.rc.Index.ChapterURL != nil {
			_, err := e.rc.Index.ChapterURL(reflect.ValueOf(indexPage))
			if err != nil {
				return nil, err
			}

			rawChapterURLs := e.Returns()
			e.ResetReturns()
			if len(rawChapterURLs) == 0 {

				return nil, ErrNotFoundChapterURLs
			}

			chapterURLs = make([]string, len(rawChapterURLs))
			for i, raw := range rawChapterURLs {
				url, err := indexPage.Url.Parse(raw)
				if err != nil {
					return nil, err
				}
				chapterURLs[i] = url.String()
			}
		}

		if len(chapterURLs) == 0 {
			return nil, ErrNotFoundChapterURLs
		}
	}

	{
		log.Info("get book chapters")
		chapterPage := make(chan *chapterPageItem)

		ctxGetPages, cancelGetPages := context.WithCancel(e.ctx)
		wgGetPages := &sync.WaitGroup{}
		wgGetPages.Add(len(chapterURLs))

		for i, chapterURL := range chapterURLs {
			go func(
				ctx context.Context,
				cp chan<- *chapterPageItem,
				i int, chapterURL string) {
				defer wgGetPages.Done()
				chapterPage, err := NewGoqueryDocument(ctx, httpool, chapterURL)
				cp <- &chapterPageItem{
					url:   chapterURL,
					page:  chapterPage,
					index: i,
					err:   err,
				}
			}(ctxGetPages, chapterPage, i, chapterURL)
		}

		// goroutine for close chapterPage chan
		go func(cp chan *chapterPageItem, wg *sync.WaitGroup) {
			wg.Wait()
			close(cp)
		}(chapterPage, wgGetPages)

		// wait goroutines closed
		defer func() {
			cancelGetPages()
			for _ = range chapterPage {
				// do nothing,only wait goroutines closed
			}
		}()

		// process chapter
		var (
			mutexProgress = &sync.Mutex{}
			got           = 0
			total         = len(chapterURLs)
		)
		mutexBC := &sync.Mutex{}
		retbook.Chapters = make([]book.Chapter, len(chapterURLs))
		for item := range chapterPage {
			if item.err != nil {
				return nil, item.err
			}

			var title string
			var content []string
			if e.rc.Chapter.Title != nil {
				_, err := e.rc.Chapter.Title(reflect.ValueOf(item.page))
				if err != nil {
					return nil, err
				}

				returns := e.Returns()
				e.ResetReturns()

				title = strings.Join(returns, " ")

			}
			if len(title) == 0 {
				title = fmt.Sprintf("chapter %d", item.index)
				log.WithFields(logrus.Fields{
					"index":    item.index,
					"url":      item.url,
					"indexURL": indexURL,
				}).Warn("not found chapter title")
			}

			if e.rc.Chapter.Content != nil {
				_, err := e.rc.Chapter.Content(reflect.ValueOf(item.page))
				if err != nil {
					return nil, err
				}

				returns := e.Returns()

				e.ResetReturns()

				for _, l := range returns {
					content = append(content, RevertRawContent(l)...)
				}

			}
			if len(content) == 0 {
				log.WithFields(logrus.Fields{
					"index":    item.index,
					"url":      item.url,
					"indexURL": indexURL,
				}).Warn("not found chapter content")
			}

			c := book.Chapter{
				Title:   title,
				Content: content,
			}

			mutexBC.Lock()
			retbook.Chapters[item.index] = c
			mutexBC.Unlock()

			mutexProgress.Lock()
			got++
			log.WithFields(logrus.Fields{
				"progress": fmt.Sprintf("%.2f%%",
					float64(got)/float64(total)*100),
				"got":   got,
				"total": total,
			}).Info("get chapters progress")
			mutexProgress.Unlock()
		}

	}

	return retbook, nil

}

type chapterPageItem struct {
	url   string
	page  *goquery.Document
	index int
	err   error
}

func NewGoqueryDocument(ctx context.Context, httpool *httpool.Httpool, url string) (doc *goquery.Document, err error) {
	resp, err := httpool.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := charset.ToUTF8(resp)
	if err != nil {
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	setDoucmentURL(doc, url)
	return doc, nil
}

func setDoucmentURL(doc *goquery.Document, URL string) {
	doc.Url, _ = url.Parse(URL)
}

func RevertRawContent(raw string) []string {
	return strings.Split(raw, "\n")
}
