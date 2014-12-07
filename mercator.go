package tiles

import "math"

type Pointer interface {
	Lon() float64
	Lat() float64
}

type Point struct {
	Lon_, Lat_ float64
}

func (p Point) Lon() float64 { return p.Lon_ }
func (p Point) Lat() float64 { return p.Lat_ }

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_numbers_to_lon..2Flat.
func getLonLatFromTileName(x, y, zoom int64) Point {
	n := math.Pow(2, float64(zoom))
	lon := (float64(x) / n * 360) - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - (2 * float64(y) / n))))
	lat := latRad * 180 / math.Pi

	return Point{lon, lat}
}

func getXY(pt Pointer, zoom float64) (float64, float64) {
	scale := math.Pow(2, zoom)
	x := ((pt.Lon() + 180) / 360) * scale * tileSize
	y := (tileSize / 2) - (tileSize*math.Log(math.Tan((math.Pi/4)+((pt.Lat()*math.Pi/180)/2)))/(2*math.Pi))*scale

	return x, y
}

func getRelativeXY(nwPt, nodePt Pointer, scale float64) (float64, float64) {
	baseX, baseY := getXY(nwPt, scale)
	nodeX, nodeY := getXY(nodePt, scale)
	x := nodeX - baseX
	y := nodeY - baseY

	return x, y
}
