package main

import (
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

type URL string

const (
	serverAddress = "http://localhost:8080/"
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var db = map[string]URL{}
var reverseIndex = map[URL]string{}

func generateString(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	var newAddress string
	if storedAddress, ok := reverseIndex[URL(body)]; ok {
		newAddress = storedAddress
	} else {
		newAddress = generateString(8)
		db[newAddress] = URL(body)
		reverseIndex[URL(body)] = newAddress
	}

	data := []byte(serverAddress + newAddress)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Error during response writing", http.StatusInternalServerError)
		return
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get is allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing shorthand path parameter", http.StatusBadRequest)
		return
	}

	url, ok := db[id]
	if !ok {
		http.Error(w, "Non-existing url-shorthand", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(url))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	serveMux := mux.NewRouter()
	serveMux.HandleFunc("/", shortenURL)
	serveMux.HandleFunc("/{id}", getURL)
	err := http.ListenAndServe(":8080", serveMux)
	if err != nil {
		panic(err)
	}
}
