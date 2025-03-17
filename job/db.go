package job

import (
	"StoreServer/config"
	"StoreServer/utils/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type DB struct {
	ColName     string
	DBName      string
	TemplateObj interface{}
	collection  *mongo.Collection
	db          *mongo.Database
}

var conn *mongo.Client

func DBConnect() {
	cfg := config.GetConfig()
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			log.Printf("Command: %v\n", e.Command)
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			log.Printf("Succeeded: %v\n", e.Reply)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			log.Fatalf("Succeeded: %v\n", e.Failure)
		},
	}
	fmt.Println(1)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MONGOURL).SetMonitor(monitor).SetServerAPIOptions(serverAPI)
	fmt.Println(2)
	conn, err = mongo.Connect(opts)
	if err != nil {
		logger.GetLogger().Fatal(err.Error())
	}

	fmt.Println(3)
	if err = conn.Ping(ctx, nil); err != nil {
		logger.GetLogger().Fatal(err.Error())
	}

	fmt.Println("Connected to MongoDB!")
}

func Disconnect() {
	if err := conn.Disconnect(context.Background()); err != nil {
		logger.GetLogger().Fatal(err.Error())
	}

	logger.GetLogger().Info("Disconnected from MongoDB!")
}

func GetDB(dbname string) *mongo.Database {
	return conn.Database(dbname)
}

func (d *DB) convertToBson(inc interface{}) (bson.M, error) {
	if inc == nil {
		return bson.M{}, nil
	}

	mar, err := bson.Marshal(inc)
	if err != nil {
		return nil, err
	}

	res := bson.M{}
	ok := bson.Unmarshal(mar, &res)
	if ok != nil {
		return nil, ok
	}

	return res, nil
}

func (d *DB) convertToObj(inc bson.M) (interface{}, error) {
	var obj interface{}

	if inc == nil {
		return obj, nil
	}

	bytes, err := bson.Marshal(inc)
	if err != nil {
		return nil, err
	}

	ok := bson.Unmarshal(bytes, obj)
	if ok != nil {
		return nil, ok
	}

	return obj, nil
}

func (d *DB) Init(dbname string) {
	d.DBName = dbname
	d.db = GetDB(d.DBName)
	d.collection = d.db.Collection(dbname)

}

func (d *DB) Create(inc interface{}) (string, error) {
	obj, err := d.convertToBson(inc)
	if err != nil {
		return "Convert Error: " + err.Error(), err
	}

	_, err = d.db.Collection(d.ColName).InsertOne(context.TODO(), obj)
	if err != nil {
		return "DB Error: " + err.Error(), err
	}

	return "Create " + d.ColName + " successfully!", nil
}
