package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/enovelhub/capture/book"
	"github.com/enovelhub/capture/charset"
	"github.com/enovelhub/capture/httpool"
	"github.com/enovelhub/capture/sitemap"
	"net"
	"net/url"
	"sync"
	"strings"
)

var (
	ErrSiteMapNotBelongHomeURL = errors.New("sitemap is not belong HomeURL" +
		"(HomeURL's Host is not equal sitemap.Domain)")
	ErrIndexURLNotExist = errors.New("indexURL is not exist in home page")
)

func Exec(
	ctx context.Context,
	httpool *httpool.Httpool,
	sitemap *sitemap.SiteMap,
	homeURL string,
	log *logrus.Logger,
) (*book.Book, error) {
	var baseURL *url.URL
	{
		log.WithFields(logrus.Fields{
			"sitemap.Domain": sitemap.Domain,
			"homeURL":        homeURL,
		}).Info("valid sitemap is belong to homeURL")

		url, err := url.Parse(homeURL)
		if err != nil {
			return nil, err
		}

		if url.Host != sitemap.Domain {
			domain, _, err := net.SplitHostPort(url.Host)
			if err != nil {
				return nil, err
			}

			if domain != sitemap.Domain {
				return nil, ErrSiteMapNotBelongHomeURL
			}
		}
		baseURL = url
	}

	retbook := &book.Book{}
	indexURL := ""

	{
		log.Info("Get home page")
		home, err := NewGoqueryDocument(ctx, httpool, homeURL)
		if err != nil {
			return nil, err
		}

		{
			log.WithFields(logrus.Fields{
				"selector": sitemap.Home.Name,
			},
			).Info("Get name")
			retbook.Name = home.Find(sitemap.Home.Name).Text()
		}

		{
			log.WithFields(logrus.Fields{
				"selector": sitemap.Home.Author,
			},
			).Info("Get author")
			retbook.Author = home.Find(sitemap.Home.Author).Text()
		}

		{
			log.WithFields(logrus.Fields{
				"selector": sitemap.Home.Index,
			},
			).Info("Get index page url")

			if sitemap.Home.Index == "" {
				indexURL = homeURL
			} else {
				url, exist := home.Find(sitemap.Home.Index).Attr("href")
				if !exist {
					return nil, ErrIndexURLNotExist
				}
				baseURL,err := baseURL.Parse(url)
				if err != nil {
					return nil,err
				}
				indexURL = baseURL.String()
			}
		}
	}
	log.WithFields(logrus.Fields{
		"Name":   retbook.Name,
		"Author": retbook.Author,
	},
	).Info("book info")

	var indexes []string

	{
		log.WithFields(logrus.Fields{
			"indexURL": indexURL,
		},
		).Info("Get index page")

		doc, err := NewGoqueryDocument(ctx, httpool, indexURL)
		if err != nil {
			return nil, err
		}

		{
			log.Info("Get indexes")

			doc.Find(sitemap.Index.Chapter).Each(func(i int, s *goquery.Selection) {
				index, exist := s.Attr("href")
				if exist {
					url,err := baseURL.Parse(index)
					if err == nil {
						indexes = append(indexes, url.String())
					}
				}
			})
		}

	}

	retbook.Chapters = make([]book.Chapter, len(indexes))
	log.WithFields(logrus.Fields{
		"book.Chapters.len": len(retbook.Chapters),
	}).Info("chapters length")

	{
		log.Info("Start get chapters")
		var (
			totals        = len(retbook.Chapters)
			progress      = 0
			progressMutex = &sync.Mutex{}
		)
		bookMutext := &sync.Mutex{}
		var (
			Err      error
			ErrMutex = &sync.Mutex{}
		)
		ctx, cancelFunc := context.WithCancel(ctx)
		wg := &sync.WaitGroup{}
		wg.Add(len(retbook.Chapters))
		for i := 0; i < len(retbook.Chapters); i++ {
			go func(
				ctx context.Context,
				cancelFunc context.CancelFunc,
				wg *sync.WaitGroup,
				i int, url string,
			) {
				defer wg.Done()
				doc, err := NewGoqueryDocument(ctx, httpool, url)
				if err != nil {
					ErrMutex.Lock()
					if Err == nil {
						Err = err
					}
					ErrMutex.Unlock()
					cancelFunc()
					return
				}

				title := doc.Find(sitemap.Chapter.Title).Text()
				rawContent := doc.Find(sitemap.Chapter.Content).Text()
				content := RevertRawContent(rawContent)

				chapter := book.Chapter{
					Title:   title,
					Content: content,
				}

				bookMutext.Lock()
				retbook.Chapters[i] = chapter
				bookMutext.Unlock()

				progressMutex.Lock()
				progress++
				nowProgress := progress
				progressMutex.Unlock()

				log.WithFields(logrus.Fields{
					"total":    totals,
					"got":      nowProgress,
					"progress": fmt.Sprintf("%.2f%%", float64(nowProgress)/float64(totals)*100),
				}).Info("Progress")

			}(ctx, cancelFunc, wg, i, indexes[i])
		}

		wg.Wait()
	}

	log.Info("Ok")
	return retbook, nil
}

func NewGoqueryDocument(ctx context.Context, httpool *httpool.Httpool, url string) (*goquery.Document, error) {
	resp, err := httpool.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := charset.ToUTF8(resp)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func RevertRawContent(raw string) []string {
	return strings.Split(raw,"\n")
}
