package main

import (
    "fmt"
    "go-postgres/routers"
    "log"
    "net/http"
)

func main() {
	r := routers.Router()

	fmt.Println("Starting server on port 8084")

	log.Fatal(http.ListenAndServe(":8084", r))
}
