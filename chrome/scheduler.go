package chrome

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Scheduler struct {
	Chrome  *Chrome
	Threads chan bool
}

var (
	scheduler *Scheduler
	once      sync.Once
)

func NewScheduler(maxThread int, chrome *Chrome) *Scheduler {
	once.Do(func() {
		scheduler = &Scheduler{
			Chrome:  chrome,
			Threads: make(chan bool, maxThread),
		}
	})

	return scheduler
}

func Screenshot(o ScreenshotOptions) (b []byte, err error) {
	scheduler.Threads <- true
	defer func() {
		<-scheduler.Threads
	}()
	if o.EndTime.Before(time.Now()) {
		return nil, errors.New("time out")
	}

	b, err = scheduler.Chrome.Screenshot(scheduler.Chrome.Ctx, o)
	if errors.Is(err, context.Canceled) {
		return nil, errors.New("time out")
	}
	return b, nil
}

func RawHtml(o NewTabOptions) (b string, err error) {

	scheduler.Threads <- true
	defer func() {
		<-scheduler.Threads
	}()
	if o.EndTime.Before(time.Now()) {
		return "", errors.New("time out")
	}

	b, err = scheduler.Chrome.RawHtml(scheduler.Chrome.Ctx, o)
	if errors.Is(err, context.Canceled) {
		return "", errors.New("time out")
	}
	return
}
