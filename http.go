package tiles

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type TileHandler struct {
	prefix string
}

// prefix should be of the form "/tiles" (without the trailing slash)
func NewTileHandler(prefix string) *TileHandler {
	return &TileHandler{
		prefix: prefix,
	}
}

func (th *TileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(th.prefix):]

	if !(strings.HasPrefix(path, "/") && strings.HasSuffix(path, ".png")) {
		w.Write([]byte("404"))
		return
	}

	xyz := strings.Split(path[1:len(path)-4], "/")
	if len(xyz) != 3 {
		w.Write([]byte("404"))
		return
	}

	xyz_ := []int64{}
	for _, value := range xyz {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			w.Write([]byte("404"))
			return
		}

		xyz_ = append(xyz_, int64(intVal))
	}

	zoom := xyz_[0]
	x := xyz_[1]
	y := xyz_[2]

	fmt.Println(getLonLat(x, y, zoom))

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))

	err := png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_numbers_to_lon..2Flat.
func getLonLat(x, y, zoom int64) (float64, float64) {
	n := math.Pow(2, float64(zoom))
	lon := (float64(x) / n * 360) - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - (2 * float64(y) / n))))
	lat := latRad * 180 / math.Pi

	return lon, lat
}
