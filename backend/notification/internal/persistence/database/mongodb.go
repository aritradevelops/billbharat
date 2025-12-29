package database

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NotInitializedErr(what string) error {
	return fmt.Errorf("%s is not initialized, have you forgot to call Connect() ?", what)
}

type Database interface {
	Connect() error
	Disconnect() error
	Health() error
	Collection(name string) *mongo.Collection
}

type MongoDB struct {
	connString string
	dbName     string
	client     *mongo.Client
}

func NewMongoDB(connString string) Database {

	return &MongoDB{
		connString: connString,
	}
}

// Connect implements Database.
func (m *MongoDB) Connect() error {

	uri, err := url.Parse(m.connString)
	if err != nil {
		return err
	}

	dbName := path.Base(uri.Path)
	if dbName == "/" {
		return fmt.Errorf("database name not found in connection string")
	}

	opts := options.Client().ApplyURI(m.connString)
	client, err := mongo.Connect(opts)
	if err != nil {
		return err
	}
	logger.Info().Msg("connected to mongodb")
	m.client = client
	return nil
}

// Disconnect implements Database.
func (m *MongoDB) Disconnect() error {
	if m.client == nil {
		return NotInitializedErr("MongoDB")
	}
	return m.client.Disconnect(context.Background())
}

// Health implements Database.
func (m *MongoDB) Health() error {
	if m.client == nil {
		return NotInitializedErr("MongoDB")
	}
	return m.client.Ping(context.Background(), nil)
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	if m.client == nil {
		return nil
	}
	return m.client.Database(m.dbName).Collection(name)
}
