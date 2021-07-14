package chrome

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

type Scheduler struct {
	Chrome  *Chrome
	Threads chan bool
}

type Task struct {
	ImageCh chan []byte
	Url     *url.URL
	EndTime time.Time
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

func (s *Scheduler) Exec(ch chan<- []byte, o ScreenshotOptions) {
	s.Threads <- true
	if o.EndTime.Before(time.Now()) {
		<-s.Threads
		return
	}

	b := s.Chrome.Screenshot(o)
	if o.EndTime.After(time.Now()) {
		ch <- b
	}

	close(ch)

	<-s.Threads
}

func Screenshot(o ScreenshotOptions) (b []byte, err error) {
	ch := make(chan []byte)
	go scheduler.Exec(ch, o)

	select {
	case b := <-ch:
		return b, nil
	case <-time.After(o.EndTime.Sub(time.Now())):
		return b, errors.New("time out")
	}

	return
}
