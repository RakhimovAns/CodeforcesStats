package utils

import (
	"context"
	"log/slog"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/valyala/fasthttp"
)

func Request(ctx context.Context, link string) (int, []byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(link)
	req.Header.SetMethod("GET")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}

	if err := client.Do(req, resp); err != nil {
		slog.Error("failed to get", "error", err.Error(), "url", link)
		return 0, nil, slerr.WithSource(err)
	}

	// Возвращаем копию тела ответа
	body := resp.Body()
	bodyCopy := make([]byte, len(body))
	copy(bodyCopy, body)

	return resp.StatusCode(), bodyCopy, nil
}
