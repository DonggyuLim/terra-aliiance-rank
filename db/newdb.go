package db

import (
	"context"
	"log"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ReturnDB() DB {
	DB := DB{
		url:        utils.LoadENV("DBURL", "db.env"),
		dbName:     utils.LoadENV("DBNAME", "db.env"),
		collection: utils.LoadENV("Collection", "db.env"),
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(DB.url).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, clientOptions)
	// ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	DB.client = client
	return DB
}

func (d *DB) Insert(account account.Account) error {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.client.Database(d.dbName).Collection(d.collection)
	_, err := collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	// fmt.Println("Insert Complete", insertResult.InsertedID)
	return nil
}

func (d *DB) FindOne(filter bson.D) (account.Account, error) {

	a := account.Account{}
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := d.client.Database(d.dbName).Collection(d.collection)
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
func (d *DB) UpdateOne(filter, update bson.D) {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.client.Database(db.dbName).Collection(db.collection)
	_, err := collection.UpdateOne(ctx, filter, update)

	utils.PanicError(err)
	// fmt.Println("Update End!")
}
