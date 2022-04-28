// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.uber.org/zap"

	"github.com/Fenritec/remark42-cluster/store/engine"
)

//Backend the backend used for mongodb
type backend struct {
	client          *mongo.Client
	db              *mongo.Database
	options         *BackendOptions
	transactionPref *options.TransactionOptions

	dbName   string
	collName string
}

//Connect connects to the database
func Connect(uri string, opts *BackendOptions) (engine.IBackend, error) {
	var err error
	ret := backend{options: opts}
	ctx, cancel := context.WithTimeout(context.Background(), ret.options.QueryTimeout)
	defer cancel()

	dbName, uri, err := parseMongoURI(uri)
	if err != nil {
		return nil, err
	}

	preadpref, err := readpref.New(readpref.NearestMode)
	if err != nil {
		return nil, err
	}

	updatePref, err := readpref.New(readpref.PrimaryMode)
	if err != nil {
		return nil, err
	}

	ret.transactionPref = options.Transaction().SetReadPreference(updatePref)

	ret.client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(uri).SetReadPreference(preadpref))
	if err != nil {
		return nil, err
	}

	if err := ret.client.Ping(ctx, preadpref); err != nil {
		return nil, err
	}

	ret.db = ret.client.Database(dbName)

	return &ret, err
}

//Disconnect disconnects from the database
func (b *backend) Disconnect() error {
	if b.client == nil {
		return nil
	}

	ctx, cancel := b.createContext()
	defer cancel()
	err := b.client.Disconnect(ctx)
	return err
}

func parseMongoURI(uri string) (db, cleanURI string, err error) {
	db = "test"
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}
	if val := u.Query().Get("engine_db"); val != "" {
		db = val
	}

	q := u.Query()
	q.Del("engine_db")
	u.RawQuery = q.Encode()
	return db, u.String(), nil
}

var defaultOptions = BackendOptions{
	QueryTimeout: 10 * time.Second,
}

//BackendOptions is the option structure for MongoBackend
type BackendOptions struct {
	//QueryTimeout is the maximum time a query should take
	QueryTimeout time.Duration
	logger       *zap.Logger
}

//DefaultOptions returns the default options for MongoBackend
func DefaultOptions(logger *zap.Logger) BackendOptions {
	ret := defaultOptions
	ret.logger = logger
	return ret
}

//CreateContext creates a context to use with DB queries
func (b *backend) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), b.options.QueryTimeout)
}

type indexInfo struct {
	key    string
	unique bool
}
