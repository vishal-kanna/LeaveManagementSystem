// package server

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	db "lms/db"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var admincollection *mongo.Collection
var leavescollection *mongo.Collection
var admincrendentials *mongo.Collection
var admintokens *mongo.Collection
var studencrentials *mongo.Collection
var studenttokens *mongo.Collection
var Token string

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var ctx = context.TODO()

func handleError(err error) {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Println("the error is", err)
	}
}
func AddStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	student := &db.Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	handleError(err)

	_, err = collection.InsertOne(context.Background(), student)
	handleError(err)
	json.NewEncoder(w).Encode("Student Added Successfully")
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
	var a db.StudentToken
	err2 := studenttokens.FindOne(context.Background(), query).Decode(&a)
	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode("Student didnot login or sign up")
		}
		log.Println("the error is", err)
	} else {
		_, err := admincollection.InsertOne(context.Background(), studentleave)
		_, err = leavescollection.InsertOne(context.Background(), studentleave)
		handleError(err)
		json.NewEncoder(w).Encode("Successfully applied for the leave")
	}
}
func ApproveRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logincheckquery := bson.M{
		"token": Token,
	}
	login := &db.AdminToken{}
	err1 := admintokens.FindOne(context.Background(), logincheckquery).Decode(login)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode("Admin didnot login ")
		} else {
			json.NewEncoder(w).Encode("Admin didnot login ")
		}
		// fmt.Println("hello")
		// log.Println(err1)
	} else {
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

}
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	admin := &db.Admin{}
	json.NewDecoder(r.Body).Decode(admin)
	// log.Println("admin is", admin.Passwd)
	query := bson.M{
		"username": admin.Username,
		"passwd":   admin.Passwd,
	}
	fmt.Println("the query is", query)
	admins := &db.Admin{}

	err := admincrendentials.FindOne(context.Background(), query).Decode(admins)
	// log.Println(cursor)
	handleError(err)
	// fmt.Println(admins)
	Admintoken := &db.AdminToken{
		Name:  admin.Username,
		Token: RandStringBytes(8),
	}
	Token = Admintoken.Token
	_, err = admintokens.InsertOne(context.Background(), Admintoken)
	handleError(err)
	json.NewEncoder(w).Encode("Admin Successfully Logged In")

	// var results []db.Admin
	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("the length of the results array is", len(results))
	// for _, result := range results {
	// 	fmt.Println(result)
	// }

}
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := bson.M{
		"token": Token,
	}
	res, err := admintokens.DeleteOne(context.Background(), query)
	handleError(err)
	if res.DeletedCount == 0 {
		json.NewEncoder(w).Encode("Admin didnot login")
	} else {
		json.NewEncoder(w).Encode("Admin Logged out successfully")
	}
}
func StudentSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentdetails := &db.StudentCredentials{}
	json.NewDecoder(r.Body).Decode(studentdetails)
	hashstring := studentdetails.Passwd
	h := sha256.New()
	h.Write([]byte(hashstring))
	studentdetails.Passwd = hex.EncodeToString(h.Sum(nil))
	_, err := studencrentials.InsertOne(context.Background(), studentdetails)
	handleError(err)
	// json.NewDecoder(w).Decode(res)
	json.NewEncoder(w).Encode("Student Signed Successfully")
}

func StudentLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//student should login using his id and password which were used at the time of sign up
	student := &db.StudentLoginCredentials{}
	json.NewDecoder(r.Body).Decode(student)
	hashstring := student.Passwd
	h := sha256.New()
	h.Write([]byte(hashstring))
	student.Passwd = hex.EncodeToString(h.Sum(nil))
	query := bson.M{
		"id":     student.Id,
		"passwd": student.Passwd,
	}
	studentde := &db.StudentCredentials{}

	err := studencrentials.FindOne(context.Background(), query).Decode(studentde)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode("Student didn't login")
		}
		log.Fatal(err)
	} else {
		studenttoken := &db.StudentToken{}
		studenttoken.Id = studentde.Id
		studenttoken.Studentoken = RandStringBytes(8)
		_, err = studenttokens.InsertOne(context.Background(), studenttoken)
		handleError(err)
		json.NewEncoder(w).Encode("Student logged in ")
	}

}
func StudentLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	studentid := &db.StudentLogout{}

	json.NewDecoder(r.Body).Decode(studentid)
	query := bson.M{
		"id": studentid.Id,
	}
	fmt.Println("student id is ", query)
	res, err := studenttokens.DeleteOne(context.Background(), query)
	handleError(err)
	if res.DeletedCount == 0 {
		json.NewEncoder(w).Encode("Student didnot login")
	} else {
		json.NewEncoder(w).Encode("Student Logged out successfully")
	}
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
	admincrendentials = client.Database("lms").Collection("admincreden")
	admintokens = client.Database("lms").Collection("admintokens")
	studencrentials = client.Database("lms").Collection("studencredentials")
	studenttokens = client.Database("lms").Collection("studenttokens")
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/addstudent", AddStudent).Methods("POST")
	myRouter.HandleFunc("/leaveReq", LeaveRequset).Methods("GET")
	myRouter.HandleFunc("/approve", ApproveRequest).Methods("GET")
	myRouter.HandleFunc("/adminlogin", AdminLogin).Methods("GET")
	myRouter.HandleFunc("/adminlogout", AdminLogout).Methods("GET")
	myRouter.HandleFunc("/studentsignup", StudentSignUp).Methods("POST")
	myRouter.HandleFunc("/studentlogin", StudentLogin).Methods("GET")
	myRouter.HandleFunc("/studnetlogout", StudentLogout).Methods("GET")
	// myRouter.HandleFunc("/stude")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", myRouter))
}
