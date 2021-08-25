package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id          int64  `bson:"id"`
	Name        string `bson:"name"`
	Dob         string `bson:"dob"`
	Address     string `bson:"address"`
	Description string `bson:"description"`
	CreatedAt   string `bson:"createdAt"`
}

var users []User

func connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:test1234@cluster0.8icei.mongodb.net/userDB?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(),
		10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	return client, ctx, cancel, err
}

func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("userDB").Collection("User")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodes)
	json.NewEncoder(w).Encode(users)
}

func createNewUserDB(w http.ResponseWriter, r *http.Request) {

	client, ctx, cancel, err := connect()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	var filter, option interface{}
	filter = bson.D{
		{"id", bson.D{{"$eq", user.Id}}},
	}

	option = bson.D{{"_id", 0}}

	cursor, err := query(client, ctx, "userDB", "User", filter, option)
	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		panic(err)
	}

	b, err := json.Marshal(results)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		str := string(b)
		if str == "null" {

			collection := client.Database("userDB").Collection("User")

			result, _ := collection.InsertOne(ctx, user)

			json.NewEncoder(w).Encode(result)
		} else {
			fmt.Fprintf(w, "User already exists")
		}

	}

}

func query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := client.Database(dataBase).Collection(col)

	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect()
	fmt.Println(cancel)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	key := vars["id"]

	i, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var filter, option interface{}

	filter = bson.D{
		{"id", bson.D{{"$eq", i}}},
	}

	option = bson.D{{"_id", 0}}

	cursor, err := query(client, ctx, "userDB", "User", filter, option)
	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		panic(err)
	}

	b, err := json.Marshal(results)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		str := string(b)
		if str == "null" {
			fmt.Fprintf(w, "No user found")
		} else {
			fmt.Fprintf(w, str)
		}

	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	client, ctx, cancel, err := connect()
	fmt.Println(cancel)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	key := vars["id"]
	collection := client.Database("userDB").Collection("User")
	idPrimitive, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	res, err := collection.DeleteOne(ctx, bson.M{"id": idPrimitive})

	if err != nil {
		fmt.Fprintf(w, "Internal error")
	} else {

		if res.DeletedCount == 0 {
			fmt.Fprintf(w, "Id not found")
		} else {
			fmt.Fprintf(w, "Deleted")
		}
	}

}

func updateUserDataDB(w http.ResponseWriter, r *http.Request) {

	client, ctx, cancel, err := connect()
	fmt.Println(cancel)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	key := vars["id"]
	collection := client.Database("userDB").Collection("User")
	idPrimitive, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var filter, option interface{}

	filter = bson.D{
		{"id", bson.D{{"$eq", idPrimitive}}},
	}

	option = bson.D{{"_id", 0}}

	cursor, err := query(client, ctx, "userDB", "User", filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		panic(err)
	}

	b, err := json.Marshal(results)

	if err != nil {
		fmt.Fprintf(w, "No user foundd")
	} else {
		str := string(b)
		if str == "null" {
			fmt.Fprintf(w, "No user foundd")
		} else {

			fmt.Fprintf(w, str)

			res, err := collection.DeleteOne(ctx, bson.M{"id": idPrimitive})
			if err != nil {
				fmt.Fprintf(w, "No recod found with that id!")

			} else {
				fmt.Println("ress")
				fmt.Println(res)
				w.Header().Set("content-type", "application/json")
				var user User
				_ = json.NewDecoder(r.Body).Decode(&user)

				result, _ := collection.InsertOne(ctx, user)
				json.NewEncoder(w).Encode(result)

			}
		}
	}

}

func handleRequests() {

	//GorilaMux
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)

	myRouter.HandleFunc("/all", returnAll)

	myRouter.HandleFunc("/addUserDB", createNewUserDB).Methods("POST")

	myRouter.HandleFunc("/getUser/{id}", returnSingleUser).Methods("Get")

	myRouter.HandleFunc("/deleteUserFromId/{id}", deleteUser).Methods("DELETE")

	myRouter.HandleFunc("/updateUserDataDB/{id}", updateUserDataDB).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	users = []User{
		User{Id: 101, Name: "Joker", Dob: "28/09/2001", Address: "Bangalore", Description: "db1", CreatedAt: "30/01/2021"},
	}
	handleRequests()
}
