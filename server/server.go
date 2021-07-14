package server

import (
	"context"
	"fmt"
	"github.com/4everland/screenshot/lib"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type (
	Server struct {
		engine *gin.Engine
		conf   Config
		server http.Server
	}

	Config struct {
		Host string
		Port int
		Mode string
	}
)

func NewServer(conf Config) Server {
	gin.SetMode(conf.Mode)

	r := gin.New()

	return Server{
		engine: Route(r),
		conf:   conf,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)

	s.engine.Use(ginzap.Ginzap(lib.Logger(), time.RFC3339, false))
	s.engine.Use(ginzap.RecoveryWithZap(lib.Logger(), true))

	s.server = http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	lib.Logger().Info("start http server "+addr, lib.HttpServerLog)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lib.Logger().Error("start http server failed, err:"+err.Error(), lib.HttpServerLog)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		lib.Logger().Error("shutdown http server failed, err:"+err.Error(), lib.HttpServerLog)
	}

	lib.Logger().Info("stop http server", lib.HttpServerLog)

	return nil
}
