package chrome

import (
	"github.com/chromedp/chromedp"
	"net/url"
	"time"
)

//LocalChromeOptions Referenced from: https://github.com/puppeteer/puppeteer/issues/3938#issuecomment-475986157
var LocalChromeOptions = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.DisableGPU,
	chromedp.NoSandbox,
	chromedp.Flag("use-gl", "swiftshader"),
	chromedp.Flag("no-zygote", true),
	chromedp.Flag("disable-setuid-sandbox", true),
)

type ScreenshotOptions struct {
	URL     *url.URL
	Width   int64
	Height  int64
	Delay   time.Duration
	EndTime time.Time
	Full    bool
}
