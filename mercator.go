package tiles

import (
	"fmt"
	"math"
)

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
	x := ((lon + 180) / 360) * scale * tileSize
	y := (180 / math.Pi * math.Log(math.Tan(math.Pi/4+lat*(math.Pi/180)/2))) * scale

	return x, y
}

func getRelativeXY(lonWest, latNorth, lon, lat, scale float64) (float64, float64) {
	baseX, baseY := getXY(lonWest, latNorth, scale)
	nodeX, nodeY := getXY(lon, lat, scale)
	x := nodeX - baseX
	y := baseY - nodeY

	if x < 0 || x >= tileSize {
		fmt.Printf("Error in X: %f %f\n", x, y)
	}

	if y < 0 || y >= tileSize {
		fmt.Printf("Error in Y: %f %f\n", x, y)
	}

	// fmt.Printf("XY: %f %f\n", x, y)
	return x, y
}
