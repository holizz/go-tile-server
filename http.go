package tiles

import (
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type TileHandler struct {
	prefix string
	font   *truetype.Font
	data   *OsmData
}

// prefix should be of the form "/tiles" (without the trailing slash)
func NewTileHandler(prefix, pbfPath, fontPath string) *TileHandler {
	// Read font
	log.Println("Parsing font file...")
	font_, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(font_)
	if err != nil {
		panic(err)
	}
	log.Println("Parsing font file... [DONE]")

	// Read PBF
	log.Println("Parsing PBF file...")
	osmData, err := parsePbf(pbfPath)
	if err != nil {
		panic(err)
	}
	log.Println("Parsing PBF file... [DONE]")

	return &TileHandler{
		prefix: prefix,
		font:   font,
		data:   osmData,
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

	lonWest, latNorth := getLonLatFromTileName(x, y, zoom)
	lonEast, latSouth := getLonLatFromTileName(x+1, y+1, zoom)

	img, err := drawTile(lonWest, latNorth, lonEast, latSouth, float64(zoom), th.font, th.data)
	if err != nil {
		panic(err)
	}

	err = png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}
