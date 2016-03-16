package tiles

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/golang/geo/s2"
	"github.com/llgcode/draw2d/draw2dimg"
)

const tileSize = 256

func DrawTile(nwPt, sePt Pointer, zoom int64, data *OsmData, debug bool) (image.Image, error) {
	t := time.Now()

	// s2.Points of tile 4 vertex
	p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), nwPt.Lon()))
	p2 := s2.PointFromLatLng(s2.LatLngFromDegrees(sePt.Lat(), nwPt.Lon()))
	p3 := s2.PointFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), sePt.Lon()))
	p4 := s2.PointFromLatLng(s2.LatLngFromDegrees(sePt.Lat(), sePt.Lon()))
	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Plot some features
	for _, feature := range data.Findex.GetFeatures(nwPt, sePt, zoom) {
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
			var (
				// previous  plotting point, nil if already added to coords or first time
				prevCoords []float64
				// previous  s2 Point. Always present, except first time
				prevS2Point s2.Point
			)
			coords := [][]float64{}
			prevWithinBounds := false
			first := true
			for _, node := range way.GetNodes(data.Nodes) {
				if first {
					first = false
					x, y := getRelativeXY(nwPt, node, float64(zoom))
					if node.Lon() > nwPt.Lon() && node.Lon() < sePt.Lon() &&
						node.Lat() < nwPt.Lat() && node.Lat() > sePt.Lat() {
						coords = append(coords, []float64{x, y})
						prevWithinBounds = true
					} else {
						//x, y := getRelativeXY(nwPt, node, float64(zoom))
						//a1 := s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon()))
						prevS2Point = s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon()))
						prevCoords = []float64{x, y}
						prevWithinBounds = false
					}
					continue
				}

				if node.Lon() > nwPt.Lon() && node.Lon() < sePt.Lon() &&
					node.Lat() < nwPt.Lat() && node.Lat() > sePt.Lat() {
					if prevWithinBounds == false {
						if len(prevCoords) > 0 {
							coords = append(coords, prevCoords)
						}
					}
					x, y := getRelativeXY(nwPt, node, float64(zoom))
					coords = append(coords, []float64{x, y})
					prevWithinBounds = true
					prevCoords = nil
				} else {
					x, y := getRelativeXY(nwPt, node, float64(zoom))
					a1 := s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon()))
					if prevWithinBounds == true {
						coords = append(coords, []float64{x, y})
						prevCoords = nil //already added to coords
					} else {

						if s2.SimpleCrossing(p1, p2, a1, prevS2Point) || s2.SimpleCrossing(p1, p3, a1, prevS2Point) || s2.SimpleCrossing(p2, p4, a1, prevS2Point) {
							if len(prevCoords) > 0 {
								coords = append(coords, prevCoords)
							}
							coords = append(coords, []float64{x, y})
							prevCoords = nil
						} else {
							if len(coords) > 0 {
								drawPolyLine(img, color.Black, coords)
								coords = coords[:0]
							}
							prevCoords = []float64{x, y}
						}
					}
					prevS2Point = a1
					prevWithinBounds = false
				}
			}

			if len(coords) > 0 {
				drawPolyLine(img, color.Black, coords)
			}

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

	// TODO check area tag ?
	// path.Close()
	path.Stroke()
}
