package pkg

import (
	"log"
	"math/rand"

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
	js, err := n.NC.JetStream()
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

	u := rand.Int()
	_, err = js.Subscribe(">", func(msg *nats.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[con]%v [seq]%v [time]%v\n", u, meta.Sequence.Stream, meta.Timestamp)
		mgo.Save(bson.D{
			{Key: "seq_id", Value: meta.Sequence.Stream},
			{Key: "message", Value: string(msg.Data)},
		})
	}, nats.Durable("con"), nats.StartSequence(1))

	if err != nil {
		log.Fatal(err)
	}
}
