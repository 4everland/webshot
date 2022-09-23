package api

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/4everland/screenshot/chrome"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RawHtmlReq struct {
	URL     string `form:"url" binding:"required"`
	Delay   int64  `form:"delay"`
	Timeout int64  `form:"timeout,default=15"`
}

func RawHtml(ctx *gin.Context) {
	var req RawHtmlReq
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

	b, err := chrome.RawHtml(chrome.NewTabOptions{
		URL:     u,
		Delay:   time.Duration(req.Delay) * time.Millisecond,
		EndTime: time.Now().Add(time.Duration(req.Timeout) * time.Second),
	})

	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	buf := bytes.NewBuffer([]byte(b))
	ctx.DataFromReader(
		http.StatusOK, int64(buf.Len()), "text/html", buf,
		map[string]string{
			"Content-Disposition": fmt.Sprintf(
				`attachment; filename="%x-%d.html"`,
				md5.Sum([]byte(u.String())), time.Now().Unix()),
		})

}
