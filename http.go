package tiles

import (
	"fmt"
	"image"
	"image/png"
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

	xyz_ := []uint64{}
	for _, value := range xyz {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			w.Write([]byte("404"))
			return
		}

		xyz_ = append(xyz_, uint64(intVal))
	}

	z := xyz_[0]
	x := xyz_[1]
	y := xyz_[2]

	fmt.Println(x, y, z)

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))

	err := png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}
