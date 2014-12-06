package tiles

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type TileHandler struct {
	prefix string
	font   *truetype.Font
}

// prefix should be of the form "/tiles" (without the trailing slash)
func NewTileHandler(prefix, fontPath string) *TileHandler {
	font_, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(font_)
	if err != nil {
		panic(err)
	}

	return &TileHandler{
		prefix: prefix,
		font:   font,
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

	lon, lat := getLonLat(x, y, zoom)

	img, err := drawTile(lon, lat, th.font)
	if err != nil {
		panic(err)
	}

	err = png.Encode(w, img)
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

func drawTile(lon, lat float64, font *truetype.Font) (image.Image, error) {
	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Top/left border
	for i := 0; i < 256; i++ {
		img.Set(0, i, image.Black)
		img.Set(i, 0, image.Black)
	}

	// Dots
	for i := 0; i < 256; i += 16 {
		for j := 0; j < 256; j += 16 {
			img.Set(i, j, image.Black)
		}
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(12)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.Black)
	ctx.SetHinting(freetype.FullHinting)

	// pt := freetype.Pt(10, 10+int(c.PointToFix32(12)>>8))
	pt := freetype.Pt(10, 20)
	_, err := ctx.DrawString(fmt.Sprintf("%f, %f", lon, lat), pt)
	if err != nil {
		panic(err)
	}

	return img, nil
}
