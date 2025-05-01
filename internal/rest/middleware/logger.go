package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"users-api/pkg/logx"
)

type content struct {
	Host         string      `json:"host"`
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	Query        string      `json:"query,omitempty"`
	CallerIP     string      `json:"caller_ip"`
	RequestBody  interface{} `json:"request_body,omitempty"`
	ResponseBody interface{} `json:"response_body,omitempty"`
	Status       int         `json:"status"`
}

func Logger() echo.MiddlewareFunc {
	var rootContent = "http_info"
	return middleware.BodyDump(func(c echo.Context, reqBody, respBody []byte) {
		if skippings(c) {
			return
		}
		var (
			req  interface{}
			resp interface{}

			ctx = c.Request().Context()
		)

		encode(ctx, reqBody, &req)
		encode(ctx, respBody, &resp)

		if span, _ := tracer.SpanFromContext(c.Request().Context()); span != nil {
			span.SetTag("http.full_path", c.Request().RequestURI)
			if len(reqBody) != 0 {
				span.SetTag("http.request_body", string(reqBody))
			}
			if len(respBody) != 0 {
				span.SetTag("http.response_body", string(respBody))
			}
		}

		rootMsg := fmt.Sprintf("[*] req-resp form %s %s", strings.ToLower(c.Request().Method), c.Path())
		content := content{
			Host:         c.Request().Host,
			Method:       c.Request().Method,
			Path:         c.Path(),
			Query:        c.Request().URL.Query().Encode(),
			CallerIP:     c.RealIP(),
			RequestBody:  req,
			ResponseBody: resp,
			Status:       c.Response().Status,
		}

		if c.Response().Status >= http.StatusBadRequest && c.Response().Status < http.StatusInternalServerError {
			logx.Warnw(ctx, rootMsg, rootContent, content)
			return
		}

		if c.Response().Status == http.StatusInternalServerError {
			logx.Errorw(ctx, rootMsg, rootContent, content)
			return
		}

		logx.Infow(ctx, rootMsg, rootContent, content)
	})
}

func skippings(e echo.Context) bool {
	if strings.Contains(e.Path(), "swagger") {
		return true
	}
	if strings.Contains(e.Path(), "healthcheck") {
		return true
	}
	return false
}

func encode(ctx context.Context, src []byte, target interface{}) {
	if len(src) == 0 {
		return
	}

	err := json.Unmarshal(src, target)
	if err != nil {
		logx.Warnf(ctx, "can't unmarshal: %s", err.Error())
	}
}
