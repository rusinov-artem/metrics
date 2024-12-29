package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	Srv     *http.Server
	address string
	logger  *zap.Logger
}

func New(logger *zap.Logger, h http.Handler, address string) *Server {
	return &Server{
		Srv: &http.Server{
			Handler: h,
		},
		address: address,
		logger:  logger,
	}
}

func (t *Server) Run() {
	ln, err := net.Listen("tcp", t.address)
	if err != nil {
		log.Fatalln(err)
	}
	t.logger.Info(fmt.Sprintf("Server started on '%s'", t.address))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = t.Srv.Serve(ln)
		if err != nil {
			log.Println(err)
		}
	}()

	// Вот тут обработка сигналов
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	t.logger.Info((<-s).String())

	ctx, closeFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeFn()
	err = t.Srv.Shutdown(ctx)
	if err != nil {
		t.logger.Error(err.Error(), zap.Error(err))
	}
	wg.Wait()
	t.logger.Info("Server stopped")
}
