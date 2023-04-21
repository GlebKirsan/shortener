package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/GlebKirsan/shortener/cmd/config"

	"github.com/go-chi/chi/v5"
)

type URL string

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var db = map[string]URL{}
var reverseIndex = map[URL]string{}

func setUp(newDB map[string]URL, newReverseIndex map[URL]string) {
	db = newDB
	reverseIndex = newReverseIndex
}

func generateString(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read path body", http.StatusBadRequest)
		return
	}

	var newAddress string
	if shortenedAddress, ok := reverseIndex[URL(body)]; ok {
		newAddress = shortenedAddress
	} else {
		newAddress = generateString(8)
		db[newAddress] = URL(body)
		reverseIndex[URL(body)] = newAddress
	}

	data := []byte(config.ResponsePrefix + "/" + newAddress)

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
	id := chi.URLParam(r, "id")
	url, ok := db[id]
	if !ok {
		http.Error(w, "Non-existing url-shorthand", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(url))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", shortenURL)
	r.Get("/{id}", getURL)
	return r
}

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(*config.ServerAddress, URLRouter()))
}
