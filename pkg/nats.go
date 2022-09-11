package pkg

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
)

type Nats struct {
	NC *nats.Conn
}

func New(nc *nats.Conn) *Nats {
	return &Nats{NC: nc}
}

func (n *Nats) Sub(mgo *Mgo) {
	ec, err := nats.NewEncodedConn(n.NC, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	js, err := ec.Conn.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	_, err = js.ConsumerInfo("LAB", "con")
	if err != nil {
		js.AddConsumer("LAB", &nats.ConsumerConfig{
			Durable: "con",
		})
	}

	_, err = js.StreamInfo("LAB")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[INF] loader running...")

	_, err = js.Subscribe(">", func(msg *nats.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INF] {\"seq_id\":%v,\"seq_recv\":\"%v\",\"subj\":\"%v\"}\n", meta.Sequence.Stream, meta.Timestamp, msg.Subject)
		mgo.Save(bson.D{
			{Key: "seq_id", Value: meta.Sequence.Stream},
			{Key: "message", Value: msg.Data},
		})
	}, nats.Durable("con"))

	if err != nil {
		log.Fatal(err)
	}
}

func (n *Nats) ChanSub(mgo *Mgo) {
	js, err := n.NC.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	_, err = js.StreamInfo("LAB")
	if err != nil {
		log.Fatal(err)
	}

	msgs := make(chan *nats.Msg, 1000)
	_, err = js.ChanSubscribe(">", msgs, nats.StartSequence(5))
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			msg := <-msgs
			meta, err := msg.Metadata()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(meta.Sequence.Stream)

			err = msg.Ack()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	wg.Wait()
}

// mongo 等到有 1000 筆時再進行寫入
