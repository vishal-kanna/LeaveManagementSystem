// package server

package main

import (
	"context"
	"encoding/json"
	"fmt"
	db "lms/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var admincollection *mongo.Collection
var leavescollection *mongo.Collection
var ctx = context.TODO()

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func AddStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	student := &db.Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	handleError(err)

	res, err := collection.InsertOne(context.Background(), student)
	handleError(err)
	json.NewEncoder(w).Encode(res)
}
func LeaveRequset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentleave := &db.LeaveRequset{}
	err := json.NewDecoder(r.Body).Decode(studentleave)
	studentleave.Status = "Pending"
	handleError(err)
	query := bson.M{
		"id": studentleave.Id,
	}
	var stu db.Student
	err1 := collection.FindOne(context.Background(), query).Decode(&stu)
	handleError(err1)
	res, err := admincollection.InsertOne(context.Background(), studentleave)
	res, err = leavescollection.InsertOne(context.Background(), studentleave)
	handleError(err)
	json.NewEncoder(w).Encode(res)
}
func ApproveRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	app := &db.Approve{}
	json.NewDecoder(r.Body).Decode(app)
	query := bson.M{
		"id": app.Id,
	}
	// fmt.Println("the approve req is", app)
	var StudentLeaveReq db.LeaveRequset
	// fmt.Println("query is", query)
	err := leavescollection.FindOne(context.Background(), query).Decode(&StudentLeaveReq)
	// fmt.Println("the student leave req", StudentLeaveReq)
	// handleError(err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	StudentLeaveReq.Status = "accepted"

	filter := bson.M{
		"$set": bson.M{
			"status": StudentLeaveReq.Status,
		},
	}
	// fmt.Println("hello")
	var a db.LeaveRequset
	_ = leavescollection.FindOneAndUpdate(context.Background(), query, filter).Decode(a)
	admincollection.FindOneAndUpdate(context.Background(), query, filter)
	if err != nil {
		log.Fatalf("err is %v", err)
	}
	if err == mongo.ErrNoDocuments {
		fmt.Println("err is", err)
	}
	json.NewEncoder(w).Encode(&a)
}
func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	handleError(err)
	fmt.Println("MongoDB connected")
	err = client.Connect(context.TODO())
	handleError(err)
	collection = client.Database("lms").Collection("student")
	admincollection = client.Database("lms").Collection("admin")
	leavescollection = client.Database("lms").Collection("leaves")
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/addstudent", AddStudent).Methods("POST")
	myRouter.HandleFunc("/leaveReq", LeaveRequset).Methods("GET")
	myRouter.HandleFunc("/approve", ApproveRequest).Methods("GET")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", myRouter))
}
