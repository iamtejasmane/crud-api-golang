package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `josn:"id"`
	Isbn     string    `josn:"isbn"`
	Title    string    `josn:"title"`
	Director *Director `josn:"director"`
}

type Director struct {
	Firstname string `josn:"firstname"`
	Lastname  string `josn:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMoive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, moive := range movies {
		if moive.ID == params["id"] {
			json.NewEncoder(w).Encode(moive)
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var moive Movie
	_ = json.NewDecoder(r.Body).Decode(&moive)

	moive.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, moive)

	json.NewEncoder(w).Encode(moive)
}

// 1. Delete the moive first and then append the new moive: NOT Recommended way
func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content-Type", "application/json")
	// params
	params := mux.Vars(r)
	// loop over the moves, range
	// delete the movie with i.d that you've sent
	// add a new movie - the movie that we send in the body of postman

	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)

			var newMovie Movie
			_ = json.NewDecoder(r.Body).Decode(&newMovie)

			newMovie.ID = params["id"]
			movies = append(movies, newMovie)

			json.NewEncoder(w).Encode(newMovie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie1", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438255", Title: "Movie2", Director: &Director{Firstname: "Johny", Lastname: "Seen"}})
	movies = append(movies, Movie{ID: "3", Isbn: "438238", Title: "Movie3", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMoive).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
