package main

import (
	"loader/pkg"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nats-io/nats.go"
)

const uri = "mongodb://root:example@localhost:27017"

func callExit() {
	log.Fatal("[INF] loader exiting")
}

func main() {

	log.SetFlags(log.Ldate | log.Lmicroseconds)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer callExit()
		<-c
		log.Panicln("[INF] initiating shutdown...")
	}()

	mgo := pkg.NewMgo(uri)

	n, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	nc := pkg.New(n)
	nc.Sub(mgo)

	wg.Wait()

}
