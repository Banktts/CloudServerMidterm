package main

import (
	"encoding/json"
	Service "example.com/CloudServerMidterm/service"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage hi!")
}

func handleRequests() {
	Router := mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/", homePage)
	Router.HandleFunc("/message", postMessage).Methods("POST")
	Router.HandleFunc("/message", updateMessage).Methods("PUT")
	Router.HandleFunc("/message/{uuid}", deleteMessage).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":80", Router))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var data Service.Message
	fmt.Print(data)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Service.CreateNewMessage(data)
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	var data Service.Message
	fmt.Print(data)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Service.UpdateMessage(data)
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params["uuid"])
	Service.DeleteMessage(params["uuid"])
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	handleRequests()
}
