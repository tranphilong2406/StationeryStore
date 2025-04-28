package job

import (
	"StoreServer/config"
	"StoreServer/utils/logger"
	"StoreServer/utils/response"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"reflect"
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

func (d *DB) NewObject() interface{} {
	tmp := reflect.TypeOf(d.TemplateObj)
	return reflect.New(tmp).Interface()
}

func (d *DB) NewList(slot int) interface{} {
	tmp := reflect.TypeOf(d.TemplateObj)
	return reflect.MakeSlice(reflect.SliceOf(tmp), 0, slot).Interface()
}

func (d *DB) Count(filter interface{}) response.Response {
	count, err := d.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return response.Response{
			Message: err.Error(),
			Data:    nil,
			Code:    http.StatusInternalServerError,
			Total:   int(count),
		}
	}

	return response.Response{
		Message: "Count query successfully!",
		Code:    http.StatusOK,
		Data:    nil,
		Total:   int(count),
	}
}

// Query get all objet in DB
func (d *DB) Query(inc interface{}, offset, limit int) response.Response {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := d.db.Collection(d.ColName).Find(context.TODO(), inc, findOptions)
	if err != nil {
		return response.Response{
			Message: "DB Error: " + err.Error(),
			Data:    nil,
			Code:    http.StatusInternalServerError,
		}
	}

	defer cursor.Close(context.TODO())

	var result = d.NewList(limit)
	err = cursor.All(context.TODO(), &result)
	if err != nil || reflect.ValueOf(result).Len() == 0 {
		return response.Response{
			Message: "No data found!",
			Data:    nil,
			Code:    http.StatusNotFound,
		}
	}

	return response.Response{
		Message: "Query " + d.ColName + " successfully!",
		Data:    result,
		Code:    http.StatusOK,
	}
}

// QueryOne get a specific object in DB
func (d *DB) QueryOne(inc interface{}) response.Response {
	return response.Response{}
}

// Create one object to db
func (d *DB) Create(inc interface{}) response.Response {
	obj, err := d.convertToBson(inc)
	if err != nil {
		return response.Response{
			Message: "Convert Error: " + err.Error(),
			Data:    nil,
			Code:    500,
		}
	}

	if obj["created_time"] == nil {
		obj["created_time"] = time.Now()
		obj["updated_time"] = obj["created_time"]
	} else {
		obj["updated_time"] = time.Now()
	}

	_, err = d.db.Collection(d.ColName).InsertOne(context.TODO(), obj)
	if err != nil {
		return response.Response{
			Message: "DB Error: " + err.Error(),
			Data:    nil,
			Code:    500,
		}
	}

	return response.Response{
		Message: "Create " + d.ColName + " successfully!",
		Data:    obj,
		Code:    http.StatusOK,
	}
}

// CreateMany insert many object to DB
func (d *DB) CreateMany(incList ...interface{}) (string, error) {
	objs := []bson.M{}
	ints := []interface{}{}

	if len(incList) == 1 {
		incList = incList[0].([]interface{})
	}

	for _, inc := range incList {
		obj, err := d.convertToBson(inc)
		if err != nil {
			return "Convert Error: " + err.Error(), err
		}

		if obj["created_time"] == nil {
			obj["created_time"] = time.Now()
			obj["updated_time"] = obj["created_time"]
		} else {
			obj["updated_time"] = time.Now()
		}

		objs = append(objs, obj)
		ints = append(ints, obj)
	}

	_, err := d.db.Collection(d.ColName).InsertMany(context.TODO(), ints)
	if err != nil {
		return "DB Error: " + err.Error(), err
	}

	return "Create " + d.ColName + "(s) successfully!", nil
}

// Update all matched item in DB
func (d *DB) Update(inc interface{}) (string, error) {
	return d.ColName, nil
}

// UpdateOne update one matched item
func (d *DB) UpdateOne(incList ...interface{}) (string, error) {
	return d.ColName, nil
}

// Delete all matched item
func (d *DB) Delete(selector interface{}) (string, error) {
	return d.ColName, nil
}

func (d *DB) DeleteOne(selector interface{}) (string, error) {
	return d.ColName, nil
}
