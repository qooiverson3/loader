package pkg

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMgo_Save(t *testing.T) {
	type fields struct {
		Client *mongo.Client
	}
	type args struct {
		event bson.D
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
			mgo := &Mgo{
				Client: tt.fields.Client,
			}
			mgo.Save(tt.args.event)
		})
	}
}
