package database

import (
	"context"
	"sync"
	"time"

	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	log    *logrus.Logger
	db     *mongo.Database
	Client *mongo.Client
}

//NewDataStore single
func NewDataStore(log *logrus.Logger, config config.Config) *mongoDB {
	var err error
	var connectOnce sync.Once
	var client *mongo.Client

	connectOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		credential := options.Credential{
			Username: config.GetDBConfig().UserName,
			Password: config.GetDBConfig().Password,
		}
		clientOpts := options.Client().ApplyURI(config.GetDBConfig().URL).SetAuth(credential).SetMaxPoolSize(20)
		client, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			_ = errors.Wrap(err, "mongo connect: ")
		}

		//db = client.Database(config.DatabaseName)
		log.Infoln("connect to mongoDB:")
	})

	if client != nil {
		mongoDataStore := new(mongoDB)
		mongoDataStore.db = client.Database(config.GetDBConfig().DBName)
		mongoDataStore.log = log
		mongoDataStore.Client = client
		return mongoDataStore
	}

	log.Fatal("failed to connect to database", err)
	return nil
}

//FindOne find one document
func (m *mongoDB) FindOne(db, coll, key string, value interface{}) (bson.M, error) {
	var result bson.M
	var err error
	var collection *mongo.Collection
	collection, err = m.Client.Database(db).Collection(coll).Clone()
	if err != nil {
		return nil, errors.Wrap(err, "mongodb find one")
	}

	filter := bson.D{primitive.E{Key: key, Value: value}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, errors.Wrap(err, "mongodb find")
}

//InsertOne insert one doc
func (m *mongoDB) InsertOne(db, coll string, value interface{}) error {
	collection := m.Client.Database(db).Collection(coll)
	_, err := collection.InsertOne(context.TODO(), value)

	return errors.Wrap(err, "mongodb insert one")
}

//InsertOne insert one doc
func (m *mongoDB) UpdateOne(db, coll string, key string, value interface{}, doc interface{}) error {
	filter := bson.D{primitive.E{Key: key, Value: value}}
	update := bson.M{"$set": doc}

	collection := m.Client.Database(db).Collection(coll)
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

	return errors.Wrap(err, "mongodb update one")
}

//DeleteOne delete one doc
func (m *mongoDB) DeleteOne(db, coll, key string, value interface{}) error {
	collection := m.Client.Database(db).Collection(coll)
	filter := bson.D{primitive.E{Key: key, Value: value}}
	_, err := collection.DeleteOne(context.TODO(), filter, nil)

	return errors.Wrap(err, "mongodb delete one")
}

/*
//CollectionCount count the doc of the Collection
func (m *mongoDB) CollectionCount() (string, int64) {
	collection := m.Client.Database(db).Collection(coll)
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

//按选项查询集合 Skip 跳过 Limit 读取数量 sort 1 ，-1 . 1 为最初时间读取 ， -1 为最新时间读取
func (m *mongoDB) CollectionDocuments(Skip, Limit int64, sort int) *mongo.Cursor {
	collection := m.Client.Database(m.Database).Collection(m.Collection)
	SORT := bson.D{{"_id", sort}} //filter := bson.D{{key,value}}
	filter := bson.D{{}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	//findOptions.SetLimit(i)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}
*/

/*
//DeleteMany: delete many
func (m *mongoDB) DeleteMany(key string, value interface{}) int64 {
	collection := m.Client.Database(m.Database).Collection(m.Collection)
	filter := bson.D{{key, value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount
}

//GetCount get count of records from DB
func (m *mongoDB) GetCount(ctx context.Context, id, date string) (string, error) {
	timeLayout := "2006-01-02 15:04:05"
	// if we need the time location, use time.LoadLocation("Local"), ParseInLocation
	startOfDay, _ := time.Parse(timeLayout, date+" 00:00:00")
	endOfDay, _ := time.Parse(timeLayout, date+" 23:59:59")

	collection := m.Client.Database(m.Database).Collection(m.Collection)
	filter := bson.M{"medallion": id, "pickup_datetime": bson.M{"$gte": startOfDay, "$lte": endOfDay}}
	//get the count
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
*/
