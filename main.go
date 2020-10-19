package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"io/ioutil"

	"github.com/gorilla/mux"
)

type book struct {
	ID     string  `json:"ID"`
	Name   string  `json:"Name"`
	Author string  `json:"Author"`
}

type allBooks []book

var books = allBooks {
	{
		ID:     "1",  
		Name:   "Beyaz Diş",
		Author: "Jack London",
	},
}

func library(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome library")
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter the Name and author")
	}

	json.Unmarshal(reqBody, &newBook)
	books = append(books, newBook)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newBook)
}

func getOneBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	
	for _, singleBook := range books {
		if singleBook.ID == bookID {
			json.NewEncoder(w).Encode(singleBook)
		}
	}
}

func getAllBook(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	var updatedBook book

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter the Name and Author")
	}
	json.Unmarshal(reqBody, &updatedBook)

	for i, singleBook := range books {
		if singleBook.ID == bookID {
			singleBook.Name = updatedBook.Name
			singleBook.Author = updatedBook.Author
			books = append(books[:i], singleBook)
			json.NewEncoder(w).Encode(singleBook)
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]

	for i, singleBook := range books {
		if singleBook.ID == bookID {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(w, "The book with ID %v has been deleted succesfully", bookID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", library)
	router.HandleFunc("/book", createBook).Methods("POST")
	router.HandleFunc("/books", getAllBook).Methods("GET")
	router.HandleFunc("/books/{id}", getOneBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PATCH")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}