package engine

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/enovelhub/capture/httpool"
	"github.com/enovelhub/capture/sitemap"
	"net/http"
	"testing"
	"time"
	"encoding/json"
)

const sitemapYaml = `
domain: "freenovelonline.com"
home:
    author: "#main-content > div > div.detail-top > p:nth-child(3) > a"
    name: "#main-content > div > div.detail-top > h2"
    index: ""
index:
    chapter: "#ztitle li a"
chapter:
    title: "#play-wrap > h3"
    content: "#game-width > div > p"
`

func TestExec(t *testing.T) {
	sitemap, err := sitemap.NewFromYaml([]byte(sitemapYaml))
	if err != nil {
		t.Error(err)
		return
	}

	logger := logrus.StandardLogger()
	logger.Level = logrus.DebugLevel
	httpool := httpool.New(http.DefaultClient, 20, time.Second*10, 3)
	homeURL := "http://freenovelonline.com/2434268-a-perfect-ten.html"

	book, err := Exec(context.Background(), httpool, sitemap, homeURL, logger)
	if err != nil {
		t.Error(err)
		return
	}

	data,_ := json.Marshal(book)
	logger.Debug(string(data))
}
