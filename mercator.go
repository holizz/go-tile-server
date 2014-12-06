package tiles

import "math"

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_numbers_to_lon..2Flat.
func getLonLatFromTileName(x, y, zoom int64) (float64, float64) {
	n := math.Pow(2, float64(zoom))
	lon := (float64(x) / n * 360) - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - (2 * float64(y) / n))))
	lat := latRad * 180 / math.Pi

	return lon, lat
}

func getXY(lon, lat, zoom float64) (float64, float64) {
	scale := math.Pow(2, zoom)
	x := (lon + 180) * scale
	y := (180 / math.Pi * math.Log(math.Tan(math.Pi/4+lat*(math.Pi/180)/2))) * scale

	return x, y
}

func getRelativeXY(lonEast, latNorth, lon, lat, scale float64) (float64, float64) {
	baseX, baseY := getXY(lonEast, latNorth, scale)
	nodeX, nodeY := getXY(lon, lat, scale)

	return nodeX - baseX, baseY - nodeY
}
