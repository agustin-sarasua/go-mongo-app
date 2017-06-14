package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/agustin-sarasua/go-mongo-app/model"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var people []model.Person

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func CreatePersonEndpoint(s *mgo.Session) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session := s.Copy()
		defer session.Close()

		params := mux.Vars(req)
		var person model.Person
		err := json.NewDecoder(req.Body).Decode(&person)
		person.Id, _ = strconv.Atoi(params["id"])

		c := session.DB("gomongoapp").C("people")
		err = c.Insert(person)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Book with this ISBN already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert book: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Location", r.URL.Path+"/"+book.ISBN)
		w.WriteHeader(http.StatusCreated)

		// people = append(people, person)
		json.NewEncoder(w).Encode(people)
	}
}

func main() {
	// Mongo
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Timeout:  60 * time.Second,
		Database: "gomongoapp",
		Username: "risk",
		Password: "risk",
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	//ensureIndex(session)
	//Fin Mongo

	//Gorilla MUX
	router := mux.NewRouter()

	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint(session)).Methods("POST")

	var p = model.Person{Id: 1, Name: "Agustin", LastName: "Sarasua", Address: nil}
	people = append(people, p)
	fmt.Println("Hello there")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("gomongoapp").C("people")

	index := mgo.Index{
		Key:        []string{"Id"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
