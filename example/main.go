package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/holizz/go-tile-server"
	"github.com/llgcode/draw2d"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3000
	}

	draw2d.SetFontFolder(".")
	tileHandler := tiles.NewTileHandler("/tiles", "isle-of-man-latest.osm.pbf")

	http.Handle("/tiles/", tileHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Printf("Listening on http://0.0.0.0:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
