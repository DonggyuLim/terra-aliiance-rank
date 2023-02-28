package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client     *mongo.Client
	url        string
	dbName     string
	collection string
	sync.Mutex
}

type Height struct {
	Atreides  int
	Harkonnen int
	Corrino   int
	Ordos     int
}

var db DB

func Connect() {

	db.url = utils.LoadENV("DBURL", "db.env")
	db.dbName = utils.LoadENV("DBNAME", "db.env")
	db.collection = utils.LoadENV("Collection", "db.env")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(db.url).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, clientOptions)
	// ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	db.client = client
	// fmt.Println(db)
	fmt.Println("============DB connect==================")
}

func Close() {
	err := db.client.Disconnect(context.TODO())
	utils.HandleErr("DB Disconnect", err)
	fmt.Println("=========Connection to MongoDB closed=============")
}
func GetCollection() *mongo.Collection {

	return db.client.Database(db.dbName).Collection(db.collection)
}

func Insert(account account.Account) error {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	_, err := collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	// fmt.Println("Insert Complete", insertResult.InsertedID)
	return nil
}

func InsertMany(data []account.Account) {

	var a []interface{}
	for _, el := range data {
		a = append(a, el)
	}
	exp := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	_, err := collection.InsertMany(ctx, a)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Insert End")
}

func FindOne(filter bson.D) (account.Account, error) {

	a := account.Account{}
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	err := collection.FindOne(ctx, filter).Decode(&a)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return a, mongo.ErrNoDocuments
		}
		log.Fatal(err)
	}
	return a, nil
}

func Find(key, value, desc string, limit int64) ([]account.Account, error) {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	findOptions := options.Find()

	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: desc, Value: -1}})
	var filter primitive.D
	if key == "" && value == "" {
		filter = bson.D{}
	} else {
		filter = bson.D{{Key: key, Value: value}}
	}

	cur, _ := collection.Find(ctx, filter, findOptions)
	var curs []account.Account
	err := cur.All(context.TODO(), &curs)
	return curs, err
}

func FindAll() ([]account.Account, error) {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	findOptions := options.Find()

	filter := bson.D{}

	cur, err := collection.Find(ctx, filter, findOptions)
	var curs []account.Account
	if err != nil {
		return curs, err
	}
	err = cur.All(context.TODO(), &curs)
	return curs, err
}

func FindAndReplace(filter, update bson.D) {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)

	result := collection.FindOneAndReplace(ctx, filter, update)
	fmt.Println("DB update")
	fmt.Println(result.Err().Error())
}

func ReplaceOne(filter bson.D, account account.Account) {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)

	_, err := collection.ReplaceOne(ctx, filter, account)

	utils.PanicError(err)
}

func UpdateOne(filter, update bson.D) {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	_, err := collection.UpdateOne(ctx, filter, update)

	utils.PanicError(err)
	// fmt.Println("Update End!")
}

func UpdateOneMap(filter bson.D, update bson.M) {

	exp := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	_, err := collection.UpdateOne(ctx, filter, update)
	utils.PanicError(err)
	// fmt.Println(err.Error())
}

func FindChain(address string, chainCode int, c *account.Reward) error {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)

	var projection bson.D
	switch chainCode {
	case 0:
		// reward = fmt.Sprintf("%s.reward", "atreides")
		projection = bson.D{{Key: "atreides", Value: 1}}
	case 1:
		// reward = fmt.Sprintf("%s.reward", "harkonnen")
		projection = bson.D{{Key: "harkonnen", Value: 1}}
	case 2:
		// reward = fmt.Sprintf("%s.reward", "corrino")
		projection = bson.D{{Key: "corrino", Value: 1}}
	case 3:
		// reward = fmt.Sprintf("%s.reward", "ordos")
		projection = bson.D{{Key: "ordos", Value: 1}}
	}

	filter := bson.D{{Key: "address", Value: utils.MakeKey(address)}}
	opts := options.FindOne().SetProjection(projection)
	err := collection.FindOne(ctx, filter, opts).Decode(c)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return mongo.ErrNoDocuments
		}
		log.Fatal(err)
	}
	return nil
}
