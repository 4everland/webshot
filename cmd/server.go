package main

import (
	"flag"
	"github.com/4everland/screenshot/chrome"
	"github.com/4everland/screenshot/server"
)

var (
	host  = flag.String("host", "0.0.0.0", "http server host")
	port  = flag.Int("port", 30080, "http server port")
	mode  = flag.String("mode", "release", "gin mode")
	max   = flag.Int("max", 10, "chrome max thread num")
	path  = flag.String("path", "", "chrome exec path")
	proxy = flag.String("proxy", "", "chrome proxy")
)

func main() {
	flag.Parse()

	scheduler := chrome.NewScheduler(*max, chrome.NewLocalChrome(*path, *proxy))
	defer func() {
		scheduler.Chrome.Cancel()
		close(scheduler.Threads)
	}()

	server := server.NewServer(server.Config{
		Host: *host,
		Port: *port,
		Mode: *mode,
	})

	server.Run()
}
