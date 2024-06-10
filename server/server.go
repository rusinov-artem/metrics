package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Srv     *http.Server
	address string
}

func New(h http.Handler, address string) *Server {
	return &Server{
		Srv: &http.Server{
			Handler: h,
		},
		address: address,
	}
}

func (t *Server) Run() {
	ln, err := net.Listen("tcp", t.address)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Server started on '%s'", t.address)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = t.Srv.Serve(ln)
		if err != nil {
			log.Println(err)
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	<-s

	ctx, closeFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeFn()
	err = t.Srv.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	wg.Wait()
	log.Println("Server stopped")
}
