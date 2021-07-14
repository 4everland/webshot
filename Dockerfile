FROM golang:1 AS builder
LABEL stage=screenshotbuilder

WORKDIR /build/screenshot
ADD . .
RUN go build -ldflags="-s -w" -o /app/screenshot cmd/server.go

FROM chromedp/headless-shell:latest
RUN apt-get update
RUN apt-get install -y dumb-init
RUN apt-get install -y ttf-wqy-microhei ttf-wqy-zenhei xfonts-wqy
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app

COPY --from=builder /app/screenshot /app/screenshot

ENTRYPOINT ["dumb-init" , "--", "./screenshot"]
