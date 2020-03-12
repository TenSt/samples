package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const stateStoreName = "statestore"
const stateURL = "http://localhost:3500/v1.0/state/" + stateStoreName

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/neworder", newOrder).Methods("POST")
	router.HandleFunc("/order", getOrder).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Go-write app is listening on port 8080.")
	log.Fatal(srv.ListenAndServe())
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Got new get order request.")
	req, err := http.NewRequest("GET", stateURL+"order", nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bodyBytes)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Finished with get order request.")
}

type neworder struct {
	Data struct {
		OrderID int `json:"orderId"`
	} `json:"data"`
}

func newOrder(w http.ResponseWriter, r *http.Request) {
	rBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var rBody neworder
	err = json.Unmarshal(rBodyBytes, &rBody)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Got new order: " + strconv.Itoa(rBody.Data.OrderID))

	data := `{"key": "order","value": "` + strconv.Itoa(rBody.Data.OrderID) + `"}`
	req, err := http.NewRequest("POST", stateURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully persisted state.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	w.WriteHeader(http.StatusOK)
}
