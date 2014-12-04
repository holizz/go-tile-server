package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/holizz/go-tile-server"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3000
	}

	http.Handle("/tiles/", tiles.NewTileHandler("/tiles"))

	http.Handle("/public/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/", indexHandler)

	fmt.Printf("Listening on http://0.0.0.0:%d/\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
