package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/altsaqif/go-graphql-new/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var connectionString string = "mongodb://localhost:27017"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetProduct(id string) *model.Product {
	productCollec := db.client.Database("go_graphql").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	var Product model.Product
	err := productCollec.FindOne(ctx, filter).Decode(&Product)
	if err != nil {
		log.Fatal(err)
	}
	return &Product
}

func (db *DB) GetProducts() []*model.Product {
	productCollec := db.client.Database("go_graphql").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var Product []*model.Product
	cursor, err := productCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &Product); err != nil {
		panic(err)
	}

	return Product
}

func (db *DB) CreateProduct(productInfo model.CreateProduct) *model.Product {
	productCollec := db.client.Database("go_graphql").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := productCollec.InsertOne(ctx, bson.M{"nama": productInfo.Nama, "stok": productInfo.Stok, "harga": productInfo.Harga})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnProduct := model.Product{ID: insertedID, Nama: productInfo.Nama, Stok: productInfo.Stok, Harga: productInfo.Harga}
	return &returnProduct
}

func (db *DB) UpdateProduct(productId string, productInfo model.UpdateProduct) *model.Product {
	productCollec := db.client.Database("go_graphql").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateProductInfo := bson.M{}

	if &productInfo.Nama != nil {
		updateProductInfo["nama"] = productInfo.Nama
	}
	if &productInfo.Stok != nil {
		updateProductInfo["stok"] = productInfo.Stok
	}
	if &productInfo.Harga != nil {
		updateProductInfo["harga"] = productInfo.Harga
	}

	_id, _ := primitive.ObjectIDFromHex(productId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateProductInfo}

	results := productCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var Product model.Product

	if err := results.Decode(&Product); err != nil {
		log.Fatal(err)
	}

	return &Product
}

func (db *DB) DeleteProduct(productId string) *model.DeleteProduct {
	productCollec := db.client.Database("go_graphql").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(productId)
	filter := bson.M{"_id": _id}
	_, err := productCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteProduct{DeletedProductID: productId}
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Users
func (db *DB) CreateUser(userInfo model.NewUser) *model.User {
	userCollec := db.client.Database("go_graphql").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	hashedPassword, err := HashPassword(userInfo.Password)
	if err != nil {
		log.Fatal(err)
	}
	inserg, err := userCollec.InsertOne(ctx, bson.M{"name": userInfo.Name, "email": userInfo.Email, "password": hashedPassword})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	hashedPasswordNew, err := HashPassword(userInfo.Password)
	if err != nil {
		log.Fatal(err)
	}
	returnUser := model.User{ID: insertedID, Name: userInfo.Name, Email: userInfo.Email, Password: hashedPasswordNew}
	return &returnUser
}
