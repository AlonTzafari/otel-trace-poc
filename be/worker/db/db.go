package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

type DB struct {
	opts   *options.ClientOptions
	client *mongo.Client
	coll   *mongo.Collection
}

func New(uri string) *DB {
	monitor := otelmongo.NewMonitor()
	return &DB{
		opts: options.Client().
			ApplyURI(uri).
			SetMonitor(monitor).
			SetAuth(options.Credential{Username: "user", Password: "password"}),
	}
}

func (db *DB) Connect(ctx context.Context) error {
	client, err := mongo.Connect(db.opts)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = client.Ping(timeoutCtx, nil)
	if err != nil {
		return nil
	}

	db.coll = client.Database("app").Collection("dicerolls")

	return nil
}

func (db *DB) Disconnect(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

type DiceRoll struct {
	ID    bson.ObjectID `bson:"_id,omitempty"`
	Value int           `bson:"value"`
}

func (db *DB) SaveDiceRoll(ctx context.Context, diceRoll DiceRoll) error {
	if diceRoll.ID == bson.NilObjectID {
		diceRoll.ID = bson.NewObjectID()
	}
	fmt.Printf("db.coll %+v", db.coll)
	_, err := db.coll.InsertOne(ctx, diceRoll)
	return err
}
