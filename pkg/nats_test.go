package pkg

import (
	"testing"

	"github.com/nats-io/nats.go"
)

func TestNats_Sub(t *testing.T) {
	type fields struct {
		NC *nats.Conn
	}
	type args struct {
		mgo *Mgo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Nats{
				NC: tt.fields.NC,
			}
			n.Sub(tt.args.mgo)
		})
	}
}
