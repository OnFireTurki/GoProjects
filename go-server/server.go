package main

import (
	"fmt"
	"net/http"
)

func server() {
	// First we create the server and init the files path
	server := http.FileServer(http.Dir("./htmls"))
	// Handle The requests
	http.Handle("/", server)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" {
			fmt.Println("404 Response is being send now")
			http.Error(w, "404 NOT-FOUND", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			fmt.Println("405 Response is being send now")
			http.Error(w, "405 Method NOT-Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Add("Accepted", "YES")
		fmt.Fprintf(w, "Hello !!")
	})
	// FORM Request
	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/form" {
			fmt.Println("404 Response is being send now")
			http.Error(w, "404 NOT-FOUND", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			fmt.Println("405 Response is being send now")
			http.Error(w, "405 Method NOT-Allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}
		w.Header().Add("Accepted", "YES")
		name := r.FormValue("Name")
		city := r.FormValue("City")
		fmt.Fprintf(w, "Hello %v From %v", name, city)
	})
	fmt.Println("Starting from port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}

}
