package chrome

import (
	"context"
	"time"

	"github.com/4everland/screenshot/lib"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type Chrome struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewLocalChrome(execPath, proxy string) *Chrome {
	if execPath != "" {
		LocalChromeOptions = append(LocalChromeOptions, chromedp.ExecPath(execPath))
	}

	if proxy != "" {
		LocalChromeOptions = append(LocalChromeOptions, chromedp.ProxyServer(proxy))
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), LocalChromeOptions...)
	return &Chrome{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func (c Chrome) Screenshot(parent context.Context, o ScreenshotOptions) (b []byte, err error) {
	timeoutCtx, cancel := context.WithTimeout(parent, time.Until(o.EndTime))
	defer cancel()

	ctx, cancel := chromedp.NewContext(timeoutCtx)
	defer cancel()

	if err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.EmulateViewport(o.Width, o.Height),
		chromedp.Navigate(o.URL.String()),
		chromedp.Sleep(o.Delay),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if o.Full {
				return chromedp.FullScreenshot(&b, 100).Do(ctx)
			}

			return chromedp.CaptureScreenshot(&b).Do(ctx)
		}),
	}); err != nil {
		lib.Logger().Error("chrome screenshot err:"+err.Error(), lib.ChromeLog)
	}

	return b, err
}

func (c Chrome) RawHtml(parent context.Context, o NewTabOptions) (b string, err error) {
	timeoutCtx, cancel := context.WithTimeout(parent, o.EndTime.Sub(time.Now()))
	defer cancel()

	ctx, cancel := chromedp.NewContext(timeoutCtx)
	defer cancel()

	if err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(o.URL.String()),
		chromedp.Sleep(o.Delay * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err := dom.GetOuterHTML().WithBackendNodeID(node.BackendNodeID).Do(ctx)
			if err == nil {
				b = html
			}
			return err
		}),
	}); err != nil {
		lib.Logger().Error("chrome catch html err:"+err.Error(), lib.ChromeLog)
	}

	return
}
