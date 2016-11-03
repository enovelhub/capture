// Httpool is the pool for http request
// It provide a easy way for request retry where error happend,
// and use package context allow caller cancel it.
// It use gorotine parallely process http request,
// the size limit the size of gorotines.
package httpool

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type request struct {
	*http.Request
	respFuture chan response
}

type response struct {
	resp *http.Response
	err  error
}

type Pool struct {
	config struct {
		sync.Mutex
		retry int // defalut retry times
		n     int // core goroutine size
	}
	client    *http.Client
	reqsQueue chan request
	closed    bool
}

func New(size, retry int, timeout time.Duration) *Pool {
	n := size
	pool := &Pool{
		config: struct {
			sync.Mutex
			retry int // defalut retry times
			n     int // core goroutine size

		}{
			retry: retry,
			n:     size,
		},
		client: &http.Client{
			Timeout: timeout,
		},
		reqsQueue: make(chan request, n),
		closed:    false,
	}

	for i := 0; i < n; i++ {
		go pool.core(n)
	}
	return pool
}

func (pool *Pool) Do(httpreq *http.Request) (*http.Response, error) {
	pool.config.Lock()
	retry := pool.config.retry
	pool.config.Unlock()

	if retry < 1 {
		retry = 1
	}

	var (
		httpresp *http.Response
		err      error
	)
	for i := 0; i < retry; i++ {
		req := request{
			Request:    httpreq,
			respFuture: make(chan response),
		}
		pool.reqsQueue <- req
		resp := <-req.respFuture
		if resp.err != nil {
			if resp.resp != nil && resp.resp.Body != nil {
				resp.resp.Body.Close()
			}

		} else {
			return resp.resp, resp.err
		}
		httpresp, err = resp.resp, resp.err
	}

	if httpresp != nil && httpresp.Body != nil {
		httpresp.Body.Close()
	}
	return nil, err
}

func (pool *Pool) Close() {
	close(pool.reqsQueue)
}

func (pool *Pool) core(index int) {
	for req := range pool.reqsQueue {
		resp, err := pool.client.Do(req.Request)
		req.respFuture <- response{
			resp: resp,
			err:  err,
		}
		close(req.respFuture)
	}
}

func debug(args ...interface{}) {
	args = append([]interface{}{"debug[httpool]"}, args)
	fmt.Fprintln(os.Stderr, args...)
}
