package api

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/4everland/screenshot/chrome"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ScreenshotReq struct {
	URL     string `form:"url" binding:"required"`
	Width   int64  `form:"width,default=1920"`
	Height  int64  `form:"height,default=1080"`
	Delay   int64  `form:"delay"`
	Full    bool   `form:"full"`
	Timeout int64  `form:"timeout,default=15"`
	Output  string `form:"output,default=raw"`
}

func Screenshot(ctx *gin.Context) {
	var req ScreenshotReq
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if strings.Index(req.URL, "http") != 0 {
		req.URL = "https://" + req.URL
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	b, err := chrome.Screenshot(chrome.ScreenshotOptions{
		URL:     u,
		Width:   req.Width,
		Height:  req.Height,
		Delay:   time.Duration(req.Delay) * time.Millisecond,
		EndTime: time.Now().Add(time.Duration(req.Timeout) * time.Second),
		Full:    req.Full,
	})

	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	switch req.Output {
	case "base64":
		ctx.String(http.StatusOK, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(b))
	case "html":
		html := fmt.Sprintf(`<html><body><img src="data:image/png;base64,%s"/></body></html>`,
			base64.StdEncoding.EncodeToString(b))
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusOK, html)
	default:
		buf := bytes.NewBuffer(b)
		ctx.DataFromReader(
			http.StatusOK, int64(buf.Len()), "image/png", buf,
			map[string]string{
				"Content-Disposition": fmt.Sprintf(
					`attachment; filename="%x-%d.png"`,
					md5.Sum([]byte(u.String())), time.Now().Unix()),
			})
	}
}
