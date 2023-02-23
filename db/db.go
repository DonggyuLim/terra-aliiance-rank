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
	db             *mongo.Client
	url            string
	dbName         string
	collection     string
	collectionName string
	sync.Mutex
}

var d DB

func Connect() {
	db := DB{}
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

	db.db = client
	d = db
	fmt.Println("============DB connect==================")
}

func Close() {
	err := d.db.Disconnect(context.TODO())
	utils.HandleErr("DB Disconnect", err)
	fmt.Println("=========Connection to MongoDB closed=============")
}
func GetCollection() *mongo.Collection {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.db.Database(d.dbName).Collection(d.collection)
}

func Insert(account account.Account) error {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
	insertResult, err := collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	fmt.Println("Insert Complete", insertResult.InsertedID)
	return nil
}

func InsertMany(data []account.Account) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	var a []interface{}
	for _, el := range data {
		a = append(a, el)
	}
	exp := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
	_, err := collection.InsertMany(ctx, a)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Insert End")
}

func FindOne(filter bson.D) (account.Account, error) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	a := account.Account{}
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
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
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
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
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
	findOptions := options.Find()

	filter := bson.D{}

	cur, err := collection.Find(ctx, filter, findOptions)
	var curs []account.Account
	if err != nil {
		return curs, err
	}
	err := cur.All(context.TODO(), &curs)
	return curs, err
}

func FindAndReplace(filter, update bson.D) {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)

	result := collection.FindOneAndReplace(ctx, filter, update)
	fmt.Println("DB update")
	fmt.Println(result.Err().Error())
}

func ReplaceOne(filter bson.D, account account.Account) {

	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)

	_, err := collection.ReplaceOne(ctx, filter, account)

	utils.PanicError(err)
}

func UpdateOne(filter, update bson.D) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	fmt.Println("Update")
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	// utils.PanicError(err)
	fmt.Println(err.Error())
	// fmt.Println("Update End!")
}

func UpdateOneMap(filter bson.D, update bson.M) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	utils.PanicError(err)
}

func FindChain(address string, chainCode int, c *account.Reward) error {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.db.Database(d.dbName).Collection(d.collectionName)

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

	filter := bson.D{{Key: "address", Value: utils.MakeAddress(address)}}
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
