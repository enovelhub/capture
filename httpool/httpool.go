// Httpool is the pool for http request
// It provide a easy way for request retry where error happend,
// and use package context allow caller cancel it.
// It use gorotine parallely process http request,
// the size limit the size of gorotines.
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
	Ctx context.Context
	Req *http.Request
	Res chan *HttpoolResponse
}

type HttpoolResponse struct {
	Res *http.Response
	Err error
}

func NewRequest(ctx context.Context, req *http.Request) *HttpoolRequest {
	return &HttpoolRequest{
		Ctx: ctx,
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
				ctx, _ := context.WithTimeout(
					httpoolReq.Ctx,
					httpool.timeout,
				)
				req := httpoolReq.Req.WithContext(ctx)
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

func (p *Httpool) Close() {
	close(p.reqQueue)
}

func (p *Httpool) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	httpoolReq := NewRequest(ctx, req)
	p.reqQueue <- httpoolReq
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case httpoolRes := <-httpoolReq.Res:
		return httpoolRes.Res, httpoolRes.Err
	}
}

func (p *Httpool) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return p.Do(ctx, req)
}
