package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Isbn     string    `json:"isbn"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	wanted := &Movie{}
	id := mux.Vars(r)["id"]
	for i := range movies {
		if movies[i].ID == id {
			wanted = &movies[i]
			break
		}
	}
	if (*wanted).ID != "" {
		json.NewEncoder(w).Encode(*wanted)
	} else {
		fmt.Fprint(w, "{\"status\":\"fail\",\"errorMessage\":\"Not Found\"}")
	}
}
func addMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMovie Movie
	json.NewDecoder(r.Body).Decode(&newMovie)
	newMovie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, newMovie)
	//fmt.Fprint(w, "{\"status\":\"ok\"}")
	json.NewEncoder(w).Encode(newMovie)
}
func removeMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	for i, element := range movies {
		if element.ID == id {
			movies = append(movies[:i], movies[i+1:]...)
			fmt.Fprint(w, "{\"status\":\"ok\"}")
			return
		}
	}
	fmt.Fprint(w, "{\"status\":\"fail\",\"errorMessage\":\"Not Found\"}")
}
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	wanted := &Movie{}
	id := mux.Vars(r)["id"]
	for i := range movies {
		if movies[i].ID == id {
			wanted = &movies[i]
			break
		}
	}
	if (*wanted).ID != "" {
		// This will allow the user to update the data that are in the request body, and it is better than removing the movies and recreate it with the new data
		json.NewDecoder(r.Body).Decode(wanted)
		json.NewEncoder(w).Encode(*wanted)
	} else {
		fmt.Fprint(w, "{\"status\":\"fail\",\"errorMessage\":\"Not Found\"}")
	}
}

func main() {
	movies = append(movies, Movie{ID: "1", Title: "Monster hunter", Isbn: "123456", Director: &Director{Firstname: "Milla", Lastname: "Jovovich"}})
	router := mux.NewRouter()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/create", addMovies).Methods("POST")
	router.HandleFunc("/movies/{id}", removeMovies).Methods("DELETE")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	fmt.Println("Starting on 8000 port")
	http.ListenAndServe(":8000", router)
}
