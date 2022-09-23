# Webshot
Webshot is a website screenshot service written in golang,
it uses Chrome Headless to generate website screenshots.

## Usage of HTTP Server
```
$ go run cmd/server.go --help
  
  -host string
        http server host (default "0.0.0.0")
  -max int
        chrome max thread num (default 10)
  -mode string
        gin mode (default "release")
  -path string
        chrome exec path
  -port int
        http server port (default 30080)
  -proxy string
        chrome proxy      
```
### Examples:
Download image from the command line.
```shell
$ curl "http://localhost:30080/screenshot?url=google.com" > google.png
```
View image in the browser.
```
http://localhost:30080/screenshot?url=google.com&output=html
```

## API
```GET /screenshot```

Name    | Type      | Default   | Description
----    | ----      | -------   | -----------
url     | string    |           | Target URL (**required**), http(s):// prefix is optional
timeout | int       | 15        | Maximum waiting time (seconds)
output  | string    | raw       | Output format (raw, base64, html)
width   | int       | 1920      | Viewport width
height  | int       | 1080      | Viewport height
full    | bool      | false     | Capture full page height
delay   | int       | 0         | Delay screenshot after page is loaded (milliseconds)

```GET /raw```

return raw html when window loaded

Name    | Type      | Default   | Description
----    | ----      | -------   | -----------
url     | string    |           | Target URL (**required**), http(s):// prefix is optional
timeout | int       | 15        | Maximum waiting time (seconds)
delay   | int       | 0         | Delay screenshot after page is loaded (milliseconds)

## Docker
### Running
```shell
# pull latest version
$ docker pull 4everland/webshot:latest

# run
$ docker run -d -p 30080:30080 --rm --name webshot 4everland/webshot
```

### Building 
```shell
$ docker build -t webshot .
```
## TODO
- [ ] Support Remote Chrome
- [ ] CLI Tools