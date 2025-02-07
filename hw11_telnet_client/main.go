package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalln("missed parameters host or port")
	}

	// Извлечение параметров.
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")

	// Реакция на внешнее прерывание утилиты
	externalClose := make(chan os.Signal, 1)
	signal.Notify(externalClose, os.Interrupt, syscall.SIGTERM)

	ctx, cancelFunc := context.WithCancel(context.Background())
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	log.Println("trying to connect to the server")
	if err := client.Connect(); err == nil {
		log.Println("Connected to server")
		go func() {
			for {
				log.Println("Send to server")
				err = client.Send()
				if err != nil {
					log.Printf("error sending to server: %s\n", err)
					cancelFunc()
				}
			}
		}()

		go func() {
			for {
				err = client.Receive()
				log.Println("Receive from server")
				if err != nil {
					log.Printf("error receiving from server: %s\n", err)
					cancelFunc()
				}
			}
		}()
	} else {
		log.Println("error connecting to server")
		cancelFunc()
		client.Close()
		close(externalClose)
		log.Fatalln(err)
	}

	select {
	case <-ctx.Done():
		log.Println("Exit by telnet done")
		cancelFunc()
		client.Close()
		close(externalClose)
	case <-externalClose:
		log.Println("Exit by external cancel")
		cancelFunc()
		client.Close()
	}
}
