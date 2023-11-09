package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "hello world"}`))
    }).Methods("GET")

    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", router)
}