package httpool

import (
	"context"
	"net/http"
	"time"
)

type Httpool struct {
	client   *http.Client
	size     int
	timeout  time.Duration
	retry    int
	reqQueue chan *HttpoolRequest
}

type HttpoolRequest struct {
	Req *http.Request
	Res chan *HttpoolResponse
}

type HttpoolResponse struct {
	Res *http.Response
	Err error
}

func NewRequest(req *http.Request) *HttpoolRequest {
	return &HttpoolRequest{
		Req: req,
		Res: make(chan *HttpoolResponse, 1),
	}
}

func New(
	client *http.Client,
	size int,
	timeout time.Duration,
	retry int,
) *Httpool {
	if size < 1 {
		size = 1
	}

	httpool := &Httpool{
		client:   client,
		size:     size,
		timeout:  timeout,
		retry:    retry,
		reqQueue: make(chan *HttpoolRequest, size),
	}

	for i := 0; i < size; i++ {
		go func(httpool *Httpool) {
			for httpoolReq := range httpool.reqQueue {
				retry := httpool.retry

			Retry:
				timeoutCtx, _ := context.WithTimeout(
					context.Background(),
					httpool.timeout,
				)
				req := httpoolReq.Req.WithContext(timeoutCtx)
				res, err := client.Do(req)
				if err != nil && retry > 0 {
					retry--
					goto Retry
				}
				httpoolReq.Res <- &HttpoolResponse{
					Res: res,
					Err: err,
				}
			}

		}(httpool)
	}

	return httpool
}

func (p *Httpool) Do(req *http.Request) (*http.Response, error) {
	httpoolReq := NewRequest(req)
	p.reqQueue <- httpoolReq
	httpoolRes := <-httpoolReq.Res
	return httpoolRes.Res, httpoolRes.Err
}

func (p *Httpool) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return p.Do(req)
}
