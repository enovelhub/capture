package engine

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/enovelhub/capture/httpool"
)

const rcfile = `
rc.WithDomain("freenovelonline.com")
rc.Home.WithAuthor(func(args){
	doc = ToGoqueryDoc(args[0])

	selector = "#main-content > div > div.detail-top > p:nth-child(3) > a"
	author = doc.Find(selector).Text()
	Return(author)
})
rc.Home.WithName(func(args) {
	doc = ToGoqueryDoc(args[0])

	selector = "#main-content > div > div.detail-top > h2"
	name = doc.Find(selector).Text()
	Return(name)

})
rc.Home.WithIndexURL(func(args) {
	doc = ToGoqueryDoc(args[0])

	homeURL,_ = Get("homeURL")
	Return(homeURL)
})

rc.Index.WithChapterURL(func(args) {
	doc = ToGoqueryDoc(args[0])

	selector = "#ztitle > li > a"
	doc.Find(selector).Each(func(i,s) {
		href,_ = s.Attr("href")
		Return(href)
	})
})

rc.Chapter.WithTitle(func(args) {
	doc = ToGoqueryDoc(args[0])

	doc.Find(".title a").Remove()
	title = doc.Find(".title").Text()
	Return(title)
})

rc.Chapter.WithContent(func(args) {
	doc = ToGoqueryDoc(args[0])

	doc.Find("#game-width p").Each(func(i,s) {
		Return(s.Text())	
	})
})
`

func TestEngine(t *testing.T) {

	log := logrus.StandardLogger()
	log.Level = logrus.DebugLevel
	en := New(context.TODO(), log, "test-rcfile", []byte(rcfile))
	if err := en.Err(); err != nil {
		t.Error(fmt.Sprintf("%#v\n", err))
	}

	log.WithField("domain", fmt.Sprintf("%#v", en.rc.Domain)).Debug("show domain")

	doc, err := goquery.NewDocument("http://localhost:6060/pkg")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = en.rc.Home.Author(reflect.ValueOf(doc))
	if err != nil {
		t.Error(err)
		return
	}

	for _, str := range en.Returns() {
		log.Debug(str)
	}
}

func TestEngineExec(t *testing.T) {
	log := logrus.StandardLogger()
	log.Level = logrus.DebugLevel
	en := New(context.TODO(), log, "test-rcfile", []byte(rcfile))
	if err := en.Err(); err != nil {
		t.Error(fmt.Sprintf("%#v\n", err))
	}

	httpool := httpool.New(http.DefaultClient, 20, time.Second*8, 3)
	defer httpool.Close()
	homeURL := "http://freenovelonline.com/241360-i-am-legend.html"

	book, err := en.Exec(httpool, homeURL)
	if err != nil {
		t.Error(fmt.Sprintf("%#v\n", err))
	}

	log.Info(fmt.Sprintf("%+v", book))
}
