package main

import (
	"encoding/json"
	Service "example.com/CloudServerMidterm/service"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage hi!")
}

func handleRequests() {
	Router := mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/", homePage)
	Router.HandleFunc("/api/messages", sync).Methods("GET")
	Router.HandleFunc("/api/messages", postMessage).Methods("POST")
	Router.HandleFunc("/api/messages/{uuid}", updateMessage).Methods("PUT")
	Router.HandleFunc("/api/messages/{uuid}", deleteMessage).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":80", Router))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var data Service.Message
	fmt.Print(data)
	err := json.NewDecoder(r.Body).Decode(&data)
	fmt.Print(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Service.CreateNewMessage(data, &w)
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	var data Service.Message
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&data)
	data.Uuid = params["uuid"]
	fmt.Println(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Service.UpdateMessage(data, &w)
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params["uuid"])
	Service.DeleteMessage(params["uuid"], &w)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	handleRequests()
}

func sync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	// convert query to int
	idx, err := strconv.Atoi(params["idx"][0])
	if err != nil {		
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qidx, err2 := strconv.Atoi(params["qidx"][0])
	if err2 != nil {		
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}
	// call sync
	messages := Service.SyncMessages(idx, qidx)
	json.NewEncoder(w).Encode(messages)
}
