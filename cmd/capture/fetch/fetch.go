package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/enovelhub/capture/engine"
	"github.com/enovelhub/capture/httpool"
)

type cmdfetch struct {
	timeout time.Duration
	retry   int
	rc      string
	url     string
	size    int
}

func (c *cmdfetch) Name() string {
	return "fetch"
}

func (c *cmdfetch) Desc() string {
	return "fetch e-novel by its home url"
}

func (c *cmdfetch) Run(args []string) error {
	fset := flag.NewFlagSet(args[0], flag.ContinueOnError)

	fset.StringVar(&c.rc, "rc", "", "rc file path")
	fset.StringVar(&c.url, "url", "", "enovel home page url")
	fset.IntVar(&c.size, "n", 20, "goroutine size for http request")
	fset.IntVar(&c.retry, "r", 3, "retry times when http request failure")
	fset.DurationVar(&c.timeout, "t", time.Second*8, "timeout about http request")

	err := fset.Parse(args[1:])
	if err != nil {
		return err
	}

	if c.rc == "" || c.url == "" {
		return errors.New("rc and url is required")
	}

	log := logrus.StandardLogger()

	httpool := httpool.New(http.DefaultClient, c.size, c.timeout, c.retry)
	// read and check rc
	rcfile, err := ioutil.ReadFile(c.rc)
	if err != nil {
		log.Error("rc file open failure")
		return err
	}
	e := engine.New(context.TODO(), log, c.rc, rcfile)
	book, err := e.Exec(httpool, c.url)
	if err != nil {
		return err
	}

	err = json.NewEncoder(os.Stdout).Encode(book)
	if err != nil {
		return err
	}

	return nil
}

func New() *cmdfetch {
	return &cmdfetch{}
}
