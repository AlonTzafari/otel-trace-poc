package db

import (
	"context"
	"errors"
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

func (db *DB) SaveDiceRoll(ctx context.Context, diceRoll DiceRoll) (bson.ObjectID, error) {
	res, err := db.coll.InsertOne(ctx, diceRoll)
	if err != nil {
		return bson.NilObjectID, err
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, errors.New("invalid id from dice roll insertion")
	}
	return id, nil
}
