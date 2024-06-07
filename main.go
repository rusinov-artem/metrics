package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8989")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server started!")
	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	go run(l, done, wg)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	<-s
	fmt.Println("Sginal arrive")
	wg.Wait()
	fmt.Println("Server stopped")
}

func run(l net.Listener, done chan struct{}, wg *sync.WaitGroup) {
out:
	for {
		select {
		case <-done:
			break out
		default:
		}

		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				b := make([]byte, 1000)
				n, err := c.Read(b)
				if err != nil {
					fmt.Println(err)
					return
				}
				if n > 0 {
					fmt.Print(string(b))
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}
	wg.Wait()
}
