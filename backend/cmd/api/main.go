package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    }).Methods("GET")

    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", router)
}