package tiles

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
)

const tileSize = 256

func DrawTile(nwPt, sePt Pointer, zoom int64, data *OsmData, debug bool) (image.Image, error) {
	t := time.Now()

	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Plot some features
	//for fName, features := range data.Features {
	for _, feature := range data.Findex.GetFeatures(nwPt, sePt, zoom, data) {
		if mz, ok := mapFeatures[feature.FName]; ok {
			if mz.MinZoom > zoom {
				continue
			}
		} else {
			continue
		}

		switch feature.Type {
		case ItemTypeNode:
			//TODO
		case ItemTypeWay:
			way := data.Ways[feature.Id]
			//TODO: this ignores ways crossing over the tile
			withinBounds := false
			for _, node := range way.GetNodes(data.Nodes) {
				if node.Lon() > nwPt.Lon() && node.Lon() < sePt.Lon() &&
					node.Lat() < nwPt.Lat() && node.Lat() > sePt.Lat() {
					withinBounds = true
					break
				}
			}

			if !withinBounds {
				continue
			}

			coords := [][]float64{}
			for _, node := range way.GetNodes(data.Nodes) {
				x, y := getRelativeXY(nwPt, node, float64(zoom))
				coords = append(coords, []float64{x, y})
			}

			drawPolyLine(img, color.Black, coords)

		case ItemTypeRelation:
			//TODO
		}
	}

	// Debugging
	if debug {
		red := color.RGBA{0xff, 0x00, 0x00, 0xff}

		// Top/left border
		for i := 0; i < tileSize; i++ {
			img.Set(0, i, red)
			img.Set(i, 0, red)
		}

		// Tile location
		err := drawText(img, red, tileSize/2, 20, fmt.Sprintf("%f, %f", nwPt.Lon(), nwPt.Lat()))
		if err != nil {
			return nil, err
		}

		// Tile location
		err = drawText(img, red, tileSize/2, tileSize-20, time.Since(t).String())
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

func drawText(img *image.RGBA, cc color.Color, x, y int, s string) error {

	path := draw2dimg.NewGraphicContext(img)
	path.SetStrokeColor(cc)
	path.SetFillColor(cc)
	path.SetDPI(72)
	path.StrokeStringAt(s, float64(x), float64(y))

	return nil
}

func drawPolyLine(img *image.RGBA, cc color.Color, coords [][]float64) {
	path := draw2dimg.NewGraphicContext(img)

	path.SetStrokeColor(cc)
	path.SetLineWidth(1)

	for i, coord := range coords {
		if i == 0 {
			path.MoveTo(coord[0], coord[1])
		} else {
			path.LineTo(coord[0], coord[1])
		}
	}

	path.Close()
	path.Stroke()
}
